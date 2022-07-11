package usecases

import (
	"context"

	"fmt"
	"math/rand"
	"time"

	"github.com/charmLd/token-generator-api/domain/entities"
	log "github.com/sirupsen/logrus"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (au *AuthUseCase) GenerateToken(ctx context.Context, userId string) (token string, err error) {
	//generate random alphanumerical number as a string

	randomStringValue := RandStringRunes(3) + RandStringRunes(4) + RandStringRunes(3)
	newToken := entities.TokenGenRequest{}
	//generate jwt token
	token, err = au.TokenAdapter.GenerateUniqueToken(ctx, newToken)

	tokenReq := entities.Token{}
	tokenReq.Token = randomStringValue
	tokenReq.GeneratedToken = token
	tokenReq.UserId = userId
	tokenReq.CreatedAt = time.Now()
	expiryTime := tokenReq.CreatedAt.Unix() + int64(au.Config.GeneratedTokenExpiry) //currentTime + ta.Cfg.RTExpiry*60*60*24 //converting days into seconds

	tokenReq.Expiry = time.Unix(expiryTime, 0)

	//Write the new token to the db
	err = au.TokenRepository.InsertUniqueToken(ctx, tokenReq)

	if err != nil {
		err = au.throwTokenIssueError(ctx)
		return
	}
	return randomStringValue, err

}

func (au *AuthUseCase) GetAllTokenDetailsByFilters(ctx context.Context, fetchDetails entities.TokenDetailsReqParam) ([]entities.Token, error) {

	return au.TokenRepository.GetAllTokenForFilter(ctx, fetchDetails)
}

func (au *AuthUseCase) Validate(ctx context.Context, validateReq entities.ValidateRequest, token string) (ok bool, err error) {
	JwtPayloadData, err := au.TokenAdapter.ValidateLoginJWToken(ctx, token)
	if err != nil {
		log.Error(ctx, "invalid user: ", err)

		return
	}

	//go func to update the user login
	defer func() {
		err = au.UserRepository.UpdateLastLogin(ctx, fmt.Sprint(JwtPayloadData.UserID))
		if err != nil {
			log.Error(ctx, "updating last login failed ", validateReq.UserId)
			return
		}
	}()
	//jwtClaims, err := au.TokenAdapter.ValidateGeneratedToken(ctx, validateReq.InviteToken)

	inviteTokenInfo, err := au.TokenRepository.FetchTokenInfo(ctx, validateReq)

	if err != nil {
		err = au.throwInviteTokenFetchError(ctx, err)
		log.Error(ctx, "token validation failed")
		return false, err
	}
	if time.Since(inviteTokenInfo.CreatedAt) > time.Second*time.Duration(au.Config.GeneratedTokenExpiry) {
		fmt.Println(ctx, "expired token", inviteTokenInfo.Expiry)
		return false, au.throwInviteTokenExpiryError(ctx, "Token expired")
	}

	return true, nil
}
