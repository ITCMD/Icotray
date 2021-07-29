package str

import "strings"

func DropWhitespaceValues(values []string) []string {
	return Filter(values, func(value string) bool {
		return len(strings.TrimSpace(value)) > 0
	})
}

func Filter(values []string, includeFn func(value string) bool) []string {
	var result []string

	for _, value := range values {
		if includeFn(value) {
			result = append(result, value)
		}
	}

	return result
}
