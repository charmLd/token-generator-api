package usecases

import (
	"context"
	appErr "github.com/charmLd/token-generator-api/domain/error"
)

func (au *AuthUseCase) throwPasswordError(ctx context.Context) error {
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "incorrect password",
		"CORE-2004",
		"incorrect password")
}

func (au *AuthUseCase) throwUserNotExistError(ctx context.Context) error {
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "invalid user",
		"CORE-2006",
		"invalid user")
}

func (au *AuthUseCase) throwTokenIssueError(ctx context.Context) error {
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "token Issue Error",
		"CORE-2012",
		"could not issue a token")
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
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "decode Error", "CORE-2009", reason)
}

func (au *AuthUseCase) throwInviteTokenFetchError(ctx context.Context, err error) error {
	return (&appErr.Error{}).New(ctx, appErr.ADAPTER, "mysql adapter", "CORE-2008", err.Error())
}
func (au *AuthUseCase) throwInviteTokenExpiryError(ctx context.Context, reason string) error {
	return (&appErr.Error{}).New(ctx, appErr.ADAPTER, "token invalid", "CORE-2007", reason)
}
