//+build linux darwin

package icon

import "icotray/internal/pkg/image"

func IsValidMimeType(mimeType string) (bool, error) {
	isValid := image.IsImageMimeType(mimeType)

	if !isValid {
		return false, errors.New("the provided file type for the icon is invalid")
	}

	return true, nil
}
