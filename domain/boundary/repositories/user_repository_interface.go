package repositories

import (
	"context"
	"time"

	"github.com/charmLd/token-generator-api/domain/entities"
)

type UserRepositoryInterface interface {
	UpdateLastLogin(ctx context.Context, userID string) (err error)
	GetLastLoginTime(ctx context.Context, userID string) (lastlogin time.Time, err error)
	GetUserByEmail(ctx context.Context, email string) (user entities.User, err error)
}
