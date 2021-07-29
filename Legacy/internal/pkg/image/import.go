package image

import (
	"fmt"
	"icotray/internal/pkg/file"
	"strings"
)

func Import(path string) ([]byte, error) {
	imageBytes, err := file.ReadBytes(path)
	if err != nil {
		return nil, err
	}

	if !IsImageMimeType(imageBytes) {
		return nil, fmt.Errorf("file '%v' is not an image type", path)
	}

	return imageBytes, nil
}

func IsImageMimeType(imageBytes []byte) bool {
	mimeType := file.ExtractMimeTypeFromBytes(imageBytes)
	return strings.HasPrefix(mimeType, "image/")
}
