package pipeline

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/filesystem"
	"github.com/Excoriate/stiletto/internal/logger"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/config"
)

func isWorkDirValid(pLog logger.Logger, workDir string) error {
	if workDir == "" {
		return nil
	}

	workDirNormalised := common.NormaliseNoSpaces(workDir)
	if _, err := filesystem.PathExist(workDirNormalised); err != nil {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, invalid working directory: %s. It does not exist",
			workDirNormalised)

		pLog.LogError("init", errMsg)

		return errors.NewPipelineConfigurationError(errMsg, err)
	}

	if err := filesystem.PathIsADirectory(workDirNormalised); err != nil {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, invalid working directory: %s it is not a directory", workDirNormalised)
		return errors.NewPipelineConfigurationError(errMsg, err)
	}

	return nil
}

func isMountDirValid(mountDir string, workDir string) error {
	mountDirNormalised := common.NormaliseNoSpaces(mountDir)

	if mountDirNormalised == "" {
		return nil // it's not passed, it's fine. It'll be set to the working directory
	}

	if _, err := filesystem.PathExist(mountDirNormalised); err != nil {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, invalid mount directory: %s. It does not exist", mountDirNormalised)
		return errors.NewPipelineConfigurationError(errMsg, err)
	}

	if err := filesystem.PathIsADirectory(mountDirNormalised); err != nil {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, invalid mount directory: %s it is not a directory", mountDirNormalised)
		return errors.NewPipelineConfigurationError(errMsg, err)
	}

	// The mountDir if passed, should be a subdirectory of the working directory
	if err := filesystem.IsSubDir(workDir, mountDir); err != nil {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, mount directory: %s is not a subdirectory of working directory: %s", mountDir, workDir)
		return errors.NewPipelineConfigurationError(errMsg, err)
	}

	return nil
}

func isTargetDirValid(targetDir string, mountDir string, workDir string) error {
	if targetDir == "" {
		return nil // it's not passed, it's fine. It'll be set to the mount directory
	}

	targetDirNormalised := common.NormaliseNoSpaces(targetDir)
	if _, err := filesystem.PathExist(targetDirNormalised); err != nil {
		return errors.NewPipelineConfigurationError("PipelineCfg cant initialise", err)
	}

	if err := filesystem.PathIsADirectory(targetDirNormalised); err != nil {
		return errors.NewPipelineConfigurationError("PipelineCfg cant initialise", err)
	}

	// The targetDir if passed, should be a subdirectory of the mount directory
	if err := filesystem.IsSubDir(mountDir, targetDir); err != nil {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, target directory: %s is not a subdirectory of mount directory: %s", targetDir, mountDir)
		return errors.NewPipelineConfigurationError(errMsg, err)
	}

	return nil
}

func isTaskNameValid(taskName string) error {
	normalisedTask := common.NormaliseStringUpper(taskName)

	// FIXME: Looks like it's redundant. Normally,
	// this parameter is validated from the UX/CLI level.
	if normalisedTask == "" {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, invalid task name: %s", taskName)
		return errors.NewPipelineConfigurationError(errMsg, nil)
	}

	return nil
}

func areEnvKeysToScanExported(envKeysToScan []string) error {
	if len(envKeysToScan) == 0 {
		return nil
	}

	err := filesystem.AreEnvVarsExportedAndSet(envKeysToScan)
	if err != nil {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, These keys are not exported: %s", envKeysToScan)
		return errors.NewPipelineConfigurationError(errMsg, err)
	}
	return nil
}

func isEnvVarsMapToSetValid(envVarsMapToSet map[string]string) error {
	if len(envVarsMapToSet) == 0 {
		return nil
	}

	for key := range envVarsMapToSet {
		if _, ok := envVarsMapToSet[key]; !ok {
			return errors.NewPipelineConfigurationError("PipelineCfg cant initialise", fmt.Errorf("env var %s is not set", key))
		}
	}

	return nil
}

func isAWSKeysExported(isAWSKeysToScan bool) error {
	if !isAWSKeysToScan {
		return nil
	}
	if _, err := filesystem.ScanAWSCredentialsEnvVars(); err != nil {
		return errors.NewPipelineConfigurationError("PipelineCfg cant initialise", err)
	}

	return nil
}

func isTFEnvVarsExported(isTFEnvVarsToScan bool) error {
	if !isTFEnvVarsToScan {
		return nil
	}

	if _, err := filesystem.ScanTerraformEnvVars(); err != nil {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, " +
			"There is no TF_VAR or related terraform environment variables exported. ")
		return errors.NewPipelineConfigurationError(errMsg, err)
	}

	return nil
}

func CheckPreConditions(args *config.PipelineOptions, pLog logger.Logger) error {
	ux := tui.TUIMessage{}

	// 1. Validate the working directory.
	if err := isWorkDirValid(pLog, args.WorkDir); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if args.WorkDir == "" {
		args.WorkDir = "."
	}

	workDirAbsolute, _ := filesystem.PathToAbsolute(args.WorkDir)
	args.WorkDirPath = workDirAbsolute

	// 2. Validate the mount directory.
	if mountDirErr := isMountDirValid(args.MountDir, args.WorkDir); mountDirErr != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", mountDirErr)
		return mountDirErr
	}

	if args.MountDir == "" {
		args.MountDir = "."
	}

	args.MountDirPath, _ = filesystem.PathToAbsolute(args.MountDir)

	// 3. Validate the target directory.
	if err := isTargetDirValid(args.TargetDir, args.MountDir, args.WorkDir); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if args.TargetDir == "" {
		args.TargetDir = args.MountDir
	}

	args.TargetDirPath, _ = filesystem.PathToAbsolute(args.TargetDir)

	if err := isTaskNameValid(args.TaskName); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := areEnvKeysToScanExported(args.EnvVarsToScanAndSet); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := isEnvVarsMapToSetValid(args.EnvKeyValuePairsToSet); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := isAWSKeysExported(args.IsAWSEnvVarKeysToScanEnabled); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := isTFEnvVarsExported(args.IsTerraformVarsScanEnabled); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	return nil
}

func New(workDir, mountDir, targetDir, taskName string, envVarKeysToScan []string,
	envVarsMapToSet map[string]string, isAWSKeysToScan bool, isTFScanEnabled bool,
	initDaggerWithWorkDirByDefault bool) (*Config,
	error) {

	logPrinter := logger.NewLogger()
	logPrinter.InitLogger()

	args := config.PipelineOptions{
		// Key directories
		WorkDir:   common.NormaliseNoSpaces(workDir),
		MountDir:  common.NormaliseNoSpaces(mountDir),
		TargetDir: common.NormaliseNoSpaces(targetDir),

		// Task identifier, that'll be used to determine what to do.
		TaskName: taskName,
		// Specific environmental options passed.
		EnvVarsToScanAndSet:   envVarKeysToScan,
		EnvKeyValuePairsToSet: envVarsMapToSet,
		EnvVarsAWSKeysToScan:  map[string]string{},
		// Scan options
		IsAWSEnvVarKeysToScanEnabled:   isAWSKeysToScan,
		IsTerraformVarsScanEnabled:     isTFScanEnabled,
		InitDaggerWithWorkDirByDefault: initDaggerWithWorkDirByDefault,
	}

	if err := CheckPreConditions(&args, logPrinter); err != nil {
		return nil, err
	}

	dirs := config.GetDefaultDirs()

	platformToArch := map[dagger.Platform]string{
		"linux/amd64": "amd64",
		"linux/arm64": "arm64",
	}

	return &Config{
		Logger:       logPrinter,
		Dirs:         *dirs,
		UXDisplay:    tui.NewTitle(),
		Platforms:    platformToArch,
		UXMessage:    tui.NewTUIMessage(),
		PipelineOpts: &args,
		Ctx:          context.Background(),
	}, nil
}
