package agent_config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"harrisonwaffel/assistant/pkg/utils"
	"os"
	"os/exec"
)

func UpdateSentences(cmd *cobra.Command, args []string) {

	// Check for a custom agent-config path
	path := viper.GetString("profile.path")
	if path == "" {
		path = "~/.local/share/voice2json/"
	}
	path = path + viper.GetString("profile.name")

	// read in old agent-config and rename it to sentences.ini.old
	err := os.Rename(path+"/sentences.ini", path+"/sentences.ini.old")
	if err != nil {
		panic(err)
	}

	// take local sentences.ini and copy it to profile path
	newSentences, err := os.ReadFile("sentences.ini")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path+"/sentences.ini", newSentences, 0644)
	if err != nil {
		panic(err)
	}
}

func Train(_ *cobra.Command, _ []string) {
	cmd := exec.Command(utils.V2j, "--debug", "--profile", viper.GetString("profile.name"), utils.Train)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(cmd.String())
		panic("failed to train!" + err.Error())
	}
}
