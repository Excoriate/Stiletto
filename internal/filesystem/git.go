package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsGitRepository(path string, checkUntilRoot bool, checkRecursivelyUntilIsFound bool) (string, error) {
	//if path is relative, convert it to absolute
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = filepath.Join(wd, path)
	}

	// Ensure the path exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", err
	}

	// Ensure the path is a directory
	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", path)
	}

	// Convert the path to an absolute path (this will join a relative path with the current working directory)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// If checkUntilRoot is true, get the root of the path
	if checkUntilRoot {
		absPath = filepath.Dir(absPath)
	}

	// If checkRecursivelyUntilIsFound is true, check each directory in the path from the path to the root
	if checkRecursivelyUntilIsFound {
		for {
			// Check for .git directory or file
			gitPath := filepath.Join(absPath, ".git")
			if _, err := os.Stat(gitPath); err == nil {
				return absPath, nil
			}

			// Move to the parent directory
			parentPath := filepath.Dir(absPath)
			if parentPath == absPath {
				break
			}
			absPath = parentPath
		}
		return "", fmt.Errorf("path is not a git repository: %s", path)
	}

	// Check for .git directory or file
	gitPath := filepath.Join(absPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		return "", fmt.Errorf("path is not a git repository: %s", path)
	}

	return absPath, nil
}
