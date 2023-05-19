package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func IsPathAbsolute(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path %s is not absolute", path)
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

	// If the child dir isn't passed absolute, make it so.
	var childDirAbs string
	if err := IsPathAbsolute(childDir); err != nil {
		childDirFull := filepath.Join(filepath.Clean(parentDir), filepath.Clean(childDir))
		childDirAbs, err = PathToAbsolute(childDirFull)
		if err != nil {
			return err
		}
	} else {
		childDirAbs = filepath.Clean(childDir)
	}

	if err := DirExist(childDirAbs); err != nil {
		return err
	}

	if err := PathIsADirectory(childDirAbs); err != nil {
		return err
	}

	relativePath, err := filepath.Rel(parentDir, childDirAbs)
	if err != nil {
		return fmt.Errorf("the child directory %s is not a subdirectory of %s", childDirAbs, parentDir)
	}

	if strings.HasPrefix(relativePath, "..") {
		return fmt.Errorf("the child directory %s is not a subdirectory of %s", childDirAbs, parentDir)
	}

	return nil
}
