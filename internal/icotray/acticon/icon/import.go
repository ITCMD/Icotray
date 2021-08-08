package icon

import (
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"fmt"
	"icotray/assets"
	"icotray/internal/pkg/image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	tempFilePrefix = "icotray_temp_icon"
)

func ImportIcon(path string) ([]byte, error) {
	if len(strings.TrimSpace(path)) > 1 {
		return image.Import(path)
	}

	return assets.IcotrayIcon, nil
}

func CreateTemporaryFile(iconBytes []byte) (path string, err error) {
	md5Checksum := md5.Sum(iconBytes)
	dataHash := hex.EncodeToString(md5Checksum[:])

	iconFilePath := filepath.Join(os.TempDir(), fmt.Sprintf("%v_%v.ico", tempFilePrefix, dataHash))

	if _, err := os.Stat(iconFilePath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(iconFilePath, iconBytes, 0644); err != nil {
			return "", fmt.Errorf("could not wirte temporary icon file to '%v': %v", iconFilePath, err)
		}
	}
	return iconFilePath, nil
}
