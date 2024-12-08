package files

import (
	"errors"
	"io"
	"os"
)

func OpenFileToByteStream(path string) ([]byte, error) {
	file, err := os.Open(path)
    if err != nil {
        return nil, errors.New(`Error opening config file: ` + err.Error())
    }
    defer file.Close()

	bytes, err := io.ReadAll(file)
    if err != nil {
        return nil, errors.New(`Error reading config file: ` + err.Error())
    }

	return bytes, err
}