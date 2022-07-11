package types

import "fmt"

// DataError is the type of errors thrown by repositories.
type DataError struct {
	msg     string
	code    int
	details string
}

// New creates a new DataError instance.
func (e *DataError) New(message string, code int, details string) error {

	return &DataError{
		msg:     message,
		code:    code,
		details: details,
	}
}

// Error returns the DataError message.
func (e *DataError) Error() string {
	return fmt.Sprintf("%s|%d|DataError|%s", e.msg, e.code, e.details)
}
