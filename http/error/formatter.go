package error

import (
	"context"
	"encoding/json"
	domainError "github.com/charmLd/token-generator-api/domain/error"
	"github.com/charmLd/token-generator-api/domain/globals"
	"github.com/charmLd/token-generator-api/http/error/types"
	"github.com/charmLd/token-generator-api/http/transport/response"
	"github.com/charmLd/token-generator-api/http/transport/response/transformers"
)

// Format error details.
func format(ctx context.Context, err error) []byte {
	wrapper := response.Error{}
	var payload interface{}

	switch err.(type) {
	case *domainError.Error:
		resolvedError, _ := err.(*domainError.Error)

		switch resolvedError.ErrorType {
		case domainError.DOMAIN, domainError.ADAPTER, domainError.SERVICE:
			payload = formatCustomError(ctx, *resolvedError)
			break
		default:
			payload = formatUnknownError(ctx, err)
			break
		}

		break
	case *types.ValidationError:

		payload = formatValidationStructureError(ctx, err)
		break
	case *types.ForbiddenError:
		payload = formatForbiddenStructureError(ctx, err)
		break
	case *types.UnAuthorizeError:
		payload = formatAuthorizationError(ctx, err)
		break
	default:
		payload = formatUnknownError(ctx, err)
		break
	}

	wrapper.Payload = payload
	message, _ := json.Marshal(wrapper)
	return message

}

// Format custom errors.
func formatCustomError(ctx context.Context, domainErr domainError.Error) (er []transformers.ErrorTransformer) {
	payload := transformers.ErrorTransformer{}
	payload.CorrelationID = ctx.Value(globals.UUIDKey)
	payload.Message = domainErr.Message
	payload.Code = domainErr.Code
	payload.DeveloperMessage = domainErr.DeveloperMessage
	er = append(er, payload)
	return er
}

// Format validation structure errors.
// These occur when the format of the sent data structure does not match the expected format.
func formatValidationStructureError(ctx context.Context, err error) (er []transformers.ErrorTransformer) {

	payload := transformers.ErrorTransformer{}
	payload.CorrelationID = ctx.Value(globals.UUIDKey)
	payload.Code = "CORE-1000"
	payload.DeveloperMessage = err.Error()
	payload.Message = "Unprocessable request"
	er = append(er, payload)
	return er
}

// Format errors of unhandled types.
func formatUnknownError(ctx context.Context, err error) (er []transformers.ErrorTransformer) {

	payload := transformers.ErrorTransformer{}
	payload.CorrelationID = ctx.Value(globals.UUIDKey)
	payload.Code = "CORE-5000"
	payload.DeveloperMessage = err.Error()
	payload.Message = "Unknown Error"
	er = append(er, payload)
	return er
}

// Format validation structure errors.
// These occur when the format of the sent data structure does not match the expected format.
func formatForbiddenStructureError(ctx context.Context, err error) (er []transformers.ErrorTransformer) {

	payload := transformers.ErrorTransformer{}
	payload.CorrelationID = ctx.Value(globals.UUIDKey)
	payload.Code = "CORE-1002"
	payload.DeveloperMessage = err.Error()
	payload.Message = "Forbidden Error"
	er = append(er, payload)
	return er
}

func formatAuthorizationError(ctx context.Context, err error) (er []transformers.ErrorTransformer) {

	payload := transformers.ErrorTransformer{}
	payload.CorrelationID = ctx.Value(globals.UUIDKey)
	payload.Code = "CORE-1003"
	payload.DeveloperMessage = err.Error()
	payload.Message = "Authorization Error"

	er = append(er, payload)
	return er
}
