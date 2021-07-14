package execute

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	urlDll   = "url.dll,FileProtocolHandler"
	runDll32 = filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
)

func Start(input string) *exec.Cmd {
	cmd := exec.Command(runDll32, urlDll, input)
	return cmd
}

func StartWith(input string, appName string) *exec.Cmd {
	cmd := exec.Command("cmd", "/C", "start", "", appName, cleanInput(input))
	return cmd
}

func cleanInput(input string) string {
	r := strings.NewReplacer("&", "^&")
	return r.Replace(input)
}
