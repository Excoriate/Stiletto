package aws

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	AWSAccessKeyID string
	AWSSecretKey   string
	AWSRegion      string
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
	Cmd.PersistentFlags().StringVarP(&AWSAccessKeyID,
		"aws-access-key-id",
		"i", "",
		"The AWS Access Key ID. If it's not set, it'll be read from the AWS_ACCESS_KEY_ID environment variable.")

	Cmd.PersistentFlags().StringVarP(&AWSSecretKey,
		"aws-secret-key",
		"k", "",
		"The AWS Secret Access Key. If it's not set, it'll be read from the AWS_SECRET_ACCESS_KEY environment variable.")

	Cmd.PersistentFlags().StringVarP(&AWSRegion,
		"aws-region",
		"", "",
		"The AWS Region. If it's not set, it'll be read from the AWS_REGION environment variable.")

	_ = viper.BindPFlag("aws-access-key-id", Cmd.PersistentFlags().Lookup("aws-access-key-id"))
	_ = viper.BindPFlag("aws-secret-access-key", Cmd.PersistentFlags().Lookup("aws-secret-key"))
	_ = viper.BindPFlag("aws-region", Cmd.PersistentFlags().Lookup("aws-region"))
}

func init() {
	Cmd.AddCommand(ECRCmd)
	Cmd.AddCommand(ECSCmd)
	addFlags()
}
