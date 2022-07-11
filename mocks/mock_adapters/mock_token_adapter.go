package mock_adapters

import (
	"context"
	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
	"github.com/charmLd/token-generator-api/domain/entities"
	"github.com/charmLd/token-generator-api/util/config"
)

type MockTokenAdapter struct {
	Cfg config.TokenConfig
}

func NewMockTokenAdapter(cfg config.TokenConfig) adapters.TokenAdapterInterface {

	return &MockTokenAdapter{
		Cfg: cfg,
	}
}
func (ta *MockTokenAdapter) GenerateUniqueToken(ctx context.Context, jwtClaims entities.TokenGenRequest) (jwToken string, err error) {

	return "exampleJWT", nil
}
func (ta *MockTokenAdapter) GenerateLoginToken(ctx context.Context, jwtClaims entities.JWTClaims) (jwToken string, err error) {
	return
}
func (ta *MockTokenAdapter) ValidateLoginJWToken(ctx context.Context, token string) (jwtClaims *entities.JWTClaims, err error) {
	return
}
func (ta *MockTokenAdapter) ValidateGeneratedToken(ctx context.Context, token string) (jwtClaims *entities.TokenGenRequest, err error) {
	return
}
func (ta *MockTokenAdapter) DecodeAuthToken(ctx context.Context, token string) (jwtClaims *entities.JWTClaims, err error) {
	return
}
func (ta *MockTokenAdapter) DecodeGeneratedToken(ctx context.Context, token string) (jwtClaims *entities.TokenGenRequest, err error) {
	return
}
