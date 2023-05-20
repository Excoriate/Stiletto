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
	parentDir = filepath.Clean(parentDir)
	childDir = filepath.Clean(childDir)

	// If the parentDir and childDir aren't passed absolute, make them so.
	parentDirAbs, err := filepath.Abs(parentDir)
	if err != nil {
		return err
	}

	var childDirAbs string
	if filepath.IsAbs(childDir) {
		childDirAbs = childDir
	} else {
		childDir = filepath.Join(parentDirAbs, childDir)
		childDirAbs, err = filepath.Abs(childDir)
		if err != nil {
			return err
		}
	}

	// Check if parentDir and childDir exists and are directories
	if err := DirExist(parentDirAbs); err != nil {
		return err
	}
	if err := PathIsADirectory(parentDirAbs); err != nil {
		return err
	}
	if err := DirExist(childDirAbs); err != nil {
		return err
	}
	if err := PathIsADirectory(childDirAbs); err != nil {
		return err
	}

	// Check if childDirAbs is a subdirectory of parentDirAbs
	relativePath, err := filepath.Rel(parentDirAbs, childDirAbs)
	if err != nil {
		return err
	}

	if strings.HasPrefix(relativePath, "..") || relativePath == "." || strings.HasPrefix(relativePath, "/") {
		return fmt.Errorf("the child directory %s is not a subdirectory of %s", childDirAbs, parentDirAbs)
	}

	return nil
}
