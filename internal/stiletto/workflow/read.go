package workflow

import (
	"github.com/Excoriate/stiletto/internal/stiletto"
)

type WorkflowFileReaderInstance struct {
	Stiletto stiletto.Stiletto
}

type WorkflowFileReader interface {
	ReadWorkflow(workflow string) error
}

func NewWorkflowFileReader() *WorkflowFileReaderInstance {
	s := stiletto.Init()

	return &WorkflowFileReaderInstance{
		Stiletto: s,
	}
}

func (w *WorkflowFileReaderInstance) ReadWorkflow(workflow string) error {
	logger := w.Stiletto.Logger

	if preConditionsErr := AreWorkflowFileValidationsPassed(logger,
		workflow); preConditionsErr != nil {
		return preConditionsErr
	}

	if schemaErr := IsWorkflowSchemaCompliant(logger, workflow); schemaErr != nil {
		return schemaErr
	}

	return nil
}
