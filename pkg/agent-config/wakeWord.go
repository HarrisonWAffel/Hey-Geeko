package agent_config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
)

func AddWakeWord(ccmd *cobra.Command, _ []string) {
	path := viper.GetString("voice2json.path") + "/etc/precise/"
	model, err := ccmd.Flags().GetString("model")
	if err != nil {
		panic(err)
	}

	params := model + ".params"
	txt := model + "txt"
	cp(model, path+model)
	cp(params, path+params)
	cp(txt, path+txt)
}

func cp(src, dst string) {
	_, err := os.Stat(dst)
	if err == nil {
		return // file already exists
	}

	f, err := os.Open(src)
	if err != nil {
		panic("could not open provided model " + src)
	}

	f2, err := os.Create(dst)
	if err != nil {
		panic("could not create new model file in v2json repository")
	}

	_, err = io.Copy(f2, f)
	if err != nil {
		panic("could not copy model file into v2json repository")
	}
}
