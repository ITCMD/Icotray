package acticon

import (
	"fmt"
	"icotray/internal/icotray/acticon/action"
	"icotray/internal/pkg/file"
	"strings"
)

type Configuration struct {
	Title       string
	IconPath    string
	HoverText   string
	ActionItems []ActionItem
	AppendQuit  bool
}

type ActionItem struct {
	Title   string
	Tooltip string
	Action  string
}

func (config *Configuration) IsValid() (bool, error) {

	if len(strings.TrimSpace(config.IconPath)) > 0 {
		doesExist, err := file.DoesExist(config.IconPath)
		if err != nil {
			return false, fmt.Errorf("there was an error while accessing the file '%v'. ensure it exists and is not locked", config.IconPath)
		}
		if !doesExist {
			return false, fmt.Errorf("the file '%v' does not exist", config.IconPath)
		}
	}

	return true, nil
}

func (actionItem *ActionItem) Execute() (string, error) {
	out, err := action.Execute(actionItem.Action)
	if err != nil {
		return "", fmt.Errorf("could not execute action: %v", err)
	}

	outStr := string(out)
	if len(outStr) < 1 {
		outStr = "<no output>"
	}
	return fmt.Sprintf("action executed: %v \n", outStr), nil
}
