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

var pipelineMsgPrefix = "Pipeline validation"

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

func isDotEnvFileValidToScan(isScanFromDotEnvFileEnabled bool, dotEnvFilePath string) (map[string]string, error) {
	if isScanFromDotEnvFileEnabled {
		envVars, err := filesystem.GetEnvVarsFromDotFile(dotEnvFilePath)
		if err != nil {
			errMsg := fmt.Sprintf("PipelineCfg cant initialise, "+
				"there is no .env file in the working directory: %s", dotEnvFilePath)
			return nil, errors.NewPipelineConfigurationError(errMsg, err)
		}

		if len(envVars) == 0 {
			errMsg := fmt.Sprintf("PipelineCfg cant initialise, "+
				"there is no env vars in the .env file: %s", dotEnvFilePath)
			return nil, errors.NewPipelineConfigurationError(errMsg, err)
		}

		return envVars, nil
	}

	return nil, nil
}

func isEnvVarsToScanFromPrefixValid(isScanFromPrefixEnabled bool, envVarsToScanFromPrefix []string) error {
	if !isScanFromPrefixEnabled {
		return nil
	}

	if len(envVarsToScanFromPrefix) == 0 {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, "+
			"there is no env vars to scan from prefix: %s", envVarsToScanFromPrefix)
		return errors.NewPipelineConfigurationError(errMsg, nil)
	}

	envVarsTo, err := filesystem.ScanEnvVarsFromPrefixes(envVarsToScanFromPrefix)
	if err != nil {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, "+
			"there is no env vars to scan from prefix: %s", envVarsToScanFromPrefix)
		return errors.NewPipelineConfigurationError(errMsg, err)
	}

	if len(envVarsTo) == 0 {
		errMsg := fmt.Sprintf("PipelineCfg cant initialise, "+
			"there is no env vars to scan from prefix: %s", envVarsToScanFromPrefix)
		return errors.NewPipelineConfigurationError(errMsg, err)
	}

	return nil
}

func CheckPreConditions(args *config.PipelineOptions, pLog logger.Logger) error {
	ux := tui.TUIMessage{}

	// 1. Validate the working directory.
	workDirCfg, err := IsWorkDirValid(args.WorkDir)
	if err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	args.WorkDir = workDirCfg.Dir
	args.WorkDirPath = workDirCfg.Path

	// 2. Validate the mount directory.
	mountDirCfg, err := IsMountOrTargetDirValid(args.MountDir, workDirCfg, "mount")
	if err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	args.MountDir = mountDirCfg.Dir
	args.MountDirPath = mountDirCfg.Path

	// 3. Validate the target directory.
	if args.TargetDir == "" {
		args.TargetDir = args.MountDir
		args.TargetDirPath = args.MountDirPath
	} else {
		targetDirCfg, err := IsMountOrTargetDirValid(args.TargetDir, workDirCfg, "target")
		if err != nil {
			ux.ShowError("VALIDATION", "Preconditions failed", err)
			return err
		}
		args.TargetDir = targetDirCfg.Dir
		args.TargetDirPath = targetDirCfg.Path
	}

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

	if _, err := isDotEnvFileValidToScan(args.IsEnvVarsToScanFromDotEnvFile,
		args.EnvVarsDotEnvFilePath); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	if err := isEnvVarsToScanFromPrefixValid(args.IsEnvVarsToScanByPrefix,
		args.EnvVarsToScanByPrefix); err != nil {
		ux.ShowError("VALIDATION", "Preconditions failed", err)
		return err
	}

	return nil
}

func New(workDir, mountDir, targetDir, taskName string, envVarKeysToScan []string,
	envVarsMapToSet map[string]string, isAWSKeysToScan bool, isTFScanEnabled bool,
	isAllEnvVarsToScan bool, dotEnvFile string, envVarsToScanByPrefix []string,
	initDaggerWithWorkDirByDefault bool) (*Config,
	error) {

	logPrinter := logger.NewLogger()
	logPrinter.InitLogger()

	var isEnvVarsToScanFromDotEnvFile bool
	if dotEnvFile == "" {
		isEnvVarsToScanFromDotEnvFile = false
	} else {
		isEnvVarsToScanFromDotEnvFile = true
	}

	var isEnvVarsToScanFromPrefix bool
	if len(envVarsToScanByPrefix) == 0 {
		isEnvVarsToScanFromPrefix = false
	} else {
		isEnvVarsToScanFromPrefix = true
	}

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
		EnvVarsDotEnvFilePath: dotEnvFile,
		EnvVarsAWSKeysToScan:  map[string]string{},
		EnvVarsToScanByPrefix: envVarsToScanByPrefix,
		EnvVarsFromDotEnvFile: map[string]string{},
		// Scan options
		IsAWSEnvVarKeysToScanEnabled:   isAWSKeysToScan,
		IsTerraformVarsScanEnabled:     isTFScanEnabled,
		InitDaggerWithWorkDirByDefault: initDaggerWithWorkDirByDefault,
		IsEnvVarsToScanFromDotEnvFile:  isEnvVarsToScanFromDotEnvFile,
		IsEnvVarsToScanByPrefix:        isEnvVarsToScanFromPrefix, // Scan env vars by prefix.
		IsAllEnvVarsToScanEnabled:      isAllEnvVarsToScan,
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
