package docker

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/api"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/config"
	"github.com/Excoriate/stiletto/pkg/task"
	"github.com/spf13/cobra"
	"os"
)

var Cmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "docker",
	Long: `The 'docker' command automates various Continuous Integration (Docker-related) tasks.
You can specify the tasks you want to perform using the provided --task flag.`,
	Example: `
  # Build a docker image from an existing DockerFile:
  stiletto docker --task=build`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Instantiate the pipeline runner, which will be used to run the tasks.
		ux := tui.TUITitle{}
		msg := tui.NewTUIMessage()

		stackName := "DOCKER"
		jobName := "BUILD"

		cliGlobalArgs, err := config.GetCLIGlobalArgs()

		if err != nil {
			panic(err)
		}

		p, j, err := api.New(&cliGlobalArgs, stackName, jobName)
		if err != nil {
			panic(err)
		}

		ux.ShowSubTitle("TASK:", cliGlobalArgs.TaskName)
		ux.ShowTaskDetails(jobName, cliGlobalArgs.TaskName, j.WorkDirPath,
			j.TargetDirPath,
			j.MountDirPath)

		err = task.RunTaskDocker(task.InitOptions{
			//Task:           GlobalTaskName,
			Task:           cliGlobalArgs.TaskName,
			Stack:          stackName,
			PipelineCfg:    p,
			JobCfg:         j,
			WorkDir:        p.PipelineOpts.WorkDir,
			MountDir:       p.PipelineOpts.MountDir,
			TargetDir:      p.PipelineOpts.TargetDir,
			ActionCommands: cliGlobalArgs.CustomCommands,
		})

		if err != nil {
			msg.ShowError("", fmt.Sprintf("Failed to run task '%s' as part of job %s on stack '%s'",
				cliGlobalArgs.TaskName, jobName, stackName), err)
			os.Exit(1)
		}
	},
}

func AddCIArguments() {
	// Add the flags to the root command.
}

func init() {
	AddCIArguments()
}
