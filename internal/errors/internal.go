package errors

import "fmt"

// InternalPipelineError is an error that is thrown when an internal Stiletto error occurs.
type InternalPipelineError struct {
	Detail string
}

// Error returns the error message.
func (e *InternalPipelineError) Error() string {
	return fmt.Sprintf("Internal PipelineCfg error: %s", e.Detail)
}

// NewInternalPipelineError returns a new InternalPipelineError.
func NewInternalPipelineError(detail string) error {
	return &InternalPipelineError{
		Detail: detail,
	}
}
