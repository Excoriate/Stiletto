package aws

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	AccessKeyID string
	SecretKey   string
	Region      string
)

var Cmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "aws",
	Long: `The 'aws' command automate and perform several aws-related actions (E.
g: push images to ECR, deploy into ECS, etc.).
You can specify the tasks you want to perform using the provided --task flag.`,
	Example: `
  # Push an image into ECR:
  stiletto aws ecr --task=push`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func addFlags() {
	Cmd.PersistentFlags().StringVarP(&AccessKeyID,
		"aws-creds-access-key-id",
		"i", "",
		"The AWS Access Key ID. If it's not set, it'll be read from the AWS_ACCESS_KEY_ID environment variable.")

	Cmd.PersistentFlags().StringVarP(&SecretKey,
		"aws-creds-secret-key",
		"k", "",
		"The AWS Secret Access Key. If it's not set, it'll be read from the AWS_SECRET_ACCESS_KEY environment variable.")

	Cmd.PersistentFlags().StringVarP(&Region,
		"aws-creds-region",
		"r", "",
		"The AWS Region. If it's not set, "+
			"it'll be read from the AWS_REGION and for the AWS_DEFAULT_REGION environment variable.")

	_ = viper.BindPFlag("aws-creds-access-key-id", Cmd.PersistentFlags().Lookup("aws-creds-access-key-id"))
	_ = viper.BindPFlag("aws-creds-secret-key", Cmd.PersistentFlags().Lookup("aws-creds-secret-key"))
	_ = viper.BindPFlag("aws-creds-region", Cmd.PersistentFlags().Lookup("aws-creds-region"))
}

func init() {
	addFlags()
	Cmd.AddCommand(ECRCmd)
	Cmd.AddCommand(ECSCmd)
}
