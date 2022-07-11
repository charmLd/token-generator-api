package error

import "fmt"

// UnAuthorizeError is the type of errors thrown by services talking to third party APIs.
type UnAuthorizeDomainError struct {
	msg     string
	code    int
	details string
}

// New creates a new UnAuthorizeError instance.
func (e *UnAuthorizeDomainError) New(message string, code int, details string) error {

	err := &UnAuthorizeDomainError{
		msg:     message,
		code:    code,
		details: details,
	}
	return err
}

// Error returns the UnAuthorizeError message.
func (e *UnAuthorizeDomainError) Error() string {
	return fmt.Sprintf("%s|%d|UnAuthorizeError|%s", e.msg, e.code, e.details)
}
