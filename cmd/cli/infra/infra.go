package infra

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Main IaC tools
	terraform  bool
	terragrunt bool

	AWSAccessKeyID string
	AWSSecretKey   string
	AWSRegion      string

	targetModule string
)

var Cmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "infra",
	Long:    `The 'infra' command automate and perform several infra-related actions using either Terraform or Terragrunt.`,
	Example: `
  stiletto infra terragrunt --plan --target-module=module1`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func addFlags() {
	Cmd.PersistentFlags().StringVarP(&AWSAccessKeyID,
		"aws-access-key-id",
		"", "",
		"The AWS Access Key ID. If it's not set, it'll be read from the AWS_ACCESS_KEY_ID environment variable.")

	Cmd.PersistentFlags().StringVarP(&AWSSecretKey,
		"aws-secret-key",
		"", "",
		"The AWS Secret Access Key. If it's not set, it'll be read from the AWS_SECRET_ACCESS_KEY environment variable.")

	Cmd.PersistentFlags().StringVarP(&AWSRegion,
		"aws-region",
		"", "",
		"The AWS Region. If it's not set, it'll be read from the AWS_REGION environment variable.")

	Cmd.PersistentFlags().BoolVarP(&terraform, "terraform", "", false, "Use Terraform.")
	Cmd.PersistentFlags().BoolVarP(&terragrunt, "terragrunt", "", false, "Use Terragrunt.")
	Cmd.PersistentFlags().StringVarP(&targetModule, "target-module", "", "", "The target module to run the command on.")

	_ = viper.BindPFlag("aws-access-key-id", Cmd.PersistentFlags().Lookup("aws-access-key-id"))
	_ = viper.BindPFlag("aws-secret-access-key", Cmd.PersistentFlags().Lookup("aws-secret-key"))
	_ = viper.BindPFlag("aws-region", Cmd.PersistentFlags().Lookup("aws-region"))
	_ = viper.BindPFlag("terraform", Cmd.PersistentFlags().Lookup("terraform"))
	_ = viper.BindPFlag("terragrunt", Cmd.PersistentFlags().Lookup("terragrunt"))
	_ = viper.BindPFlag("target-module", Cmd.PersistentFlags().Lookup("target-module"))
}

func init() {
	Cmd.AddCommand(TgCmd)
	addFlags()
}
