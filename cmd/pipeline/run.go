package pipeline

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/display"
	"github.com/Excoriate/stiletto/internal/stiletto"
	"github.com/Excoriate/stiletto/internal/stiletto/workflow"
)

type RunInstance struct {
	Stiletto stiletto.Stiletto
}

type Runner interface {
	RunWorkflowFile(path string) error
}

func (i *RunInstance) RunWorkflowFile(workflowPath string) error {
	wr := workflow.NewWorkflowFileReader()
	err := wr.ReadWorkflow(workflowPath)
	if err != nil {
		display.UXError("workflow", fmt.Sprintf("workflow file %s is not valid", workflowPath), err)
		return err
	}

	return nil
}
