package usecases

import (
	"context"
	"fmt"
	//log "github.com/sirupsen/logrus"

	//"github.com/charmLd/token-generator-api/domain/entities"
	appErr "github.com/charmLd/token-generator-api/domain/error"
)

// AdminRevokeToken handles the business usecase for invalidate auth tokens by admin
func (au *AuthUseCase) AdminRevokeToken(ctx context.Context, token string) (err error) {
	jwtPayload, err := au.TokenAdapter.DecodeGeneratedToken(ctx, token)
	if err != nil {
		return au.throwJWTDecodeError(ctx, err.Error())
	}
	fmt.Println("uniques value:", jwtPayload.UniqueString)
	err = au.TokenRepository.Revoke(ctx, jwtPayload.UniqueString)
	if err != nil {
		return au.throwRevokeError(ctx, err.Error())
	}

	return
}
func (au *AuthUseCase) throwRevokeError(ctx context.Context, reason string) error {
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "revoke Error",
		"CORE-2011",
		reason)
}

func (au *AuthUseCase) throwUnauthorizedError(ctx context.Context, reason string) error {
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "unauthorized Error", "CORE-2010", reason)
}
func (au *AuthUseCase) throwJWTDecodeError(ctx context.Context, reason string) error {
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "decode Error", "CORE-2010", reason)
}