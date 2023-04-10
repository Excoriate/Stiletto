package errors

import "fmt"

const actionCfgError = "Action configuration error: "
const actionExecError = "Action execution error: "

type ActionCfgError struct {
	Details string
	Err     error
}

func (e *ActionCfgError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", actionCfgError, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", actionExecError, e.Details)
}

func NewActionCfgError(details string, err error) *ActionCfgError {
	return &ActionCfgError{
		Details: details,
		Err:     err,
	}
}

type ActionExecError struct {
	Details string
	Err     error
}

func (e *ActionExecError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", actionExecError, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", actionExecError, e.Details)
}

func NewActionExecError(details string, err error) *ActionExecError {
	return &ActionExecError{
		Details: details,
		Err:     err,
	}
}
