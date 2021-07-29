package root

import (
	"fmt"
	"icotray/pkg/cmd/add"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	helpCredits = `
Program by Mnoronen for ITCMD https://github.com/ITCMD/icotray
`
)

func addChildCommands() {
	cmd.AddCommand(add.Cmd)
}

func addCredits() {
	defaultTemplate := cmd.HelpTemplate()
	watermarkedTemplate := fmt.Sprintf(`%v%v`, defaultTemplate, helpCredits)

	cmd.SetHelpTemplate(watermarkedTemplate)
}

func configureFlags() {
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.icotray.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".icotray")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		_, err := fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if err != nil {
			return
		}
	}
}
