package file

import (
	"bufio"
	"fmt"
	"os"
)

func ReadBytes(path string) (bytes []byte, err error) {
	// open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file '%v'", path)
	}

	// close the file after reading the bytes
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			// override the named return values
			bytes = nil
			err = nil
		}
	}(file)

	// read the stats of the file
	var fileInfo os.FileInfo
	fileInfo, err = file.Stat()
	if err != nil {
		return nil, err
	}

	// initialize the resulting byte slice
	var fileSize int64 = fileInfo.Size()
	fileBytes := make([]byte, fileSize)

	// read the actual bytes from the file into the byte slice
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("could not read bytes of the file '%v'", path)
	}

	return fileBytes, nil
}
