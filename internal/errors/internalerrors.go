package errors

import "fmt"

// InternalStilettoError is an error that is thrown when an internal Stiletto error occurs.
type InternalStilettoError struct {
	Detail string
}

// Error returns the error message.
func (e *InternalStilettoError) Error() string {
	return fmt.Sprintf("Internal Stiletto error: %s", e.Detail)
}

// NewInternalStilettoError returns a new InternalStilettoError.
func NewInternalStilettoError(detail string) error {
	return &InternalStilettoError{
		Detail: detail,
	}
}
