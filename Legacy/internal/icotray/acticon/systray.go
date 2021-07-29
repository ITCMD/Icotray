package acticon

import (
	"errors"
	"fmt"
	"icotray/internal/icotray/acticon/icon"
	"icotray/internal/pkg/dstruct"
	"icotray/internal/pkg/dstruct/channel"
	"icotray/internal/pkg/file"
	"strings"

	"github.com/getlantern/systray"
)

type SystrayResult struct {
	Error error
}

func CreateFromConfig(config *Configuration) (err error) {
	if isValid, err := config.IsValid(); !isValid || err != nil {
		if err != nil {
			return err
		}
		return errors.New("the acticon configuration is invalid. check you input")
	}

	result := &SystrayResult{}

	systray.Run(bindSystrayData(onSystrayReady, config, result), bindSystrayData(onSystrayExit, config, result))

	return result.Error
}

func bindSystrayData(boundFn func(config *Configuration, result *SystrayResult), config *Configuration, result *SystrayResult) func() {
	return func() {
		boundFn(config, result)
	}
}

func onSystrayReady(config *Configuration, result *SystrayResult) {
	configureActicon(config, result)
	addMenuItemsAndWaitForClicks(config, result)
}

func onSystrayExit(config *Configuration, result *SystrayResult) {
	// nothing to do when systray exits
}

func configureActicon(config *Configuration, result *SystrayResult) {
	//! the icon needs to be set before the other fields
	//! because "github.com/getlantern/systray" panics otherwise
	if err := setSystrayIcon(config); err != nil {
		result.Error = err
		systray.Quit()
		return
	}

	if len(strings.TrimSpace(config.Title)) > 0 {
		systray.SetTitle(config.Title)
	}

	if len(strings.TrimSpace(config.HoverText)) > 0 {
		systray.SetTooltip(config.HoverText)
	}
}

func setSystrayIcon(config *Configuration) error {
	imageBytes, err := icon.ImportIcon(config.IconPath)
	if err != nil {
		return err
	}

	mimeType := file.ExtractMimeTypeFromBytes(imageBytes)
	if isValid, err := icon.IsValidMimeType(mimeType); !isValid || err != nil {
		if err != nil {
			return err
		}

		return fmt.Errorf("the provided file type for the icon is invalid")
	}

	systray.SetTemplateIcon(imageBytes, imageBytes)

	return nil
}

func addMenuItemsAndWaitForClicks(config *Configuration, result *SystrayResult) {
	clickedItemIndexes := addMenuItems(config)

	quitClicked := make(chan struct{})
	if config.AppendQuit {
		if len(config.ActionItems) > 0 {
			systray.AddSeparator()
		}
		menuItem := systray.AddMenuItem("Quit", "Quit the application")

		quitClicked = menuItem.ClickedCh
	}

	for {
		select {
		case clickedItemIndex := <-clickedItemIndexes:
			itemIndex := clickedItemIndex.(int)
			actionItem := config.ActionItems[itemIndex]

			out, err := actionItem.Execute()
			if err != nil {
				result.Error = err
				systray.Quit()
				return
			}
			fmt.Print(out)
		case <-quitClicked:
			systray.Quit()
			return
		}
	}
}

func addMenuItems(config *Configuration) (itemIndexes <-chan dstruct.Any) {
	var clickedItemIndexes []<-chan dstruct.Any

	for index, actionItem := range config.ActionItems {
		// copy index in order to keep it in the context of this action item
		itemIndex := index
		menuItem := systray.AddMenuItem(actionItem.Title, actionItem.Tooltip)

		clickedItemId := channel.MapReceivingVoidVendor(menuItem.ClickedCh, func() dstruct.Any {
			return itemIndex
		})

		clickedItemIndexes = append(clickedItemIndexes, clickedItemId)
	}

	return channel.MergeReceivingWithGoroutines(clickedItemIndexes...)
}
