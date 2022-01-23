package goutil

import (
	"fmt"
	"os"
)

// OpenFile is a function for open file with path & filename
func OpenFile(path, filename string) (*os.File, error) {
	os.Mkdir(path, 0777)
	return os.OpenFile(
		fmt.Sprintf("%s/%s", path, filename),
		os.O_CREATE|os.O_APPEND|os.O_RDWR,
		0666)
}

// OpenFileNewest is a function for open file with path & filename and close the oldest os.File
func OpenFileNewest(oldest *os.File, path, filename string) (*os.File, error) {
	file, err := os.OpenFile(
		fmt.Sprintf("%s/%s", path, filename),
		os.O_CREATE|os.O_APPEND|os.O_RDWR,
		0666)
	if err != nil {
		return nil, err
	}

	if err = oldest.Close(); err != nil {
		return nil, err
	}

	return file, nil
}
