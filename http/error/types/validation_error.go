package types

// ValidationError is the type of errors thrown by the request validator.
type ValidationError struct {
	details string
}

// New creates a new ValidationError instance.
func (e *ValidationError) New(details string) error {

	return &ValidationError{
		details: details,
	}
}

// Error returns the ValidationError message.
func (e *ValidationError) Error() string {
	return e.details
}
