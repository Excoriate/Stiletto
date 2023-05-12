package infra

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Main IaC tools
	tgPlan    bool
	tgPlanAll bool

	tgApply    bool
	tgApplyAll bool

	tgDestroy    bool
	tgDestroyAll bool
)

var TgCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "terragrunt",
	Long:    `The 'terragrunt' command automate and perform several infra-related actions using Terragrunt.`,
	Example: `
	  stiletto infra terragrunt --plan --target-module=module1`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func addTgFlags() {
	TgCmd.Flags().BoolVarP(&tgPlan, "plan", "", false, "Run a plan.")
	TgCmd.Flags().BoolVarP(&tgPlanAll, "plan-all", "", false, "Run a plan-all.")
	TgCmd.Flags().BoolVarP(&tgApply, "apply", "", false, "Run an apply.")
	TgCmd.Flags().BoolVarP(&tgApplyAll, "apply-all", "", false, "Run an apply-all.")
	TgCmd.Flags().BoolVarP(&tgDestroy, "destroy", "", false, "Run a destroy.")
	TgCmd.Flags().BoolVarP(&tgDestroyAll, "destroy-all", "", false, "Run a destroy-all.")

	_ = viper.BindPFlag("plan", TgCmd.Flags().Lookup("plan"))
	_ = viper.BindPFlag("plan-all", TgCmd.Flags().Lookup("plan-all"))
	_ = viper.BindPFlag("apply", TgCmd.Flags().Lookup("apply"))
	_ = viper.BindPFlag("apply-all", TgCmd.Flags().Lookup("apply-all"))
	_ = viper.BindPFlag("destroy", TgCmd.Flags().Lookup("destroy"))
	_ = viper.BindPFlag("destroy-all", TgCmd.Flags().Lookup("destroy-all"))
}
