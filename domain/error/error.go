package error

import (
	"context"
	"fmt"
	"github.com/charmLd/token-generator-api/domain/globals"
)

type errorType int

var (
	DOMAIN  errorType = 1
	ADAPTER errorType = 2
	SERVICE errorType = 3
	//DATA_ERROR errorType = iota
	//MIDDLEWARE_ERROR
	//SERVER_ERROR
	//SERVICE_ERROR
	//VALIDATION_ERROR
	//DOMAIN_ERROR
	//ADAPTER_ERROR
	//UNKNOWN_ERROR
)

// Error is the type of errors thrown by business logic.
type Error struct {
	CorrelationId    interface{}
	Code             string
	Message          string
	DeveloperMessage string
	ErrorType        errorType
}

// New creates a new DomainError instance.
func (e *Error) New(ctx context.Context, t errorType, message string, code string, details string) error {

	err := &Error{}

	err.CorrelationId = ctx.Value(globals.UUIDKey)
	err.Code = code
	err.Message = message
	err.DeveloperMessage = details
	err.ErrorType = t

	return err
}

// Error returns the DataError message.
func (e *Error) Error() string {
	return fmt.Sprintf("%s|%s|%s|Error|%s", e.CorrelationId, e.Code, e.Message, e.DeveloperMessage)
}
