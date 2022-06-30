package cmd

import (
	"github.com/spf13/cobra"
	"harrisonwaffel/assistant/pkg/agent-config"
)

// profileCmd represents the agent-config command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Voice Profiles and update intents",
}

var updateCmd = &cobra.Command{
	Use:     "update-sentences",
	Aliases: []string{"u"},
	Short:   "replace the current profiles sentences.ini with your custom version. The old sentences.ini will be renamed to sentences.ini.old",
	Run:     agent_config.UpdateSentences,
}

var trainCmd = &cobra.Command{
	Use:   "train",
	Short: "train the voice model on sentences.ini",
	Run:   agent_config.Train,
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(updateCmd)
	profileCmd.AddCommand(trainCmd)
}
