package pipeline

import (
	"github.com/Excoriate/stiletto/internal/logger"
	"github.com/Excoriate/stiletto/pkg/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestIsWorkDirValid(t *testing.T) {
	currentDir, _ := os.Getwd()

	t.Run("Workdir is empty so it resolve current dir", func(t *testing.T) {
		workDirCfg, err := IsWorkDirValid("")
		assert.NoError(t, err, "The isWorkDirValid should not return an error")

		// expected object shape
		currentDir, _ := os.Getwd()
		cfg := config.PipelineDirs{
			Dir:  ".",
			Path: currentDir,
		}

		assert.Equal(t, cfg, workDirCfg, "The workDirCfg should be empty")
		assert.Equal(t, ".", workDirCfg.Dir, "The workDirCfg should be empty")
		assert.Equal(t, currentDir, workDirCfg.Path, "The workDirCfg should be empty")
	})

	// The workdir is passed as a "." string
	t.Run("Workdir is passed as a \".\" string", func(t *testing.T) {
		workDirCfg, err := IsWorkDirValid(".")
		assert.NoError(t, err, "The isWorkDirValid should not return an error")

		// expected object shape
		currentDir, _ := os.Getwd()
		cfg := config.PipelineDirs{
			Dir:  ".",
			Path: currentDir,
		}

		assert.Equal(t, cfg, workDirCfg, "The workDirCfg should be empty")
		assert.Equal(t, ".", workDirCfg.Dir, "The workDirCfg should be empty")
		assert.Equal(t, currentDir, workDirCfg.Path, "The workDirCfg should be empty")
	})

	t.Run("Workdir is passed as a relative path", func(t *testing.T) {
		// The workdir is passed as a relative path
		relativePathThatExist := "newfolder"
		// Create this folder in the current directory temporately.
		// This folder will be deleted after the test

		err := os.Mkdir(relativePathThatExist, 0755)
		assert.NoError(t, err, "The isWorkDirValid should not return an error")

		workDirCfg, err := IsWorkDirValid(relativePathThatExist)
		assert.NoError(t, err, "The isWorkDirValid should not return an error")

		// expected object shape

		cfg := config.PipelineDirs{
			Dir:  relativePathThatExist,
			Path: currentDir + "/" + relativePathThatExist,
		}

		assert.Equal(t, cfg, workDirCfg, "The workDirCfg should be empty")
		assert.Equal(t, relativePathThatExist, workDirCfg.Dir, "The workDirCfg should be empty")
		assert.Equal(t, currentDir+"/"+relativePathThatExist, workDirCfg.Path, "The workDirCfg should be empty")

		// Remove the folder
		err = os.Remove(relativePathThatExist)
		assert.NoError(t, err, "The isWorkDirValid should not return an error")
	})

	// The relative path doesn't exist, it didn't pass the validation.
	t.Run("Workdir is passed as a relative path that doesn't exist", func(t *testing.T) {
		relativePathThatDoesntExist := "newfolderthatdoesntexist"
		workDirCfg, err := IsWorkDirValid(relativePathThatDoesntExist)
		assert.Error(t, err, "The isWorkDirValid should return an error")
		assert.Equal(t, config.PipelineDirs{}, workDirCfg, "The workDirCfg should be empty")
	})

	// The workdir is passed, as an absolute path
	t.Run("Workdir is passed as an absolute path", func(t *testing.T) {
		absolutePathFromCurrentDirTwoLevelsBelow := "/../../"

		workDirCfg, err := IsWorkDirValid(absolutePathFromCurrentDirTwoLevelsBelow)
		assert.NoError(t, err, "The isWorkDirValid should not return an error")

		assert.Equal(t, absolutePathFromCurrentDirTwoLevelsBelow, workDirCfg.Path, "The workDirCfg should be empty")
	})
}

func TestIsMountOrTargetDirValid(t *testing.T) {
	logPrinter := logger.NewLogger()
	logPrinter.InitLogger()
	currentDir, _ := os.Getwd()

	t.Run("MountDir is passed empty, so it inherits the workDir values", func(t *testing.T) {
		workDirCfg, err := IsWorkDirValid(".")
		assert.NoError(t, err, "The isWorkDirValid should not return an error")

		mountDirCfg, err := IsMountOrTargetDirValid("", workDirCfg, "mount")

		assert.NoError(t, err, "The IsMountOrTargetDirValid should not return an error")
		assert.Equal(t, workDirCfg.Dir, mountDirCfg.Dir, "The mounted directory should be the same as the workdir")
		assert.Equal(t, workDirCfg.Path, mountDirCfg.Path, "The mounted directory should be the same as the workdir")
	})

	t.Run("MountDir is passed as a \".\" string", func(t *testing.T) {
		workDirCfg, err := IsWorkDirValid(".")
		assert.NoError(t, err, "The isWorkDirValid should not return an error")

		mountDirCfg, err := IsMountOrTargetDirValid(".", workDirCfg, "mount")

		assert.NoError(t, err, "The IsMountOrTargetDirValid should not return an error")
		assert.Equal(t, ".", mountDirCfg.Dir, "The mounted directory should be the same as the workdir")
		assert.Equal(t, currentDir, mountDirCfg.Path, "The mounted directory should be the same as the workdir")
	})

	// The mount dir is passed, and it's a relative path.
	t.Run("MountDir is passed as a relative path", func(t *testing.T) {
		newFolderName := "newfolder"
		err := os.Mkdir(newFolderName, 0755)

		assert.NoError(t, err, "The isWorkDirValid should not return an error")
		mountDirPath := newFolderName

		workDirCfg, err := IsWorkDirValid(".")
		assert.NoError(t, err, "The isWorkDirValid should not return an error")

		mountDirCfg, err := IsMountOrTargetDirValid(mountDirPath, workDirCfg, "mount")

		assert.NoError(t, err, "The IsMountOrTargetDirValid should not return an error")
		assert.NotEmpty(t, mountDirCfg.Dir, "The mounted directory should be the same as the workdir")
		assert.NotEmpty(t, mountDirCfg.Path, "The mounted directory should be the same as the workdir")
		assert.Equal(t, mountDirCfg.Dir, newFolderName, "The mounted directory should be the same as the workdir")
		assert.Equal(t, mountDirCfg.Path, currentDir+"/"+newFolderName, "The mounted directory should be the same as the workdir")

		// Remove the folder
		err = os.Remove(newFolderName)
		assert.NoError(t, err, "The isWorkDirValid should not return an error")
	})

	// Mount dir is passed as an absolute path
	t.Run("MountDir is passed as an absolute path", func(t *testing.T) {
		absolutePathFromCurrentDirTwoLevelsBelow := "/../../"
		workDir := filepath.Join(currentDir, absolutePathFromCurrentDirTwoLevelsBelow)
		workDirAbs, err := filepath.Abs(workDir)

		workDirCfg, err := IsWorkDirValid(workDirAbs)
		assert.NoError(t, err, "The isWorkDirValid should not return an error")

		mountDirPath := "pkg/pipeline"
		mountDirPathAbs := filepath.Join(workDirCfg.Path, mountDirPath)
		mountDirCfg, err := IsMountOrTargetDirValid(mountDirPathAbs, workDirCfg, "mount")

		assert.NoError(t, err, "The IsMountOrTargetDirValid should not return an error")
		assert.NotEmpty(t, mountDirCfg.Dir, "The mounted directory should be the same as the workdir")
		assert.NotEmpty(t, mountDirCfg.Path, "The mounted directory should be the same as the workdir")
		assert.Equal(t, mountDirCfg.Path, mountDirPathAbs,
			"The mounted directory should be the same as the workdir")
	})
}
