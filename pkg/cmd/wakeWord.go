/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	agent_config "harrisonwaffel/assistant/pkg/agent-config"

	"github.com/spf13/cobra"
)

// wakeWordCmd represents the wakeWord command
var wakeWordCmd = &cobra.Command{
	Use:   "wake-word",
	Short: "commands to update the wake-word",
	Run:   agent_config.AddWakeWord,
}

func init() {
	rootCmd.AddCommand(wakeWordCmd)
	wakeWordCmd.PersistentFlags().StringP("model", "m", "hey-geeko.pb", "the filename of the model should be added to voice2json. model must be stored in the root directory of this project")
}
