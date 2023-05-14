package task

import (
	"context"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/filesystem"
	"github.com/Excoriate/stiletto/pkg/config"
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
	// Check if module exist.
	if err := filesystem.DirIsNotEmpty(tgModuleDirValue); err != nil {
		return InfraTerraGruntActionArgs{}, errors.NewActionCfgError("Failed to run this Terragrunt action, "+
			"the target module dir is empty.", err)
	}

	args.TargetModuleDir = tgModuleDirValue

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

func (a *InfraTerraGruntAction) Plan() (Output, error) {
	_, _ = a.GetOptions()
	// Getting all the requirements.
	uxLog := a.Task.GetPipelineUXLog()
	opts, err := a.GetOptions()

	if err != nil {
		errMsg := "Failed to run action: 'Plan' - Cannot pass the 'action' arguments validations"
		uxLog.ShowError(a.prefix, errMsg, err)
		return Output{}, errors.NewActionCfgError(errMsg, err)
	}

	container := a.Task.GetJobContainerDefault()
	client := a.Task.GetClient()
	ctx := a.Task.GetJob().Ctx
	preRequiredFiles := []string{opts.TgConfigFile}

	mountedContainer, mntErr := a.Task.MountDir(opts.TargetModuleDir, client, container,
		preRequiredFiles, ctx)

	if mntErr != nil {
		return Output{}, mntErr
	}

	cmds := [][]string{opts.Commands}

	if err := a.Task.RunCmdInContainer(mountedContainer, cmds, false, ctx); err != nil {
		return Output{}, err
	}

	return Output{}, nil
}

func NewInfraTerraGruntAction(task CoreTasker, prefix string) *InfraTerraGruntAction {
	return &InfraTerraGruntAction{
		Task:   task,
		prefix: prefix,
		Id:     common.GetUUID(),
		Name:   "Perform infrastructure changes using Terragrunt",
	}
}
