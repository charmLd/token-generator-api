package repositories

import (
	"context"
	"github.com/charmLd/token-generator-api/domain/entities"
)

type TokenRepositoryInterface interface {
	InsertUniqueToken(ctx context.Context, token entities.Token) (err error)
	Revoke(ctx context.Context, userId string) (err error)
	GetAllTokenForFilter(ctx context.Context, fetchDetailsFilters entities.TokenDetailsReqParam) (tokenDetailsArray []entities.Token, err error)
	CreateNewToken(ctx context.Context, token entities.Token) (err error)
	FetchTokenInfo(ctx context.Context, tokenDetails entities.ValidateRequest) (tokenDetailsArray entities.Token, err error)
}
