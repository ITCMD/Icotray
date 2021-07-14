package add

import (
	"errors"
	"github.com/spf13/cobra"
	"icotray/internal/icotray/acticon"
)

func validateArgs(cmd *cobra.Command, args []string) error {
	// as the arguments are optional for now, return nil if none were provided
	if len(args) < 1 {
		return nil
	}

	if len(args) > 1 {
		return errors.New("only one argument for 'identifier' is allowed")
	}

	// validate the 'identifier' argument
	if isValid, err := acticon.Identifier(args[0]).IsValid(); !isValid || err != nil {
		return errors.New("the provided identifier is invalid. only use alphanumerical characters and '-' or '_'")
	}

	return nil
}

func (config *Configuration) isValid() (bool, error) {
	// when using the separate lists for the item names and actions
	// check if both have the equal amount of elements
	if len(config.actionItemNames) != len(config.actionItemActions) {
		return false, errors.New("the length of the menu-items passed with the 'item-name' flag does not match the number of 'item-action' flags")
	}

	return true, nil
}
