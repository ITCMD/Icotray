package action

import (
	"fmt"
	"icotray/internal/pkg/execute"
	"os/exec"
	"regexp"
	"strings"
)

var (
	commandActionPrefix = "cmd:"
	whitespacePadding   = "#whtspcpddng#"
)

func getCommand(action string) *exec.Cmd {
	commandAction := getCommandAction(action)

	// if commandAction is empty, the action will be treated as a regular command
	if commandAction == "" {
		return regularCommand(action)
	} else {
		return commandActionCommand(commandAction)
	}
}

func regularCommand(action string) *exec.Cmd {
	return execute.Start(action)
}

func commandActionCommand(action string) *exec.Cmd {
	unescapedWhitespaceRe := regexp.MustCompile(`([^\\])\s`)
	escapedWhitespaceRe := regexp.MustCompile(`\\\s`)
	commandletsRe := regexp.MustCompile(fmt.Sprintf(`%v\s`, whitespacePadding))

	// add padding of '#' to the unescaped whitespaces in the action because go does not support lookarounds
	actionWithPaddedWhitespace := unescapedWhitespaceRe.ReplaceAllString(action, fmt.Sprintf(`$1%v `, whitespacePadding))

	// replace the escaped whitespaces in the padded action with normal ones
	normalizedActionWithPaddedWhitespace := escapedWhitespaceRe.ReplaceAllString(actionWithPaddedWhitespace, " ")

	// split the padded action by the whitespaces which are padded
	commandlets := commandletsRe.Split(normalizedActionWithPaddedWhitespace, -1)

	// construct the command from the first value as the command name and the rest as arguments
	cmd := exec.Command(commandlets[0], commandlets[1:]...)

	return cmd
}

func getCommandAction(action string) string {
	isCommandAction := strings.HasPrefix(action, commandActionPrefix)

	if isCommandAction {
		return strings.TrimPrefix(action, commandActionPrefix)
	}
	return ""
}
