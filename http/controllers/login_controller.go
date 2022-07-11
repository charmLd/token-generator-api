package controllers

import (
	"net/http"

	"github.com/charmLd/token-generator-api/domain/usecases"
	"github.com/charmLd/token-generator-api/http/error"
	"github.com/charmLd/token-generator-api/http/transport/request"
	"github.com/charmLd/token-generator-api/http/transport/request/unpackers"
	"github.com/charmLd/token-generator-api/http/transport/response"
	"github.com/charmLd/token-generator-api/http/transport/response/transformers"
	"github.com/charmLd/token-generator-api/http/validator"
)

func (ctl *BaseController) EmailLoginControllerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	emailLoginRequest := unpackers.EmailLoginRequest{}

	err := request.Unpack(r, &emailLoginRequest)
	if err != nil {
		error.Handle(ctx, err, w)
		return
	}
	// validate unpacked data
	errs := validator.Validate(ctx, emailLoginRequest)
	if errs != nil {
		error.HandleUnpackValidationErrors(ctx, errs, w)
		return
	}
	if err := validator.IsEmailValid(emailLoginRequest.Email); err != nil {
		error.Handle(ctx, err, w)
		return
	}
	loginUseCaseReq := usecases.EmailLoginRequest{
		Email:    emailLoginRequest.Email,
		Password: emailLoginRequest.Password,
	}

	token, err := ctl.AuthUseCase.AuthenticateByEmail(ctx, loginUseCaseReq)
	if err != nil {
		error.Handle(ctx, err, w)
		return
	}

	var loginResponse transformers.LoginResponse
	loginResponse.JWT = token.GeneratedToken
	loginResponse.Status = "SUCCESS"

	response.Send(w, response.Transform(loginResponse), http.StatusCreated)
}
