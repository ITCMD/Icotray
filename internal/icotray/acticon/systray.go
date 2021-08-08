package acticon

import (
	"errors"
	"github.com/lxn/walk"
	"log"
)

func CreateFromConfig(config *Configuration) (error error) {
	// validate the configuration before creating the window and notifyIcon
	if isValid, err := config.IsValid(); !isValid || err != nil {
		if err != nil {
			return err
		}
		return errors.New("the acticon configuration is invalid. check you input")
	}

	// either a walk.MainWindow or walk.Dialog is required in order to use the message loop
	// however, the window will not be visible by default
	mainWindow, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	// create the actual notifyIcon and attach it to the mainWindow
	notifyIcon, err := createAndConfigureNotifyIcon(mainWindow, config)
	if err != nil {
		return err
	}

	// dispose the notifyIcon once the message loop is closed or the application exits
	defer func(notifyIcon *walk.NotifyIcon) {
		err := notifyIcon.Dispose()
		if err != nil {
			error = err
		}
	}(notifyIcon)

	if err := registerDisposalOfNotifyIconAtExit(notifyIcon); err != nil {
		return err
	}

	// add the menuItems from the configuration
	if err := addMenuItems(config, notifyIcon); err != nil {
		return err
	}

	mainWindow.Run()

	return nil
}
