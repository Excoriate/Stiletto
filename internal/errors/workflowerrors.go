package errors

import "fmt"

// WorkflowFileNotFoundError represents an error when a workflow file is not found.
type WorkflowFileNotFoundError struct {
	Filename string
}

// Error implements the error interface for WorkflowFileNotFoundError.
func (e *WorkflowFileNotFoundError) Error() string {
	return fmt.Sprintf("Workflow file not found: %s", e.Filename)
}

// NewWorkflowFileNotFoundError creates a new WorkflowFileNotFoundError with the given filename.
func NewWorkflowFileNotFoundError(filename string) *WorkflowFileNotFoundError {
	return &WorkflowFileNotFoundError{Filename: filename}
}

// WorkflowInvalidExtensionError represents an error when a workflow file has an invalid extension.
type WorkflowInvalidExtensionError struct {
	Filename string
}

// Error implements the error interface for WorkflowInvalidExtensionError.
func (e *WorkflowInvalidExtensionError) Error() string {
	return fmt.Sprintf("Workflow file has invalid extension: %s", e.Filename)
}

// NewWorkflowInvalidExtensionError creates a new WorkflowInvalidExtensionError with the given filename.
func NewWorkflowInvalidExtensionError(filename string) *WorkflowInvalidExtensionError {
	return &WorkflowInvalidExtensionError{Filename: filename}
}

// WorkflowCompilationError represents an error when a workflow file fails to compile.
type WorkflowCompilationError struct {
	Filename string
	Err      error
}

// Error implements the error interface for WorkflowCompilationError.
func (e *WorkflowCompilationError) Error() string {
	return fmt.Sprintf("Workflow file %s failed to compile, with error: %s", e.Filename, e.Err)
}

// NewWorkflowCompilationError creates a new WorkflowCompilationError with the given filename and error.
func NewWorkflowCompilationError(filename string, err error) *WorkflowCompilationError {
	return &WorkflowCompilationError{Filename: filename, Err: err}
}

// WorkflowSchemaVerificationError represents an error when a workflow file fails to verify against the schema.
type WorkflowSchemaVerificationError struct {
	Filename string
	Err      error
}

// Error implements the error interface for WorkflowSchemaVerificationError.
func (e *WorkflowSchemaVerificationError) Error() string {
	return fmt.Sprintf("Workflow file %s failed to verify against schema, with error: %s", e.Filename, e.Err)
}

// NewWorkflowSchemaVerificationError creates a new WorkflowSchemaVerificationError with the given filename and error.
func NewWorkflowSchemaVerificationError(filename string, err error) *WorkflowSchemaVerificationError {
	return &WorkflowSchemaVerificationError{Filename: filename, Err: err}
}
