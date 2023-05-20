package pipeline

import (
	"github.com/Excoriate/stiletto/internal/logger"
	"github.com/Excoriate/stiletto/pkg/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckPreConditions(t *testing.T) {
	logPrinter := logger.NewLogger()
	logPrinter.InitLogger()

	var args *config.PipelineOptions
	args = &config.PipelineOptions{}

	t.Run("WorkDirIsEmpty", func(t *testing.T) {
		args.WorkDir = ""
		args.TaskName = "workdir-is-empty"
		err := CheckPreConditions(args, logPrinter)

		assert.NoError(t, err, "The CheckPreConditions should not return an error")
	})

	// The workdir is passed as a "." string
	t.Run("WorkDirIsDot", func(t *testing.T) {
		args.WorkDir = "."
		args.TaskName = "workdir-is-dot"
		err := CheckPreConditions(args, logPrinter)

		assert.NoError(t, err, "The CheckPreConditions should not return an error")
	})

	// The workdir is passed as an invalid path, so it's expected an error.
	t.Run("WorkDirIsInvalid", func(t *testing.T) {
		args.WorkDir = "invalid-path"
		args.TaskName = "workdir-is-invalid"
		err := CheckPreConditions(args, logPrinter)

		assert.Error(t, err, "The CheckPreConditions should return an error")
	})

	// The workdir is passed as an absolute path, that does exist. No error is expected.
	t.Run("WorkDirIsAbsolute", func(t *testing.T) {
		currentDir, _ := os.Getwd()
		args.WorkDir = currentDir
		args.TaskName = "workdir-is-absolute"
		err := CheckPreConditions(args, logPrinter)

		assert.NoError(t, err, "The CheckPreConditions should not return an error")
	})

	// The workdir is valid, but the mountDir does not exist. An error is expected.
	t.Run("MountDirDoesNotExist", func(t *testing.T) {
		currentDir, _ := os.Getwd()
		args.WorkDir = currentDir
		args.MountDir = "invalid-mount-dir"
		args.TaskName = "mount-dir-does-not-exist"
		err := CheckPreConditions(args, logPrinter)

		assert.Error(t, err, "The CheckPreConditions should return an error")
	})

	// The workdir is valid, and the mountDir is passed as ".". No error is expected.
	t.Run("MountDirIsDot", func(t *testing.T) {
		currentDir, _ := os.Getwd()
		args.WorkDir = currentDir
		args.MountDir = "."
		args.TaskName = "mount-dir-is-dot"
		err := CheckPreConditions(args, logPrinter)

		assert.NoError(t, err, "The CheckPreConditions should not return an error")
	})

	// The workdir is valid, the mountDir is passed as an absolute path but it doesn't exist. An error is expected.
	t.Run("MountDirIsAbsolute", func(t *testing.T) {
		currentDir, _ := os.Getwd()
		args.WorkDir = currentDir
		args.MountDir = filepath.Join(currentDir, "invalid-mount-dir")
		args.TaskName = "mount-dir-is-absolute"
		err := CheckPreConditions(args, logPrinter)

		assert.Error(t, err, "The CheckPreConditions should return an error")
	})

	// The workdir is valid, the mountdir is also passed as a relative path that actually exist,
	//but the targetdir does not exist.
	t.Run("TargetDirDoesNotExist", func(t *testing.T) {
		currentDir, _ := os.Getwd()
		args.WorkDir = currentDir
		args.MountDir = "."
		args.TargetDir = "invalid-target-dir"
		args.TaskName = "target-dir-does-not-exist"
		err := CheckPreConditions(args, logPrinter)

		assert.Error(t, err, "The CheckPreConditions should return an error")
	})

	// The targetDir is a relative path that actually exist, so no error is expected.
	t.Run("TargetDirIsRelative", func(t *testing.T) {
		currentDir, _ := os.Getwd()
		args.WorkDir = currentDir
		args.MountDir = "."
		args.TargetDir = "."
		args.TaskName = "target-dir-is-relative"
		err := CheckPreConditions(args, logPrinter)

		assert.NoError(t, err, "The CheckPreConditions should not return an error")
	})
}
