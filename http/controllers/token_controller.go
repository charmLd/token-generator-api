package controllers

import (
	"github.com/charmLd/token-generator-api/domain/entities"
	"github.com/charmLd/token-generator-api/http/error"

	"github.com/charmLd/token-generator-api/http/transport/request"
	"github.com/charmLd/token-generator-api/http/transport/request/unpackers"
	"github.com/charmLd/token-generator-api/http/transport/response"
	"github.com/charmLd/token-generator-api/http/transport/response/transformers"
	"github.com/charmLd/token-generator-api/http/validator"
	"github.com/gorilla/mux"
	"net/http"
)

func (ctl *BaseController) AdminTokenGenerateControllerFun(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	ctx := r.Context()

	tokenUnpacker := unpackers.TokenCreateRequest{}

	err := request.Unpack(r, &tokenUnpacker)
	if err != nil {
		error.Handle(ctx, err, w)
		return
	}
	// validate unpacked data
	errs := validator.Validate(ctx, tokenUnpacker)
	if errs != nil {
		error.HandleUnpackValidationErrors(ctx, errs, w)
		return
	}

	token, err := ctl.AuthUseCase.GenerateToken(ctx, tokenUnpacker.UserID)
	if err != nil {
		error.Handle(ctx, err, w)
		return
	}

	var tokenRes transformers.TokenResponse
	tokenRes.Token = token
	tokenRes.Status = "SUCCESS"

	response.Send(w, response.Transform(tokenRes), http.StatusCreated)

}

func (ctl *BaseController) AdminTokenRevokeControllerFun(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	ctx := r.Context()

	// unpack request
	adminRevokeUnpacker := unpackers.AdminRevokeUnpackers{}

	err := request.Unpack(r, &adminRevokeUnpacker)
	if err != nil {
		error.Handle(ctx, err, w)
		return
	}
	// validate unpacked data
	errs := validator.Validate(ctx, adminRevokeUnpacker)
	if errs != nil {
		error.HandleUnpackValidationErrors(ctx, errs, w)
		return
	}

	err = ctl.AuthUseCase.AdminRevokeToken(ctx, adminRevokeUnpacker.InviteToken)
	if err != nil {
		error.Handle(ctx, err, w)
		return
	}

	// transform
	revokeResponse := transformers.SuccessTransformer{
		Status: "SUCCESS",
	}
	// send response
	response.Send(w, response.Transform(revokeResponse), http.StatusCreated)
}

func (ctl *BaseController) AdminTokenFetchControllerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := validator.IsUserIDValidUInt(userID); err != nil {
		error.Handle(ctx, err, w)
		return
	}

	Qparams := r.URL.Query()

	var blacklisted []string
	var blacklistedOk bool

	if len(Qparams) != 0 {

		if len(Qparams["blacklisted"]) == 0 {
			blacklisted = append(blacklisted, "")
			blacklistedOk = false
		} else {
			blacklisted, blacklistedOk = Qparams["blacklisted"]
			if err := validator.IsBoolValid(blacklisted[0]); err != nil {
				error.Handle(ctx, err, w)
				return
			}
		}

		fetchDetails := entities.TokenDetailsReqParam{
			Balcklisted: entities.TokenParams{
				IsOK:  blacklistedOk,
				Value: blacklisted[0],
			},
			UserId: userID,
		}
		tokenResponseArray, err := ctl.AuthUseCase.GetAllTokenDetailsByFilters(ctx, fetchDetails)

		if err != nil {
			error.Handle(ctx, err, w)
			return
		}
		tRes := transformers.TokenDetailResponse{
			UserId: userID,
			Tokens: tokenResponseArray,
		}
		response.Send(w, response.Transform(tRes), http.StatusOK)
	}

}
