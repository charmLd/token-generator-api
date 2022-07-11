package mock_repositories

import (
	"context"
	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
	"github.com/charmLd/token-generator-api/domain/boundary/repositories"
	"github.com/charmLd/token-generator-api/domain/entities"
	appErr "github.com/charmLd/token-generator-api/domain/error"
	"time"
)

type MockTokenRepository struct {
	DBAdapter   adapters.DBAdapterInterface
	transaction adapters.TransactionInterface
	Location    *time.Location
}

func NewMockTokenRepository(dbAdapter adapters.DBAdapterInterface, loc *time.Location) repositories.TokenRepositoryInterface {

	return &MockTokenRepository{
		DBAdapter: dbAdapter,
		Location:  loc,
	}
}
func (tr *MockTokenRepository) Revoke(ctx context.Context, userId string) (err error) {
	return
}
func (tr *MockTokenRepository) InsertUniqueToken(ctx context.Context, token entities.Token) (err error) {
	if token.UserId == "1" {
		return nil
	}
	if token.UserId == "2" {
		return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "token Issue Error",
			"CORE-2012",
			"could not issue a token")
	}
	return
}

func (t *MockTokenRepository) GetAllTokenForFilter(ctx context.Context, fetchDetailsFilters entities.TokenDetailsReqParam) (tokenDetailsArray []entities.Token, err error) {
	return
}
func (tr *MockTokenRepository) CreateNewToken(ctx context.Context, token entities.Token) (err error) {
	return
}
func (t *MockTokenRepository) FetchTokenInfo(ctx context.Context, tokenDetails entities.ValidateRequest) (tokendetail entities.Token, err error) {
	return
}
