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

func IsPathAbsolute(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path %s is not absolute", path)
	}

	return nil
}

func IsPathRelative(path string) error {
	if filepath.IsAbs(path) {
		return fmt.Errorf("path %s is not relative", path)
	}

	return nil
}

func IsRelativeChildPath(parentPath string, childPath string) error {
	// Check if the child path is absolute
	if filepath.IsAbs(childPath) {
		return fmt.Errorf("child path %s is not relative", childPath)
	}

	// Resolve the child path relative to the parent path
	absChildPath := filepath.Join(parentPath, childPath)

	// Check if the absolute child path exists
	if _, err := os.Stat(absChildPath); os.IsNotExist(err) {
		return fmt.Errorf("child path %s does not exist", absChildPath)
	}

	// Check if the absolute child path is indeed a child of the parent path
	if !filepath.HasPrefix(absChildPath, parentPath) {
		return fmt.Errorf("child path %s is not a child of parent path %s", absChildPath, parentPath)
	}

	return nil
}
