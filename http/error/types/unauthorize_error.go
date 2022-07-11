package types

import "fmt"

// UnAuthorizeError is the type of errors thrown by services talking to third party APIs.
type UnAuthorizeError struct {
	msg string
}

// New creates a new UnAuthorizeError instance.
func (e *UnAuthorizeError) New(message string) (err error) {
	err = &UnAuthorizeError{
		msg: message,
	}
	return
}

// Error returns the UnAuthorizeError message.
func (e *UnAuthorizeError) Error() string {
	return fmt.Sprintf("|UnAuthorizeError|%s", e.msg)
}
