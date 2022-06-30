package agent

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"harrisonwaffel/assistant/pkg/conversation"
	"harrisonwaffel/assistant/pkg/utils"
	"io"
	"os/exec"
	"time"
)

const (
	LowerConfidenceThreshold = 0.4
	UpperConfidenceThreshold = 0.7
)

type Voice2JsonCommand struct {
	Profile string
	Debug   bool
}

func (v Voice2JsonCommand) ListenForWakeWord() (chan string, error) {

	var cmd *exec.Cmd
	if viper.GetString("wakeWord.ModelName") != "" && viper.GetString("wakeWord.ModelName") != "hey-mycroft-2.pb" {
		cmd = exec.Command(utils.V2j, utils.WaitForWakeWord, "--model", viper.GetString("voice2json.path")+"/etc/precise/"+viper.GetString("wakeWord.modelName"))
	} else {
		cmd = exec.Command(utils.V2j, utils.WaitForWakeWord)
	}

	outputChan := make(chan string)
	var output bytes.Buffer
	cmd.Stdout = &output

	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	go cmd.Wait() // we constantly listen for the wake word

	// todo; memory profile, can this run for a long time?
	if v.Debug {
		go func() {
			fmt.Println("---")
			fmt.Println(cmd.String())
			fmt.Println("---")
			// print out each occurrence of the wake word
			scanner := bufio.NewScanner(&output)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line)
			}
		}()
	}

	// waiting for 2 seconds ensures we don't get a false wake-word activation
	time.Sleep(2 * time.Second)
	output.Reset()
	go func() {
		for {
			if !CurrentlyListening {
				e := output.Bytes()
				if string(e) != "" {
					go utils.PlayOpen()
					outputChan <- string(e)
					output.Reset()
				}
			}
		}
	}()

	return outputChan, nil
}

func (v Voice2JsonCommand) ListenWithIntent(prettyPrint bool) (chan conversation.V2JResponse, error) {
	streamCmd := exec.Command(utils.V2j, "--profile", v.Profile, utils.StreamInput)
	classifyCmd := exec.Command(utils.V2j, "--profile", v.Profile, utils.DetermineIntent)
	fmt.Println(streamCmd.String())
	fmt.Println(classifyCmd.String())
	// pipe commands together to asynchronously stream intents
	r, w := io.Pipe()
	streamCmd.Stdout = w
	classifyCmd.Stdin = r

	intentOutput, err := classifyCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := streamCmd.Start(); err != nil {
		panic(err)
	}

	if err := classifyCmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(intentOutput)
	input := make(chan conversation.V2JResponse)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			// we always listen to avoid latency
			// incurred by starting up ASR.
			// We only send the information once we
			// get a wake-word. If we're thinking of
			// a response, we don't listen or act on wake-words.
			if CurrentlyListening {
				resp := conversation.V2JResponse{}
				json.Unmarshal([]byte(line), &resp)

				// only accept a response if we are more than 40% confident of what we heard
				// this is technically intent confidence, and not speech confidence (which doesn't seem to work).
				if resp.Likelihood >= UpperConfidenceThreshold {
					if prettyPrint {
						go func() {
							speechInputChan <- resp
						}()
					}
					go utils.PlayClose()
					input <- resp
				}

				// if we are between 40% and 70% sure, we should ask if they can say it again. (this sort of works, but not well)
				if resp.Likelihood > LowerConfidenceThreshold && resp.Likelihood < UpperConfidenceThreshold {
					// this could be improved imho
					CurrentlyListening = false
					go v.speak("Can you say that again?")
					CurrentlyListening = true
				}
			}
		}
	}()

	go streamCmd.Wait() // we are always listening so this will never complete
	return input, nil
}
