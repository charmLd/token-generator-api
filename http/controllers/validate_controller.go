package controllers

import (
	"errors"
	"github.com/charmLd/token-generator-api/domain/usecases"
	"github.com/charmLd/token-generator-api/http/error"
	"github.com/charmLd/token-generator-api/http/error/types"
	"github.com/charmLd/token-generator-api/http/transport/request"
	"github.com/charmLd/token-generator-api/http/transport/request/unpackers"
	"github.com/charmLd/token-generator-api/http/transport/response"
	"github.com/charmLd/token-generator-api/http/transport/response/transformers"
	"github.com/charmLd/token-generator-api/http/validator"
	"net/http"
	"strings"
)

func (ctl *BaseController) PublicTokenValidateCOntrollerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	token := r.Header.Get("Authorization")
	if token == "" {
		error.Handle(r.Context(), (&types.ForbiddenError{}).New("authentication header invalid"), w)
		return
	}
	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) != 2 {
		error.Handle(r.Context(), (&types.ForbiddenError{}).New("authentication header token not inserted"), w)
		return
	}
	token = strings.TrimSpace(splitToken[1])

	validateRequest := unpackers.ValidateRequest{}
	err := request.Unpack(r, &validateRequest)
	if err != nil {
		error.Handle(ctx, err, w)
		return
	}

	// validate unpacked data
	errs := validator.Validate(ctx, validateRequest)
	if errs != nil {
		error.HandleUnpackValidationErrors(ctx, errs, w)
		return
	}

	validateUseCaseRequest := usecases.ValidateRequest{AuthToken: validateRequest.Jwt}
	valid, err := ctl.AuthUseCase.Validate(ctx, validateUseCaseRequest, token)
	if err != nil {
		error.Handle(ctx, err, w)
		return
	}

	if !valid {
		//send the appropriate error
		error.Handle(ctx, errors.New("invalid token"), w)
		return
	}

	var validateRes transformers.ValidateResponse
	validateRes.Success = "SUCCESS"
	response.Send(w, response.Transform(validateRes), http.StatusCreated)
}
