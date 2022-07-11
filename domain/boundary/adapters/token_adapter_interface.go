package adapters

import (
	"context"

	"github.com/charmLd/token-generator-api/domain/entities"
)

// JwtAdapterInterface defines an interface for jwt provider
type TokenAdapterInterface interface {
	GenerateUniqueToken(ctx context.Context, jwtClaims entities.TokenGenRequest) (jwToken string, err error)
	GenerateLoginToken(ctx context.Context, jwtClaims entities.JWTClaims) (jwToken string, err error)
	ValidateGeneratedToken(ctx context.Context, token string) (jwtClaims *entities.TokenGenRequest, err error)
	ValidateLoginJWToken(ctx context.Context, token string) (jwtClaims *entities.JWTClaims, err error)
	DecodeAuthToken(ctx context.Context, token string) (jwtClaims *entities.JWTClaims, err error)
	DecodeGeneratedToken(ctx context.Context, token string) (jwtClaims *entities.TokenGenRequest, err error)
}
