package cli

import (
	"github.com/spf13/cobra"
)

var (
	PipelineName string
)

var PipelineCMD = &cobra.Command{
	Version: "v0.0.1",
	Use:     "pipeline",
	Long:    "asd",
	Example: "asd",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func setChildCommands() {
	PipelineCMD.AddCommand(PipelineRunCMD)
}

func addCArguments() {
	//PipelineCMD.Flags().StringVarP(&PipelineName,
	//	"name",
	//	"n", "",
	//	"Name of the pipeline to run.")
	//
	//_ = PipelineCMD.MarkFlagRequired("name")
}

func init() {
	addCArguments()
	setChildCommands()
}
