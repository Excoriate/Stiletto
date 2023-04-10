package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func PathExist(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path %s does not exist", path)
		}
		return nil, fmt.Errorf("error checking path %s: %s", path, err.Error())
	}

	return info, nil
}

func PathToAbsolute(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error converting path %s to absolute path: %s", path, err.Error())
	}

	return absolutePath, nil
}
