package error

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/charmLd/token-generator-api/http/transport/response/transformers"

	domainError "github.com/charmLd/token-generator-api/domain/error"
	"github.com/charmLd/token-generator-api/http/error/types"
	"github.com/charmLd/token-generator-api/http/transport/response"
	log "github.com/sirupsen/logrus"
)

// Handle handles all errors globally.
func Handle(ctx context.Context, err error, w http.ResponseWriter) {

	switch err.(type) {
	case *types.ServerError:
		log.Error(ctx, " http.error.handler.Handle", "Server Error", err)
		response.Send(w, format(ctx, err), http.StatusInternalServerError)
		break
	case *types.AdapterError, *types.MiddlewareError, *types.DataError, *types.ServiceError, *domainError.Error:
		log.Error(ctx, " http.error.handler.Handle ", "Other Error", err)
		response.Send(w, format(ctx, err), http.StatusBadRequest)
		break
	case *types.ValidationError:
		log.Error(ctx, " http.error.handler.Handle ", "Validation Structure Error", err)
		response.Send(w, format(ctx, err), http.StatusUnprocessableEntity)
		break
	case *types.ForbiddenError:
		log.Error(ctx, " http.error.handler.Handle ", "ForbiddenError Error", err)
		response.Send(w, format(ctx, err), http.StatusForbidden)
		break
	case *types.UnAuthorizeError, *domainError.UnAuthorizeDomainError:
		log.Error(ctx, " http.error.handler.Handle ", "UnAuthorizeError Error", err)
		response.Send(w, format(ctx, err), http.StatusUnauthorized)
		break
	default:
		log.Error(ctx, " http.error.handler.Handle ", "Unknown Error", err)
		response.Send(w, format(ctx, err), http.StatusInternalServerError)
		break
	}

	return
}

// HandleValidationErrors specifically handles validation errors thrown by the validator.
func HandleValidationErrors(ctx context.Context, errs []error, w http.ResponseWriter) {

	wrapper := response.Error{}
	var errArray []transformers.ErrorTransformer
	for _, err := range errs {
		fom := formatValidationStructureError(ctx, err)
		errArray = append(errArray, fom[0])
	}
	wrapper.Payload = errArray
	message, _ := json.Marshal(wrapper)
	log.Error(ctx, " http.error.handler.HandleValidationErrors ", "Validation Errors", string(message))
	response.Send(w, message, http.StatusUnprocessableEntity)

	return
}
func HandleUnpackValidationErrors(ctx context.Context, errs []transformers.ErrorTransformer, w http.ResponseWriter) {

	wrapper := response.Error{}
	wrapper.Payload = errs
	message, _ := json.Marshal(wrapper)
	log.Error(ctx, " http.error.handler.HandleValidationErrors ", "Validation Errors", string(message))
	response.Send(w, message, http.StatusUnprocessableEntity)

	return
}
