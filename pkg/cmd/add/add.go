package add

import (
	"errors"
	"fmt"
	"icotray/internal/icotray/acticon"
	"icotray/internal/icotray/acticon/builder"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "add [<identifier>]",
	Short: "Adds an acticon to the system tray",
	Long: `
Adds an acticon to the system tray.
The configuration of the acticon may either be passed
through arguments and flags or created using the interactive mode.


# INTERACTIVE
If using the latter, icotray will prompt the values needed for
configuring the acticon. You may sill preconfigure other options when using
the interactive mode. You then have the option to overwrite the fields.
Please note that this feature may not work in all shells.


# ICON
If no own icon is provided via the command flag, the default icon will be used.
The file type of the icons depends on the operating system.
For example: Windows only accepts .ico files for the icons in the system tray.


# ACTION
The actions of the command can be provided as a list of key-value pairs.
The provided actions will be shown in a list format when clicking on the icon
in the system tray. Using the '--quittable' flag will ad an option to quit the program.
The key represents the title of the item and the key the action which will be run
when clicking on the item.

## RUN MODES
### DEFAULT
By default the actions will be run using the 'open with default program' method.
The concrete method will be chosen depending on the operating system:
    OSX         :  "open"
    Windows     :  "start"    ->  rundll32.exe url.dll,FileProtocolHandler
    Linux/Other :  "xdg-open"

### COMMAND
By adding a 'cmd:' prefix to the action, it will be executed as a command.
For this the first value after 'cmd:' has to be the command / program 
to use for the action (e.g. TASKKILL, TREE, bash ...).

After the first value, the following values are splitted into arguments by whitespace.

So >>"Name"="TASKKILL /PID 1234321"<< will result in: 
    Program     :  "TASKKILL"
    Arguments   :  ["/PID", "1234321"]

Sometimes it will be necessary to keep multiple values together as a single argument.
For this case you will have to escape the whitespace using a '\' (backslash).

So >>"Name"="cmd:bash -c echo\ \"Hello\ World\"\ >\ ~/myfile.txt"<< will result in
    Program     :  "bash"
    Arguments   :  ["-c", "echo Hello World > ~/myfile.txt"]

## Default Action
By providing an action with the '--default' / '-d' flag, the action will be
interpreted as the 'default' action. The action passed with the flag will
be run when double-clicking the acticon. In order to open the context menu,
the acticon will have to be right-clicked.


# EXAMPLES
## Basic acticon with some actions
icotray add --title "Acticon Title" --actions "Run xy"="xy","Start Firefox"="firefox"

## The actions may also be provided through separate flags
icotray add --title "Acticon Title" --actions "Run xy"="xy" --actions "Start Firefox"="firefox"

## The menu item names may also be provided as separate list flags
icotray add --title "Acticon Title" --item-name "Run xy" --item-action "xy" --item-name "Start Firefox" --item-action "firefox"

## Acticon with custom icon
icotray add --title "Acticon Title" --icon "/path/to/icon"

## Show message when hovering the acticon
icotray add --title "Acticon Title" --hover "Hovertext"
`,
	Args: validateArgs,
	RunE: runCommand,
}

func init() {
	configureFlags()
}

func runCommand(cmd *cobra.Command, args []string) error {
	if isValid, err := config.isValid(); !isValid || err != nil {
		if err != nil {
			return err
		}
		return errors.New("the configuration is invalid. check your input")
	}

	acticonConfig := config.toActiconConfig()

	if config.runInteractively {
		var err error
		if acticonConfig, err = builder.BuildInteractively(acticonConfig); err != nil {
			return err
		}
	}

	if config.printCommand {
		fmt.Println(acticonConfigToCommand(acticonConfig))
	}

	if err := acticon.CreateFromConfig(acticonConfig); err != nil {
		return err
	}

	return nil
}
