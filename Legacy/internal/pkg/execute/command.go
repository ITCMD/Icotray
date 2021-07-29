// +build !windows,!darwin

package execute

import (
	"os/exec"
)

// http://sources.debian.net/src/xdg-utils/1.1.0~rc1%2Bgit20111210-7.1/scripts/xdg-open/
// http://sources.debian.net/src/xdg-utils/1.1.0~rc1%2Bgit20111210-7.1/scripts/xdg-mime/

func Start(input string) *exec.Cmd {
	return exec.Command("xdg-open", input)
}

func StartWith(input string, appName string) *exec.Cmd {
	return exec.Command(appName, input)
}
