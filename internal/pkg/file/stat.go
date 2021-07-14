package file

import (
	"net/http"
	"os"
)

func ExtractMimeTypeFromBytes(fileBytes []byte) string {
	return http.DetectContentType(fileBytes)
}

func DoesExist(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}
