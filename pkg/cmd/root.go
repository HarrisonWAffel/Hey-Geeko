package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "geeko",
	Short: "A simple voice assistant",
	Long:  `A simple voice assistant build using open-source software designed to run on raspberry pi's`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFile = "pkg/cmd/.agent-config.yml"
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("COBRACLISAMPLES")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
}
