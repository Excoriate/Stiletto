package infra

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/api"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/config"
	"github.com/Excoriate/stiletto/pkg/task"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	// Main IaC tools
	tgPlan    bool
	tgPlanAll bool

	tgApply    bool
	tgApplyAll bool

	tgDestroy    bool
	tgDestroyAll bool

	tgCustomCommands []string
	tgModuleDir      string
)

var msg = tui.NewTUIMessage()
var ux = tui.TUITitle{}

var TgCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "terragrunt",
	Long:    `The 'terragrunt' command automate and perform several infra-related actions using Terragrunt.`,
	Example: `
	  stiletto infra terragrunt --plan --target-module=module1
      # Also, the task can be set instead of the flag:
      stiletto infra terragrunt --task=plan --target-module=module1
`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Fail-fast for inconsistent flags
		var tempTgCustomCommands []string
		var argErr error

		if viper.GetStringSlice("tg-commands") != nil {
			tempTgCustomCommands = viper.GetStringSlice("tg-commands")
		}

		tempTaskName := viper.GetString("task")

		if len(tempTgCustomCommands) > 0 && tempTaskName != "" {
			msg.ShowError("", "The task name and custom commands cannot be set at the same time",
				nil)
			argErr = fmt.Errorf("the task name and custom commands cannot be set at the same time")
		}

		if !tgApply && !tgPlan && !tgDestroy && !tgPlanAll && !tgApplyAll && !tgDestroyAll && len(tempTgCustomCommands) == 0 && tempTaskName == "" {
			msg.ShowError("", "No task was set. Please set a task to run.", nil)
			argErr = fmt.Errorf("no task was set. Please set a task to run")
		}

		if tgPlan && tgApply && tgDestroy && tgPlanAll && tgApplyAll && tgDestroyAll {
			msg.ShowError("", "Only one task can be set at a time", nil)
			argErr = fmt.Errorf("only one task can be set at a time")
		}

		if (tgPlan || tgApply || tgDestroy || tgPlanAll || tgApplyAll || tgDestroyAll) && tempTaskName != "" {
			msg.ShowError("", "The task name and task flags cannot be set at the same time", nil)
			argErr = fmt.Errorf("the task name and task flags cannot be set at the same time")
		}

		// Set variable task in 'viper' if it's not set
		if tempTaskName == "" {
			if tgPlan {
				viper.Set("task", "plan")
				viper.Set("tg-commands", []string{"plan"})
			}

			if tgApply {
				viper.Set("task", "apply")
				viper.Set("tg-commands", []string{"apply"})
			}

			if tgDestroy {
				viper.Set("task", "destroy")
				viper.Set("tg-commands", []string{"destroy"})
			}

			if tgPlanAll {
				viper.Set("task", "plan-all")
				viper.Set("tg-commands", []string{"run-all, plan"})
			}

			if tgApplyAll {
				viper.Set("task", "apply-all")
				viper.Set("tg-commands", []string{"run-all, apply"})
			}

			if tgDestroyAll {
				viper.Set("task", "destroy-all")
				viper.Set("tg-commands", []string{"run-all, destroy"})
			}
		} else {
			viper.Set("task", tempTaskName)
			// If the flag is any of the 'run-all' related ones, set the commands accordingly
			if tempTaskName == "plan-all" || tempTaskName == "apply-all" || tempTaskName == "destroy-all" {
				viper.Set("tg-commands", []string{"run-all", tempTaskName})
			} else {
				viper.Set("tg-commands", []string{tempTaskName})
			}
		}

		// If custom commands are set, add the custom commands to the tg-commands slice.
		if len(tempTgCustomCommands) > 0 {
			viper.Set("tg-commands", append(viper.GetStringSlice("tg-commands"), tempTgCustomCommands...))
		}

		return argErr
	},
	Run: func(cmd *cobra.Command, args []string) {
		stackName := "INFRA:TERRAGRUNT"
		jobName := "IAC"

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

		err = task.RunTaskInfraTerraGrunt(task.InitOptions{
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

func addTgFlags() {
	TgCmd.Flags().BoolVarP(&tgPlan, "plan", "", false, "Run a plan.")
	TgCmd.Flags().BoolVarP(&tgPlanAll, "plan-all", "", false, "Run a plan-all.")
	TgCmd.Flags().BoolVarP(&tgApply, "apply", "", false, "Run an apply.")
	TgCmd.Flags().BoolVarP(&tgApplyAll, "apply-all", "", false, "Run an apply-all.")
	TgCmd.Flags().BoolVarP(&tgDestroy, "destroy", "", false, "Run a destroy.")
	TgCmd.Flags().BoolVarP(&tgDestroyAll, "destroy-all", "", false, "Run a destroy-all.")
	TgCmd.Flags().StringSliceVarP(&tgCustomCommands, "tg-commands", "", []string{}, "Run custom commands.")
	TgCmd.Flags().StringVarP(&tgModuleDir, "target-module", "", "", "Target module directory.")

	err := TgCmd.MarkFlagRequired("target-module")
	if err != nil {
		log.Fatal(err)
	}

	_ = viper.BindPFlag("plan", TgCmd.Flags().Lookup("plan"))
	_ = viper.BindPFlag("plan-all", TgCmd.Flags().Lookup("plan-all"))
	_ = viper.BindPFlag("apply", TgCmd.Flags().Lookup("apply"))
	_ = viper.BindPFlag("apply-all", TgCmd.Flags().Lookup("apply-all"))
	_ = viper.BindPFlag("destroy", TgCmd.Flags().Lookup("destroy"))
	_ = viper.BindPFlag("destroy-all", TgCmd.Flags().Lookup("destroy-all"))
	_ = viper.BindPFlag("tg-commands", TgCmd.Flags().Lookup("tg-commands"))
	_ = viper.BindPFlag("target-module", TgCmd.Flags().Lookup("target-module"))
}

func init() {
	addTgFlags()
}
