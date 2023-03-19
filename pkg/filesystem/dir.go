package filesystem

import (
	"os"
)

func CheckIfDirExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return err
	}

	return nil
}
