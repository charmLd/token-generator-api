package mock_repositories

import (
	"context"
	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
	"github.com/charmLd/token-generator-api/domain/boundary/repositories"
	"github.com/charmLd/token-generator-api/domain/entities"
	"time"
)

type MockUserRepository struct {
	DBAdapter   adapters.DBAdapterInterface
	transaction adapters.TransactionInterface
	Location    *time.Location
}

func NewMockUserRepository(dbAdapter adapters.DBAdapterInterface, loc *time.Location) repositories.UserRepositoryInterface {

	return &MockUserRepository{
		DBAdapter: dbAdapter,
		Location:  loc,
	}
}
func (usr *MockUserRepository) UpdateLastLogin(ctx context.Context, userID string) (err error) {
	if userID == "1" {
		return nil
	}
	return
}
func (usr *MockUserRepository) GetLastLoginTime(ctx context.Context, userID string) (lastlogin time.Time, err error) {
	return
}
func (usr *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (user entities.User, err error) {
	return
}
