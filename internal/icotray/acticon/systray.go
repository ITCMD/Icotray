package acticon

import (
	"errors"
	"fmt"
	"icotray/internal/icotray/acticon/icon"
	"icotray/internal/pkg/dstruct/channel"
	"icotray/internal/pkg/file"
	"log"
	"strings"
	"time"

	"github.com/lxn/walk"
)

func CreateFromConfig(config *Configuration) error {
	if isValid, err := config.IsValid(); !isValid || err != nil {
		if err != nil {
			return err
		}
		return errors.New("the acticon configuration is invalid. check you input")
	}

	return createWindowAndNotifyIcon(config)
}

func createWindowAndNotifyIcon(config *Configuration) (error error) {
	// either a walk.MainWindow or walk.Dialog is required in order to use the message loop
	// however, the window will not be visible by default
	mainWindow, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	// create the actual notifyIcon and attach it to the mainWindow
	notifyIcon, err := config.createAndConfigureNotifyIcon(mainWindow)
	if err != nil {
		return err
	}

	// dispose the notify icon once the message loop is closed
	defer func(notifyIcon *walk.NotifyIcon) {
		err := notifyIcon.Dispose()
		if err != nil {
			error = err
		}
	}(notifyIcon)

	// add the menuItems from the configuration
	if err := addMenuItems(config, notifyIcon); err != nil {
		return err
	}

	// Run the message loop.
	mainWindow.Run()

	return nil
}

func (config *Configuration) createAndConfigureNotifyIcon(mainWindow *walk.MainWindow) (*walk.NotifyIcon, error) {
	// Create the notify icon and make sure we clean it up on exit.
	notifyIcon, err := walk.NewNotifyIcon(mainWindow)
	if err != nil {
		return nil, fmt.Errorf("could not create windows notifyicon: %v", err)
	}

	// get the icon from the configuration
	notifyIconIcon, err := config.getIcon()
	if err != nil {
		return nil, err
	}

	// set the icon of the notifyIcon
	if err := notifyIcon.SetIcon(notifyIconIcon); err != nil {
		return nil, fmt.Errorf("could not set the icon of the windows notifyicon: %v", err)
	}

	// set the notifyIcon tooltip if it is set
	if len(strings.TrimSpace(config.HoverText)) > 0 {
		if err := notifyIcon.SetToolTip(config.HoverText); err != nil {
			return nil, fmt.Errorf("could not set the tooltip of the windows notifyicon: %v", err)
		}
	}

	if len(strings.TrimSpace(config.DefaultAction)) > 0 {
		defaultActionItem := &ActionItem{
			Action: config.DefaultAction,
		}

		notifyIconClicked := make(chan interface{})
		notifyIconDoubleClicked := channel.Capacity(notifyIconClicked, 2, time.Millisecond * 500, false)

		go func() {
			for range notifyIconDoubleClicked {
				executeAction(*defaultActionItem)
			}
		}()

		// attach a mouse click handler to the notifyIcon for the default action
		notifyIcon.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
			if button == walk.LeftButton {
				notifyIconClicked <- button
			}
		})
	}

	// as the notifyIcon is hidden initially, the visibility must be set to true
	if err := notifyIcon.SetVisible(true); err != nil {
		return nil, fmt.Errorf("could not change the visibility of the windows notifyicon %v", err)
	}

	return notifyIcon, nil
}

func addMenuItems(config *Configuration, notifyIcon *walk.NotifyIcon) error {
	for _, actionItem := range config.ActionItems {
		action, err := addMenuItem(actionItem, notifyIcon)
		if err != nil {
			return err
		}

		// report any triggers to the channel
		action.Triggered().Attach(func() {
			executeAction(actionItem)
		})
	}

	if config.AppendQuit {
		quitActionItem := ActionItem{
			Title:   "Quit",
			Tooltip: "Quit this icotray instance",
			Action:  "",
		}

		if err := notifyIcon.ContextMenu().Actions().Add(walk.NewSeparatorAction()); err != nil {
			return fmt.Errorf("could not add separator action to context menu: %v", err)
		}

		quitAction, err := addMenuItem(quitActionItem, notifyIcon)
		if err != nil {
			return err
		}

		quitAction.Triggered().Attach(func() {
			quit()
		})
	}

	return nil
}

func addMenuItem(actionItem ActionItem, notifyIcon *walk.NotifyIcon) (*walk.Action, error) {
	// create the action
	action := walk.NewAction()
	if err := action.SetText(actionItem.Title); err != nil {
		return nil, fmt.Errorf("could not set the title '%v' for the menu item", actionItem.Title)
	}

	// set the action tooltip
	if err := action.SetToolTip(actionItem.Tooltip); err != nil {
		return nil, fmt.Errorf("could not set the tooltip '%v' for the menu item", actionItem.Tooltip)
	}

	// add the action to the context menu
	if err := notifyIcon.ContextMenu().Actions().Add(action); err != nil {
		return nil, fmt.Errorf("could not add the menu item '%v' to the acticon", actionItem.Title)
	}

	return action, nil
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

func executeAction(actionItem ActionItem) {
	out, err := actionItem.Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(out)
}

func quit() {
	walk.App().Exit(0)
}
