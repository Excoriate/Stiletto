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
	ecrRepositoryName string
	ecrRegistryName   string
	imageTag          string
	dockerFileName    string
)

var ECRCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "ecr",
	Long: `The 'ecr' command automates and implement actions on top of AWS Elastic Container
Registry`,
	Example: `
  # Push an image into ECR:
  stiletto aws ecr --task=push`,
	Run: func(cmd *cobra.Command, args []string) {

		msg := tui.NewTUIMessage()
		ux := tui.TUITitle{}

		stackName := "AWS"
		jobName := "ECR"

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

		err = task.RunTaskAWSECR(task.InitOptions{
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

func addECRCmdFlags() {
	ECRCmd.Flags().StringVarP(&ecrRepositoryName, "ecr-repository", "", "",
		"The name of the ECR repository")
	ECRCmd.Flags().StringVarP(&imageTag, "tag", "", "latest",
		"The tag of the image to be pushed. If not specified, it will default to 'latest'")
	ECRCmd.Flags().StringVarP(&dockerFileName, "dockerfile", "", "",
		"The name of the Dockerfile. If not specified, it will default to 'Dockerfile'")
	ECRCmd.Flags().StringVarP(&ecrRegistryName, "ecr-registry", "", "",
		"The name of the ECR registry.")

	err := ECRCmd.MarkFlagRequired("ecr-repository")
	if err != nil {
		panic(err)
	}

	err = ECRCmd.MarkFlagRequired("ecr-registry")
	if err != nil {
		panic(err)
	}

	_ = viper.BindPFlag("ecr-repository", ECRCmd.Flags().Lookup("ecr-repository"))
	_ = viper.BindPFlag("ecr-registry", ECRCmd.Flags().Lookup("ecr-registry"))
	_ = viper.BindPFlag("tag", ECRCmd.Flags().Lookup("tag"))
	_ = viper.BindPFlag("dockerfile", ECRCmd.Flags().Lookup("dockerfile"))
}

func init() {
	addECRCmdFlags()
}
