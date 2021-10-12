package acticon

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func createAndConfigureNotifyIcon(mainWindow *walk.MainWindow, config *Configuration) (*walk.NotifyIcon, error) {
	// attach a new notifyIcon to the passed in window
	notifyIcon, err := walk.NewNotifyIcon(mainWindow)
	if err != nil {
		return nil, fmt.Errorf("could not create windows notifyicon: %v", err)
	}

	// get the icon from the configuration
	notifyIconIcon, err := config.getIcon()
	if err != nil {
		return nil, err
	}

	// set the icon of the actual notifyIcon
	if err := notifyIcon.SetIcon(notifyIconIcon); err != nil {
		return nil, fmt.Errorf("could not set the icon of the windows notifyicon: %v", err)
	}

	// if the tooltip is configured and not empty set it on the notifyIcon
	if len(strings.TrimSpace(config.HoverText)) > 0 {
		if err := notifyIcon.SetToolTip(config.HoverText); err != nil {
			return nil, fmt.Errorf("could not set the tooltip of the windows notifyicon: %v", err)
		}
	}

	configureDefaultAction(mainWindow, notifyIcon, config)

	// as the notifyIcon is hidden initially, the visibility must be set to true
	if err := notifyIcon.SetVisible(true); err != nil {
		return nil, fmt.Errorf("could not change the visibility of the windows notifyicon %v", err)
	}

	return notifyIcon, nil
}

func addMenuItems(config *Configuration, notifyIcon *walk.NotifyIcon) error {
	for _, actionItem := range config.ActionItems {
		actionItem := actionItem

		action, err := addMenuItem(actionItem, notifyIcon)
		if err != nil {
			return err
		}

		// report any triggers to the channel
		action.Triggered().Attach(func() {
			actionItem.Execute()
		})
	}

	if err := configureQuitAction(notifyIcon, config); err != nil {
		return err
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

func configureQuitAction(notifyIcon *walk.NotifyIcon, config *Configuration) error {
	if config.AppendQuit {
		// the static configuration of the quit action
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
			walk.App().Exit(0)
		})
	}

	return nil
}

func configureDefaultAction(window *walk.MainWindow, notifyIcon *walk.NotifyIcon, config *Configuration) {
	var primaryClickActionFn func()

	// if the default action is configured it will be run on click
	// otherwise the context menu is opened
	if len(strings.TrimSpace(config.DefaultAction)) > 0 {
		defaultActionItem := &ActionItem{
			Action: config.DefaultAction,
		}

		primaryClickActionFn = func() {
			defaultActionItem.Execute()
		}
	} else {
		primaryClickActionFn = func() {
			win.SendMessage(window.Handle(), win.WM_APP, 0, win.WM_CONTEXTMENU)
		}
	}

	// attach a mouse click handler to the notifyIcon for the default action
	notifyIcon.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			primaryClickActionFn()
		}
	})
}

func registerDisposalOfNotifyIconAtExit(notifyIcon *walk.NotifyIcon) (error error) {
	// wait for common termination signals, dispose the notifyIcon
	// and exit the application in order to remove zombie notifyIcons
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-exitSignal
		signal.Stop(exitSignal)
		close(exitSignal)

		if err := notifyIcon.Dispose(); err != nil {
			error = fmt.Errorf("could not dispose the notifyicon")
		}

		walk.App().Exit(0)
		os.Exit(0)
	}()

	return nil
}
