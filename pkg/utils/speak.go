package utils

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type TTSProvider string

// Speak will speak-out the provided sentences using the configured TTS program.
// Passing a slice of sentences will cause a pause between the speaking of each sentence.
func Speak(sentence ...string) {
	for _, s := range sentence {
		makeTTSReq(s)
	}
}

func makeTTSReq(sentence string) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:5500/api/tts", nil)
	if err != nil {
		panic("could not make tts request")
	}
	q := req.URL.Query()
	q.Add("text", sentence)
	q.Add("cache", "true")
	q.Add("voice", viper.GetString("tts.voice"))
	q.Add("vocoder", "low")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("Could not speak: " + err.Error()) // if tts doesn't work theres no reason to continue
	}
	cmd := exec.Command("aplay")
	stdIn, err := cmd.StdinPipe()
	if err != nil {
		panic("could not get stdin for aplay: " + err.Error())
	}
	defer resp.Body.Close()
	audio, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("could not parse TTS response")
	}
	if err := cmd.Start(); err != nil {
		panic("could not start aplay for TTS")
	}
	stdIn.Write(audio)
	cmd.Wait()
}
