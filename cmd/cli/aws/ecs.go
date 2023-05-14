package aws

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/api"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/config"
	"github.com/Excoriate/stiletto/pkg/task"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	ecsService     string
	ecsCluster     string
	taskDefinition string // Map to the task definition family name.
	imageURL       string
	releaseVersion string
)

var ECSCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "ecs",
	Long: `The 'ecs' command automates and implement several Elastic Container Service actions,
E.g.: 'deploy'`,
	Example: `
  # Deploy a new version of a task running in a ECS service:
  stiletto aws ecs --task=deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Instantiate the pipeline runner, which will be used to run the tasks.
		msg := tui.NewTUIMessage()
		ux := tui.TUITitle{}

		stackName := "AWS"
		jobName := "ECS"

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

		err = task.RunTaskAWSECS(task.InitOptions{
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

func addECSCmdFlags() {
	ECSCmd.Flags().StringVarP(&ecsService, "ecs-service", "", "",
		"The name of the ECS service to be deployed.")

	ECSCmd.Flags().StringVarP(&ecsCluster, "ecs-cluster", "", "",
		"The name of the ECS cluster to be deployed.")

	ECSCmd.Flags().StringVarP(&taskDefinition, "task-definition", "", "",
		"The name of the ECS task definition to be deployed.")

	ECSCmd.Flags().StringVarP(&imageURL, "image-url", "", "",
		"The URL of the image to be deployed.")

	ECSCmd.Flags().StringVarP(&releaseVersion, "release-version", "", "",
		"The tag or version of the (container) image to be deployed. If not specified, "+
			"the default value is 'latest'.")

	err := ECSCmd.MarkFlagRequired("ecs-service")
	if err != nil {
		panic(err)
	}

	err = ECSCmd.MarkFlagRequired("ecs-cluster")
	if err != nil {
		panic(err)
	}

	err = ECSCmd.MarkFlagRequired("task-definition")
	if err != nil {
		panic(err)
	}

	_ = viper.BindPFlag("ecs-service", ECSCmd.Flags().Lookup("ecs-service"))
	_ = viper.BindPFlag("ecs-cluster", ECSCmd.Flags().Lookup("ecs-cluster"))
	_ = viper.BindPFlag("task-definition", ECSCmd.Flags().Lookup("task-definition"))
	_ = viper.BindPFlag("image-url", ECSCmd.Flags().Lookup("image-url"))
	_ = viper.BindPFlag("release-version", ECSCmd.Flags().Lookup("release-version"))
}

func init() {
	addECSCmdFlags()
}
