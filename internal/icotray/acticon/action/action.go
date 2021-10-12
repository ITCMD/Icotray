package action

import "fmt"

func ExecuteVerbosely(action string) {
	out, err := ExecuteAndPrepareOutput(action)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(out)
}

func ExecuteAndPrepareOutput(action string) (string, error) {
	out, err := Execute(action)
	if err != nil {
		return "", fmt.Errorf("could not execute action: %v", err)
	}

	outStr := string(out)
	if len(outStr) < 1 {
		outStr = "<no output>"
	}
	return fmt.Sprintf("action executed: %v \n", outStr), nil
}
