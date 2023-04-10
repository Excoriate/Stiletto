package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func DirExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", dir)
	}

	return nil
}

func PathIsADirectory(path string) error {
	info, err := PathExist(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	return nil
}

func DirIsNotEmpty(dir string) error {
	if err := DirExist(dir); err != nil {
		return err
	}

	if err := PathIsADirectory(dir); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		return fmt.Errorf("failed to read content of directory %s", dir)
	}

	if len(entries) == 0 {
		return fmt.Errorf("directory %s is empty", dir)
	}

	return nil
}

func DirIsValid(dir string) error {
	if err := DirExist(dir); err != nil {
		return err
	}

	if err := PathIsADirectory(dir); err != nil {
		return err
	}

	return nil
}

func IsSubDir(parentDir string, childDir string) error {
	if err := DirExist(parentDir); err != nil {
		return err
	}

	if err := PathIsADirectory(parentDir); err != nil {
		return err
	}

	if err := DirExist(childDir); err != nil {
		return err
	}

	if err := PathIsADirectory(childDir); err != nil {
		return err
	}

	parentDir = filepath.Clean(parentDir)
	childDir = filepath.Clean(childDir)

	relativePath, err := filepath.Rel(parentDir, childDir)
	if err != nil {
		return fmt.Errorf("the child directory %s is not a subdirectory of %s", childDir, parentDir)
	}

	if filepath.HasPrefix(relativePath, "..") {
		return fmt.Errorf("the child directory %s is not a subdirectory of %s", childDir, parentDir)
	}

	return nil
}
