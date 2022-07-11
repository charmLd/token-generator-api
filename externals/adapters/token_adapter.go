package adapters

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmLd/token-generator-api/domain/entities"
	error2 "github.com/charmLd/token-generator-api/domain/error"
	"github.com/charmLd/token-generator-api/util/config"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type tokenAdapter struct {
	Cfg config.TokenConfig
}

func NewTokenAdapter(cfg config.TokenConfig) *tokenAdapter {

	return &tokenAdapter{
		Cfg: cfg,
	}
}

//todo For future improvements
func (ta *tokenAdapter) GenerateUniqueToken(ctx context.Context, jwtClaims entities.TokenGenRequest) (jwToken string, err error) {

	var alg *jwt.SigningMethodHMAC
	//using HS256 Hashing Algorithm
	alg = jwt.SigningMethodHS256

	//jwt claims
	currentTime := time.Now().Unix()

	jwtClaims.IssuedAt = currentTime
	jwtClaims.ExpireAt = currentTime + int64(ta.Cfg.TokenExpiry)*60*60*24 //currentTime + ta.Cfg.RTExpiry*60*60*24 //converting days into seconds
	jwtClaims.IsBlacklisted = false

	t := jwt.NewWithClaims(alg, &jwtClaims)

	jwToken, err = t.SignedString([]byte(ta.Cfg.TokenSecretKey))
	if err != nil {

		log.Error(ctx, err, "signing jw token with secret key failed")
		return
	}

	log.Trace(ctx, "token is generated ")

	return
}

func (ta *tokenAdapter) GenerateLoginToken(ctx context.Context, jwtClaims entities.JWTClaims) (jwToken string, err error) {

	/* 	var alg *jwt.SigningMethodHMAC
	   	if ta.Cfg.HashAlgorithm == "HS256" {
	   		alg = jwt.SigningMethodHS256
	   	} else {
	   		log.Error(ctx, "invalid hashing method")
	   		return "", "", errors.New("invalid hashing method")
	   	} */
	var alg *jwt.SigningMethodHMAC
	//using HS256 Hashing Algorithm
	alg = jwt.SigningMethodHS256

	//jwt claims
	currentTime := time.Now().Unix()
	jwtClaims.IssuedAt = currentTime
	jwtClaims.ExpireAt = currentTime + int64(ta.Cfg.LoginTokenExpiry)*60*60*24 //currentTime + ta.Cfg.RTExpiry*60*60*24 //converting days into seconds

	t := jwt.NewWithClaims(alg, &jwtClaims)

	jwToken, err = t.SignedString([]byte(ta.Cfg.LoginSecretKey))
	if err != nil {

		log.Error(ctx, err, "signing jw token with secret key failed")
		return "", err
	}

	log.Trace(ctx, fmt.Sprintf("tokens generated for user: %v", jwtClaims.UserID))

	return
}

func (ta *tokenAdapter) ValidateLoginJWToken(ctx context.Context, token string) (jwtClaims *entities.JWTClaims, err error) {

	jwtClaims = &entities.JWTClaims{}
	fmt.Println(ta.Cfg.LoginSecretKey, " key")

	tokenStr, err := jwt.ParseWithClaims(token, jwtClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(ta.Cfg.LoginSecretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Error(ctx, err, "parse with claims failed")

			return jwtClaims, error2.InvalidSignature{}
		}

		log.Error(ctx, err, "parse with claims failed")
		return
	}

	if !tokenStr.Valid {
		log.Debug(ctx, "invalid token", jwtClaims.UserID)
		return jwtClaims, error2.UnauthorizedTokenError{}
	}

	if time.Now().After(time.Unix(jwtClaims.ExpireAt, 0)) {
		log.Debug(ctx, "expired token", jwtClaims.UserID)
		return jwtClaims, error2.ExpiredTokenError{}
	}
	log.Trace(ctx, "token validation is successful", jwtClaims.UserID)
	return jwtClaims, nil
}

func (ta *tokenAdapter) ValidateGeneratedToken(ctx context.Context, token string) (jwtClaims *entities.TokenGenRequest, err error) {

	jwtClaims = &entities.TokenGenRequest{}

	tokenStr, err := jwt.ParseWithClaims(token, jwtClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(ta.Cfg.TokenSecretKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Error(ctx, err, "parse with claims failed")

			return jwtClaims, error2.InvalidSignature{}
		}

		log.Error(ctx, err, "parse with claims failed")
		return
	}

	if !tokenStr.Valid {
		log.Debug(ctx, "invalid token", jwtClaims.IssuedAt)
		return jwtClaims, error2.UnauthorizedTokenError{}
	}

	if time.Now().After(time.Unix(jwtClaims.ExpireAt, 0)) {
		log.Debug(ctx, "expired token", jwtClaims.ExpireAt)
		return jwtClaims, error2.ExpiredTokenError{}
	}

	log.Trace(ctx, "token validation is successful")

	//check if the token
	return jwtClaims, nil
}

func (ta *tokenAdapter) DecodeAuthToken(ctx context.Context, token string) (jwtClaims *entities.JWTClaims, err error) {
	jwtSplit := strings.Split(token, ".")
	if len(jwtSplit) != 3 {
		log.Error(ctx, "", "token format error ")
		err = errors.New("invalid token")
		return
	}
	var data []byte
	data, err = base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(jwtSplit[1])
	if err != nil {
		log.Error(ctx, "Auth Token Base64 decode Error :", data, err)
		return
	}
	err = json.Unmarshal(data, &jwtClaims)
	if err != nil {
		log.Error(ctx, "Auth Token Unmarshal Error :", string(data), err)
		return
	}
	if time.Now().After(time.Unix(jwtClaims.ExpireAt, 0)) {
		log.Error(ctx, "expired token", token)
		err = errors.New("expired token")
		return
	}
	return
}

func (ta *tokenAdapter) DecodeGeneratedToken(ctx context.Context, token string) (jwtClaims *entities.TokenGenRequest, err error) {
	jwtSplit := strings.Split(token, ".")
	if len(jwtSplit) != 3 {
		log.Error(ctx, "", "token format error ")
		err = errors.New("invalid token")
		return
	}
	var data []byte
	data, err = base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(jwtSplit[1])
	if err != nil {
		log.Error(ctx, "Auth Token Base64 decode Error :", data, err)
		return
	}
	err = json.Unmarshal(data, &jwtClaims)
	if err != nil {
		log.Error(ctx, "Auth Token Unmarshal Error :", string(data), err)
		return
	}
	if time.Now().After(time.Unix(jwtClaims.ExpireAt, 0)) {
		log.Error(ctx, "expired token", token)
		err = errors.New("expired token")
		return
	}
	return
}
