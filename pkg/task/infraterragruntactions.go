package task

import (
	"context"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/filesystem"
	"github.com/Excoriate/stiletto/pkg/config"
	"path/filepath"
)

type InfraTerraGruntAction struct {
	Task   CoreTasker
	prefix string // How the UX messages should be prefixed
	Id     string // The ID of the task
	Name   string // The name of the task
	Ctx    context.Context
}

type InfraTerraGruntActionOptions struct {
	AWSAccessKey  string
	AWSSecretKey  string
	AWSRegion     string
	CommandsToRun []string // Either standard commands, or custom commands passed if applicable.
}

type InfraTerraGruntActionArgs struct {
	AWSRegion       string
	AWSAccessKey    string
	AWSSecretKey    string
	TargetModuleDir string
	Commands        []string
	TgConfigFile    string
}

type TerrGruntActionRunner interface {
	GetOptions() (InfraTerraGruntActionArgs, error)
	Plan() (Output, error)
	Apply() (Output, error)
	Destroy() (Output, error)
	Validate() (Output, error)
	RunTGCommand(commands [][]string) (Output, error)
}

func (a *InfraTerraGruntAction) GetOptions() (InfraTerraGruntActionArgs, error) {
	ux := a.Task.GetPipelineUXLog()
	var args InfraTerraGruntActionArgs

	isAWSScanOptEnabled := a.Task.GetPipeline().PipelineOpts.IsAWSEnvVarKeysToScanEnabled
	if isAWSScanOptEnabled {

		ux.ShowWarning(a.prefix,
			"This task is configured to scan for AWS credentials that were scanned from the"+
				" environment variables. ")

		awsCred := a.Task.GetJob().EnvVarsAWSScanned

		args = InfraTerraGruntActionArgs{
			AWSRegion:    awsCred["AWS_REGION"],
			AWSAccessKey: awsCred["AWS_ACCESS_KEY_ID"],
			AWSSecretKey: awsCred["AWS_SECRET_ACCESS_KEY"],
		}
	}

	// If module dir is not set, then it'll fail.
	cfg := config.Cfg{}
	tgModuleDirCfg, err := cfg.GetFromViper("target-module")
	if err != nil {
		return InfraTerraGruntActionArgs{}, errors.NewActionCfgError("Failed to run this Terragrunt action, "+
			"it cannot find the target module dir.", err)
	}

	tgModuleDirValue := tgModuleDirCfg.Value.(string)
	workDirPath := a.Task.GetPipeline().PipelineOpts.WorkDirPath

	// The target module dir should be a relative of the working directory
	if err := filesystem.IsRelativeChildPath(workDirPath, tgModuleDirValue); err != nil {
		return InfraTerraGruntActionArgs{}, errors.NewActionCfgError(
			"Failed to validate the target module directory. "+
				"The target module passed %s is not a child of the working directory %s. "+
				"TerraGrunt require a proper git repository to resolve paths", err)
	}

	// The target module dir should be a valid git repository
	targetModuleDirWithWorkDir := filepath.Join(workDirPath, tgModuleDirValue)
	if _, err := filesystem.IsGitRepository(targetModuleDirWithWorkDir, false, true); err != nil {
		return InfraTerraGruntActionArgs{}, errors.NewActionCfgError(
			"Failed to validate the target module directory. "+
				"The target module passed %s is not a git repository. "+
				"TerraGrunt require a proper git repository to resolve paths", err)
	}

	args.TargetModuleDir = tgModuleDirValue
	ux.ShowInfo(a.prefix, "The target module dir is: "+workDirPath)

	// If commands are passed, validate them and use them
	tgCmds, err := cfg.GetFromViperOrDefault("tg-commands", []string{})
	if err != nil {
		return InfraTerraGruntActionArgs{}, errors.NewActionCfgError("Failed to run this Terragrunt action, "+
			"it cannot fetch or retrieve the 'tg-commands' configuration.", err)
	}

	tgCommands := tgCmds.Value.([]string)

	if err := common.ValidateTerragruntCommands(tgCommands); err != nil {
		return InfraTerraGruntActionArgs{}, errors.NewActionCfgError("Failed to run this Terragrunt action, "+
			"the 'tg-commands' configuration is invalid.", err)
	}

	// If commands are passed, use them always adding the 'terragrunt' command at first
	tgCommands = append([]string{"terragrunt"}, tgCommands...)
	args.Commands = tgCommands
	args.TgConfigFile = "terragrunt.hcl"

	return args, nil
}

func (a *InfraTerraGruntAction) RunTGCommand(commands [][]string) (Output, error) {
	uxLog := a.Task.GetPipelineUXLog()

	// Fetch action's configuration
	opts, err := a.GetOptions()

	if err != nil {
		errMsg := "Failed to run action: 'RunTGCommand' - Cannot pass the 'action' arguments validations"
		uxLog.ShowError(a.prefix, errMsg, err)
		return Output{}, errors.NewActionCfgError(errMsg, err)
	}

	// Reference required objects (container, client, context, etc.)
	container := a.Task.GetJobContainerDefault()
	client := a.Task.GetClient()
	ctx := a.Task.GetJob().Ctx
	preRequiredFiles := []string{opts.TgConfigFile}

	// Inherit the environment variables from the job.
	preConfiguredContainer, preCfgErr := a.Task.SetEnvVarsFromJob(container)
	if preCfgErr != nil {
		uxLog.ShowError(a.prefix, "Failed to run action: 'RunTGCommand' - Cannot set the environment variables from the job", preCfgErr)
		return Output{}, errors.NewActionCfgError("Failed to run action: 'RunTGCommand' - Cannot set the environment variables from the job", preCfgErr)
	}

	// Mount required directories.
	workDirPath := a.Task.GetPipeline().PipelineOpts.WorkDirPath
	configuredContainer, mntErr := a.Task.MountDir(workDirPath, opts.TargetModuleDir, client,
		preConfiguredContainer, preRequiredFiles, ctx)

	if mntErr != nil {
		return Output{}, mntErr
	}

	// Run the commands.
	var cmdsToRun [][]string
	if len(commands) > 0 {
		cmdsToRun = commands
	} else {
		cmdsToRun = [][]string{opts.Commands}
	}

	if err := a.Task.RunCmdInContainer(configuredContainer, cmdsToRun, false, ctx); err != nil {
		return Output{}, err
	}

	return Output{}, nil
}

func (a *InfraTerraGruntAction) Plan() (Output, error) {
	inspectCfgFile := []string{"cat", "terragrunt.hcl"}
	planCmd := []string{"terragrunt", "plan"}

	cmds := [][]string{inspectCfgFile, planCmd}
	return a.RunTGCommand(cmds)
}

func (a *InfraTerraGruntAction) Apply() (Output, error) {
	inspectCfgFile := []string{"cat", "terragrunt.hcl"}
	applyCmd := []string{"terragrunt", "apply", "-auto-approve"}

	cmds := [][]string{inspectCfgFile, applyCmd}
	return a.RunTGCommand(cmds)
}

func (a *InfraTerraGruntAction) Destroy() (Output, error) {
	inspectCfgFile := []string{"cat", "terragrunt.hcl"}
	destroyCmd := []string{"terragrunt", "destroy", "-auto-approve"}

	cmds := [][]string{inspectCfgFile, destroyCmd}
	return a.RunTGCommand(cmds)
}

func (a *InfraTerraGruntAction) Validate() (Output, error) {
	inspectCfgFile := []string{"cat", "terragrunt.hcl"}
	validateCmd := []string{"terragrunt", "validate"}

	cmds := [][]string{inspectCfgFile, validateCmd}
	return a.RunTGCommand(cmds)
}

func NewInfraTerraGruntAction(task CoreTasker, prefix string) *InfraTerraGruntAction {
	return &InfraTerraGruntAction{
		Task:   task,
		prefix: prefix,
		Id:     common.GetUUID(),
		Name:   "Perform infrastructure changes using Terragrunt",
	}
}
