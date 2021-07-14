//+build windows

package icon

import "errors"

func IsValidMimeType(mimeType string) (bool, error) {
	isValid := mimeType == "image/vnd.microsoft.icon" || mimeType == "image/x-icon"

	if !isValid {
		return false, errors.New("windows systray only accepts '.ico' files")
	}

	return true, nil
}
