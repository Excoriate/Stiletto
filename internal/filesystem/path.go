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

func PathGetWorkDirAbsolute(workDir string) (string, error) {
	var workDirPath string
	if workDir == "." {
		wd, err := os.Getwd()

		if err != nil {
			return "", fmt.Errorf("error getting current working directory: %s", err.Error())
		}

		workDirPath = wd
	} else {
		// if it's a relative path, convert it to an absolute path
		if !filepath.IsAbs(workDir) {
			wd, err := os.Getwd()
			if err != nil {
				return "", fmt.Errorf("error getting current working directory: %s", err.Error())
			}

			workDirPath = filepath.Join(wd, workDir)
		} else {
			workDirPath = workDir
		}
	}

	return workDirPath, nil
}
