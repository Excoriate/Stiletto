package config

import (
	"os"
)

var (
	currentDir, _ = os.Getwd()
	homeDir, _    = os.UserHomeDir()
)

type DefaultDirs struct {
	CurrentDir           string
	BinaryDir            string
	GitRepositoryRootDir string
	HomeDir              string
	BuildDirInContainer  string
}

func GetDefaultDirs() *DefaultDirs {
	return &DefaultDirs{
		CurrentDir:           currentDir,
		BinaryDir:            currentDir,
		GitRepositoryRootDir: "", // To be resolved later on.
		BuildDirInContainer:  "/build",
		HomeDir:              homeDir,
	}
}
