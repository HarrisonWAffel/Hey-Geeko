package utils

import (
	"os/exec"
)

const (
	V2j             = "voice2json" // core command
	Train           = "train-profile"
	Download        = "download-profile"
	StreamInput     = "transcribe-stream"
	DetermineIntent = "recognize-intent"
	WaitForWakeWord = "wait-wake"
	SpeakCommand    = "speak-sentence"
	MaryTTSServer   = "marytts-server:59125" // server is assumed to be running
)

// PlayOpen initially tried to use "github.com/faiface/beep",
// but exec is just faster
func PlayOpen() {
	c := exec.Command("aplay", "pkg/microphone-open.wav")
	c.Output()
}

func PlayClose() {
	c := exec.Command("aplay", "pkg/microphone-close.wav")
	c.Output()
}
