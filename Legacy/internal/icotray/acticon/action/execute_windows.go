// +build windows

package action

import "syscall"

func Execute(action string) ([]byte, error) {
	cmd := getCommand(action)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	return cmd.CombinedOutput()
}
