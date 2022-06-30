package agent

import (
	"encoding/json"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"harrisonwaffel/assistant/pkg/agent/cmds/music"
	"harrisonwaffel/assistant/pkg/agent/cmds/weather"
	"harrisonwaffel/assistant/pkg/conversation"
	"strings"
	"time"
)

func init() {
	speechInputChan = make(chan conversation.V2JResponse)
	WordChan = make(chan string)
}

var (
	CurrentlyListening = false
	CurrentlyThinking  = false

	// following channels are only used for pretty print

	speechInputChan chan conversation.V2JResponse
	WordChan        chan string
)

type agentContext struct {
	TTSProvider  string
	WeatherCache weather.Cache
	MusicClient  *music.MusicClient
	PrettyPrint  bool
	Voice2JsonCommand
}

func StartAgent(cmd *cobra.Command, args []string) {
	if err := coreLoop(buildCtx(cmd)); err != nil {
		panic(err)
	}
}

func buildCtx(cmd *cobra.Command) agentContext {
	ctx := agentContext{
		TTSProvider:  "espeak",
		WeatherCache: weather.WeatherCache,
		MusicClient:  music.NewMusicClient(),
		PrettyPrint:  false,
		Voice2JsonCommand: Voice2JsonCommand{
			Profile: viper.GetString("profile.name"),
		},
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		panic(err)
	}
	ctx.Debug = debug

	pretty, err := cmd.Flags().GetBool("pretty-print")
	if err != nil {
		panic(err)
	}

	ctx.PrettyPrint = pretty

	return ctx
}

func coreLoop(ctx agentContext) error {
	go startupMessage(&ctx)

	// geeko is always listening, it sends input via the returned channel
	input, err := ctx.ListenWithIntent(ctx.PrettyPrint)
	if err != nil {
		panic(err)
	}

	incomingSpeech, err := ctx.ListenForWakeWord()
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-incomingSpeech:
			if !CurrentlyListening {

				musicPlaying := false
				if ctx.MusicClient.Ctl != nil && !ctx.MusicClient.Ctl.Paused {
					musicPlaying = true
					ctx.MusicClient.TogglePause()
				}

				CurrentlyListening = true
				resp := <-input
				CurrentlyListening = false
				handleIntent(&ctx, resp)

				if musicPlaying {
					ctx.MusicClient.TogglePause()
				}
			}
		}
	}

	return err
}

func startupMessage(ctx *agentContext) {
	fmt.Println("Starting Up Geeko...")
	s := `
If you're using the default wake-word, say "Hey, Geeko" to activate Geeko and make a command.
Background noise WILL reduce accuracy! Further training the provided model, or making your own, can help accuracy.

NOTE: On first startup, the wake-word will take longer to recognize, after the first time it should be faster.

`
	fmt.Println(s)

	if ctx.PrettyPrint {
		prettyPrint(ctx)
	}
}

func prettyPrint(ctx *agentContext) {
	f := figure.NewFigure("Hey, Geeko", "doom", true)
	fin := f.String() + "\n" + "A simple voice assistant created during SUSE Hack-week"
	header := pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Sprint(fin)
	area, _ := pterm.DefaultArea.WithCenter().Start()
	var j conversation.V2JResponse
	var words []string

	x := ""
	area.Clear()
	var oldPanels string
	for {
		time.Sleep(2 * time.Second)
		x = fmt.Sprintf("Currently Listening: %t\n", CurrentlyListening)
		x = x + fmt.Sprintf("Currently Thinking : %t\n", CurrentlyThinking)
		x = x + fmt.Sprintf("Music Playing      : %t\n", ctx.MusicClient.Ctl != nil)
		select {
		case m := <-speechInputChan:
			j = m
		case w := <-WordChan:
			if len(words) == 50 {
				words = []string{}
			}
			words = append(words, w)
		default:
		}
		panel := pterm.DefaultBox.Sprintf(header)
		panel1 := pterm.DefaultBox.WithTitle("speaker status").Sprint(x)
		panel2 := ""
		if ctx.Debug {
			js, _ := json.MarshalIndent(j, "", " ")
			panel2 = pterm.DefaultBox.WithTitle("most recent JSON").Sprint(string(js))
		} else {
			panel2 = pterm.DefaultBox.WithTitle("most recent Input").Sprint(strings.Join(j.RawTokens, " "))
		}
		panel3 := pterm.DefaultBox.WithTitle("recent speech output").WithTitleBottomCenter().Sprint(strings.Join(words, "\n"))
		panels, _ := pterm.DefaultPanel.WithPanels(pterm.Panels{
			{{Data: panel}},
			{{Data: panel1}, {Data: panel3}},
			{{Data: panel2}},
		}).Srender()

		if oldPanels == "" {
			oldPanels = panels
			area.Update(panels)
		} else if oldPanels != panels {
			area.Update(panels)
		}
	}
}
