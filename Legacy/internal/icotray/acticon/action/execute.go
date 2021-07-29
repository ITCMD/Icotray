// +build !windows

package action

func Execute(action string) ([]byte, error) {
	cmd := getCommand(action)
	return cmd.CombinedOutput()
}
