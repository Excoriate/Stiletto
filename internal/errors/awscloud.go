package errors

import "fmt"

const awsCloudConfigErrorPrefix = "AWS configuration error: "
const awsCloudExecutionErrorPrefix = "AWS execution error: "

type AWSConfigurationError struct {
	Details string
	Err     error
}

func (e *AWSConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", awsCloudConfigErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", awsCloudConfigErrorPrefix, e.Details)
}

func NewAWSCfgError(details string, err error) *AWSConfigurationError {
	return &AWSConfigurationError{
		Details: details,
		Err:     err,
	}
}

type AWSExecutionError struct {
	Details string
	Err     error
}

func (e *AWSExecutionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %s", awsCloudExecutionErrorPrefix, e.Details, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", awsCloudExecutionErrorPrefix, e.Details)
}

func NewAWSExecutionError(details string, err error) *AWSExecutionError {
	return &AWSExecutionError{
		Details: details,
		Err:     err,
	}
}
