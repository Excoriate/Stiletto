package cli

import (
	"github.com/Excoriate/stiletto/cmd/pipeline"
	"github.com/spf13/cobra"
	"os"
)

var WorkflowFilePath string

var PipelineRunCMD = &cobra.Command{
	Version: "v0.0.1",
	Use:     "run",
	Long:    "asd",
	Example: "asd",
	Run: func(cmd *cobra.Command, args []string) {
		i := pipeline.RunInstance{}
		err := i.RunWorkflowFile(WorkflowFilePath)

		if err != nil {
			os.Exit(1)
		}
	},
}

func addPipelineRunArgs() {
	PipelineRunCMD.Flags().StringVarP(&WorkflowFilePath,
		"file",
		"f", "",
		"File path where the <.job-something.yaml> (or equivalent) resides")

	_ = PipelineRunCMD.MarkFlagRequired("file")
}

func init() {
	addPipelineRunArgs()
}
