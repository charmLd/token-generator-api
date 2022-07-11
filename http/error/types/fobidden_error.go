package types

import "fmt"

// ForbiddenError is the type of errors thrown by services talking to third party APIs.
type ForbiddenError struct {
	msg string
}

// New creates a new ForbiddenError instance.
func (e *ForbiddenError) New(message string) (err error) {
	err = &ForbiddenError{
		msg: message,
	}
	return
}

// Error returns the ForbiddenError message.
func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("|ForbiddenError|%s", e.msg)
}
