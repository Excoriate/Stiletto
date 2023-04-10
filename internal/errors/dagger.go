package errors

import "fmt"

const daggerEngineErrorPrefix = "Dagger engine error: "
const daggerConfigurationErrorPrefix = "Dagger configuration error: "

type DaggerEngineError struct {
	Details string
	Err     error
}

func (e *DaggerEngineError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", daggerEngineErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", daggerEngineErrorPrefix, e.Details)
}

func NewDaggerEngineError(details string, err error) *DaggerEngineError {
	return &DaggerEngineError{
		Details: details,
		Err:     err,
	}
}

type DaggerConfigurationError struct {
	Details string
	Err     error
}

func (e *DaggerConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", daggerConfigurationErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", daggerConfigurationErrorPrefix, e.Details)
}

func NewDaggerConfigurationError(details string, err error) *DaggerConfigurationError {
	return &DaggerConfigurationError{
		Details: details,
		Err:     err,
	}
}
