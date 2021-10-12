package acticon

import (
	"fmt"
	"github.com/lxn/walk"
	"icotray/internal/icotray/acticon/action"
	"icotray/internal/icotray/acticon/icon"
	"icotray/internal/pkg/file"
	"strings"
)

type Configuration struct {
	Title         string
	IconPath      string
	HoverText     string
	DefaultAction string
	ActionItems   []ActionItem
	AppendQuit    bool
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

func (config *Configuration) getIcon() (*walk.Icon, error) {
	var resultIcon *walk.Icon

	// get the bytes of the icon
	imageBytes, err := icon.ImportIcon(config.IconPath)
	if err != nil {
		return nil, err
	}

	// verify the file mime type of the bytes
	mimeType := file.ExtractMimeTypeFromBytes(imageBytes)
	if isValid, err := icon.IsValidMimeType(mimeType); !isValid || err != nil {
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("the provided file type for the icon is invalid")
	}

	tempFilePath, err := icon.CreateTemporaryFile(imageBytes)
	if err != nil {
		return nil, err
	}

	// create a walk.Icon from the image
	resultIcon, err = walk.NewIconFromFile(tempFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not create the icon from the image: %v", err)
	}

	return resultIcon, nil
}

func (actionItem *ActionItem) Execute() {
	action.ExecuteVerbosely(actionItem.Action)
}
