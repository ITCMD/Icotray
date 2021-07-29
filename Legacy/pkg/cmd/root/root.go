package root

import (
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "icotray",
	Short: "Create custom 'acticons' in the system tray",
	Long: `
Icotray is a CLI tool to create custom acticons in the system tray.
An acticon is an icon in the systemtray which has one or multiple actions associated with it.

The acticon will be available in the tray as long as the program is running
`,
	CompletionOptions: cobra.CompletionOptions{
		//DisableDefaultCmd: true,
	},
}

func Execute() {
	cobra.CheckErr(cmd.Execute())
}

func init() {
	setVersion()
	addCredits()
	addChildCommands()
	configureFlags()

	cobra.OnInitialize(initConfig)
}
