package filesystem

import (
	"fmt"
	"os"
)

func FileExist(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory, not a file", filePath)
	}

	return nil
}
