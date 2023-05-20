package pipeline

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/filesystem"
	"github.com/Excoriate/stiletto/pkg/config"
	"path/filepath"
)

func IsWorkDirValid(workDir string) (config.PipelineDirs, error) {
	workDirCfg := config.PipelineDirs{}

	// If the workDir is empty, we will use the current directory.
	if workDir == "" || workDir == "." {
		workDirCfg.Dir = "."
		workDirCfg.Path = config.GetDefaultDirs().CurrentDir
		return workDirCfg, nil
	}

	currentDir := config.GetDefaultDirs().CurrentDir

	// Scenario 1: workdir is passed, but as a relative path. Check if it's relative
	if !filepath.IsAbs(workDir) {
		workDirPath := filepath.Join(currentDir, workDir)

		workDirAbs, err := filesystem.PathToAbsolute(common.NormaliseNoSpaces(workDirPath))
		if err != nil {
			errMsg := fmt.Sprintf("%s: invalid working directory: %s. "+
				"Cant convert it to an absolute path. ", pipelineMsgPrefix, workDirPath)

			return workDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
		}

		if err := filesystem.DirIsValid(workDirAbs); err != nil {
			errMsg := fmt.Sprintf("%s: invalid working directory: %s. "+
				"It does not exist or it is an invalid directory", pipelineMsgPrefix, workDirAbs)

			return workDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
		}

		workDirCfg.Dir = workDir
		workDirCfg.Path = workDirAbs

		return workDirCfg, nil
	}

	// Scenario 2: workdir is passed, and it's an absolute path.
	var workDirAbs string
	if err := filesystem.IsPathAbsolute(workDir); err != nil {
		workDirAbs, err = filesystem.PathToAbsolute(common.NormaliseNoSpaces(workDir))
		if err != nil {
			errMsg := fmt.Sprintf("%s: invalid working directory: %s. "+
				"Cant convert it to an absolute path. ", pipelineMsgPrefix, workDir)

			return workDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
		}
	} else {
		workDirAbs = common.NormaliseNoSpaces(workDir)
	}

	if err := filesystem.DirIsValid(workDirAbs); err != nil {
		errMsg := fmt.Sprintf("%s: invalid working directory: %s. "+
			"It does not exist or it is an invalid directory", pipelineMsgPrefix, workDir)

		return workDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
	}

	workDirCfg.Dir = workDir
	workDirCfg.Path = workDirAbs

	return workDirCfg, nil
}

func IsMountOrTargetDirValid(childDir string, workDirCfg config.PipelineDirs,
	typeOfDir string) (config.PipelineDirs, error) {
	childDirCfg := config.PipelineDirs{}

	// If the childDir is empty or ".", default to workDir configuration.
	if childDir == "" {
		childDirCfg.Dir = workDirCfg.Dir
		childDirCfg.Path = workDirCfg.Path
		return childDirCfg, nil
	}

	if childDir == "." {
		childDirCfg.Dir = "."
		childDirCfg.Path = config.GetDefaultDirs().CurrentDir
		return childDirCfg, nil
	}

	// Scenario 1: childDir is passed as a relative path. Join it with the workDir path and check its validity.
	if !filepath.IsAbs(childDir) {
		childDirPath := filepath.Join(workDirCfg.Path, childDir)
		childDirAbs, err := filesystem.PathToAbsolute(common.NormaliseNoSpaces(
			childDirPath))

		if err != nil {
			errMsg := fmt.Sprintf("%s: invalid %s directory: %s. "+
				"Can't convert it to an absolute path. ", pipelineMsgPrefix, typeOfDir,
				childDirPath)

			return childDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
		}

		if err := filesystem.DirIsValid(childDirAbs); err != nil {
			errMsg := fmt.Sprintf("%s: invalid %s directory: %s. "+
				"It does not exist or it is an invalid directory", pipelineMsgPrefix,
				typeOfDir, childDirAbs)

			return childDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
		}

		if err := filesystem.IsSubDir(workDirCfg.Path, childDir); err != nil {
			errMsg := fmt.Sprintf("%s: invalid %s directory: %s. "+
				"It is not a subdirectory of working directory %s", pipelineMsgPrefix,
				typeOfDir, childDirAbs, workDirCfg.Path)

			return childDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
		}

		childDirCfg.Dir = childDir
		childDirCfg.Path = childDirAbs

		return childDirCfg, nil
	}

	// Scenario 2: childDir is passed as an absolute path. Check its validity.
	var childDirAbs string
	if err := filesystem.IsPathAbsolute(childDir); err != nil {
		childDirAbs, err = filesystem.PathToAbsolute(common.NormaliseNoSpaces(childDir))
		if err != nil {
			errMsg := fmt.Sprintf("%s: invalid %s directory: %s. "+
				"Can't convert it to an absolute path. ", pipelineMsgPrefix, typeOfDir, childDir)

			return childDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
		}
	} else {
		childDirAbs = common.NormaliseNoSpaces(childDir)
	}

	if err := filesystem.DirIsValid(childDirAbs); err != nil {
		errMsg := fmt.Sprintf("%s: invalid %s directory: %s. "+
			"It does not exist or it is an invalid directory", pipelineMsgPrefix, typeOfDir,
			childDirAbs)

		return childDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
	}

	if err := filesystem.IsSubDir(workDirCfg.Path, childDirAbs); err != nil {
		errMsg := fmt.Sprintf("%s: invalid mount directory: %s. "+
			"It is not a subdirectory of working directory %s", pipelineMsgPrefix, childDirAbs, workDirCfg.Path)

		return childDirCfg, errors.NewPipelineConfigurationError(errMsg, err)
	}

	childDirCfg.Dir = childDir
	childDirCfg.Path = childDirAbs

	return childDirCfg, nil
}
