package usecases

import (
	"context"
	//log "github.com/sirupsen/logrus"
)

// AdminRevokeToken handles the business usecase for invalidate auth tokens by admin
func (au *AuthUseCase) AdminRevokeToken(ctx context.Context, userId string) (err error) {
	err = au.TokenRepository.Revoke(ctx, userId)
	if err != nil {
		return au.throwRevokeError(ctx, err.Error())
	}

	return
}
