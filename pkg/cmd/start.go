package cmd

import (
	"github.com/spf13/cobra"
	"harrisonwaffel/assistant/pkg/agent"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the voice agent",
	Run:   agent.StartAgent,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("debug", "d", false, "output debug logs")
	startCmd.Flags().BoolP("pretty-print", "p", false, "pretty print speaker info (requires large terminal window)")

}
