package usecases

import (
	"context"

	"fmt"
	"math/rand"
	"time"

	"github.com/charmLd/token-generator-api/domain/entities"
	appErr "github.com/charmLd/token-generator-api/domain/error"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type ValidateRequest struct {
	AuthToken string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func (au *AuthUseCase) GenerateToken(ctx context.Context, userId string) (token string, err error) {
	//generate random alphanumerical number as a string
	newToken := entities.TokenGenRequest{}
	//newToken.UniqueString = userId + RandStringRunes(3) // uniques string with 3 characters
	newToken.UniqueString = uuid.New().String()
	//generate new token
	token, err = au.TokenAdapter.GenerateUniqueToken(ctx, newToken)
	fmt.Println("new token: ", token)

	tokenReq := entities.Token{}
	tokenReq.ID = newToken.UniqueString
	tokenReq.GeneratedToken = token
	tokenReq.UserId = userId
	tokenReq.CreatedAt = time.Now()
	expiryTime := tokenReq.CreatedAt.Unix() + int64(au.Config.GeneratedTokenExpiry)*60*60*24 //currentTime + ta.Cfg.RTExpiry*60*60*24 //converting days into seconds

	tokenReq.Expiry = time.Unix(expiryTime, 0)

	//Write the new token to the db
	err = au.TokenRepository.InsertUniqueToken(ctx, tokenReq)
	fmt.Println("insert error:", err)
	if err != nil {
		err = au.throwTokenIssueError(ctx)
		return
	}
	return token, err

}

func (au *AuthUseCase) GetAllTokenDetailsByFilters(ctx context.Context, fetchDetails entities.TokenDetailsReqParam) ([]entities.Token, error) {

	return au.TokenRepository.GetAllTokenForFilter(ctx, fetchDetails)
}

func (au *AuthUseCase) Validate(ctx context.Context, validateReq ValidateRequest, token string) (ok bool, err error) {

	JwtPayloadData, err := au.TokenAdapter.ValidateLoginJWToken(ctx, token)

	if err != nil {
		log.Error(ctx, "invalid user: ", err)

		return
	}
	fmt.Println("user : ", JwtPayloadData.UserID)
	//go func to update the user login
	defer func() {
		err = au.UserRepository.UpdateLastLogin(ctx, fmt.Sprint(JwtPayloadData.UserID))
		if err != nil {
			log.Error(ctx, "updating last login failed ", JwtPayloadData.UserID)
			return
		}
	}()

	jwtClaims, err := au.TokenAdapter.ValidateGeneratedToken(ctx, validateReq.AuthToken)
	fmt.Println("claim:", jwtClaims)
	if err != nil {
		//these 2 errors should be proceeded to http to give necessary error codes
		if err.Error() == appErr.UnAuthorizedTokenStr {
			err = au.throwUnauthorizedError(ctx, err.Error())
			return
		}
		if err.Error() == appErr.ExpiredTokenStr {
			err = au.throwUnauthorizedError(ctx, "expired jwt token")
			return
		}
		err = au.throwUnauthorizedError(ctx, "token validation failed")
		log.Error(ctx, "token validation failed")
		return
	}

	/*
		//fetching user from refresh claims
		user, userExists, err := au.UserRepository.GetUserByUserIdAndAppName(ctx, jwtClaims.UserID, jwtClaims.AppID)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("validating token failed for user %v", jwtClaims.UserID))
			return
		}

		if !userExists {
			log.Error(ctx, "user does not exist", fmt.Sprintf("validating token failed for %v", jwtClaims.UserID))
			err = au.throwUserNotExistError(ctx)
			return
		}

		// Check if the user is banned or not
		if user.IsBlacklisted {
			log.Error(ctx, "user is blacklisted", fmt.Sprintf("validating token failed for %v", jwtClaims.UserID))
			err = au.throwUserBlacklistedError(ctx)
			return
		}

		//returns tokens which are not blacklisted
		_, exists, err := au.TokenRepository.GetTokenByTokenID(ctx, jwtClaims.TokenID)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("validating failed for %v", jwtClaims.UserID))
			return
		}

		if !exists {
			log.Error(ctx, "token might have been blacklisted or expired", jwtClaims.TokenID, fmt.Sprintf("validating token failed for %v", jwtClaims.UserID))
			err = au.throwInvalidJwtTokenIDError(ctx)
			return
		}

		log.Trace(ctx, "jwt token validation is successful", jwtClaims.UserID) */
	return true, nil
}

