package icon

import (
	_ "embed"
	"icotray/assets"
	"icotray/internal/pkg/image"
	"strings"
)


func ImportIcon(path string) ([]byte, error) {
	if len(strings.TrimSpace(path)) > 1 {
		return image.Import(path)
	}

	return assets.IcotrayIcon, nil
}
