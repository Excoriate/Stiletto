package cli

import (
	"context"
	"os"
)

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "stiletto",
	Long:    "asd",
	Example: "asd",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(PipelineCMD)
}
