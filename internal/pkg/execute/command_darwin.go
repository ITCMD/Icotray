// +build darwin

package execute

import (
	"os/exec"
)

func Start(input string) *exec.Cmd {
	return exec.Command("open", input)
}

func StartWith(input string, appName string) *exec.Cmd {
	return exec.Command("open", "-a", appName, input)
}
