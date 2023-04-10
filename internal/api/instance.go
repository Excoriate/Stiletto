package api

import (
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/config"
	"github.com/Excoriate/stiletto/pkg/job"
	"github.com/Excoriate/stiletto/pkg/pipeline"
)

func New(cliArgs *config.CLIGlobalArgs, stack, jobName string) (*pipeline.Config, *job.Job, error) {
	msg := tui.NewTUIMessage()
	ux := tui.TUITitle{}

	stackNormalised := common.NormaliseStringUpper(stack)
	jobNormalised := common.NormaliseStringUpper(jobName)

	config.ShowCLITitle()

	// Pipeline instance.
	p, err := pipeline.New(cliArgs.WorkingDir, cliArgs.MountDir,
		cliArgs.TargetDir, cliArgs.TaskName,
		cliArgs.ScanEnvVarKeys,
		cliArgs.EnvKeyValuePairsToSetString, cliArgs.ScanAWSKeys,
		cliArgs.ScanTerraformVars, cliArgs.InitDaggerWithWorkDirByDefault)

	if err != nil {
		msg.ShowError("INIT", "Failed pipeline initialization", err)
		return nil, nil, err
	}

	ux.ShowSubTitle(stackNormalised, jobNormalised)
	ux.ShowInitDetails(jobNormalised, cliArgs.TaskName, p.PipelineOpts.WorkDirPath,
		p.PipelineOpts.TargetDirPath, p.PipelineOpts.MountDirPath)

	// 2. Initialising the job.
	j, err := job.NewJob(p, job.InitOptions{
		Name:  cliArgs.TaskName,
		Stack: "AWS",

		// Pipeline reference.
		PipelineCfg: p,

		// Critical directories to be resolved.
		WorkDir:   p.PipelineOpts.WorkDir,
		TargetDir: p.PipelineOpts.TargetDir,
		MountDir:  p.PipelineOpts.MountDir,

		// Environmental configuration
		ScanAWSEnvVars:       cliArgs.ScanAWSKeys,
		ScanTerraformEnvVars: cliArgs.ScanTerraformVars,
		EnvVarsToSet:         cliArgs.EnvKeyValuePairsToSetString,
		EnvVarsToScan:        cliArgs.ScanEnvVarKeys,
	})

	if err != nil {
		msg.ShowError("INIT", "Failed job initialization", err)
		return nil, nil, err
	}

	return p, j, nil
}
