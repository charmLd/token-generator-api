package usecases

import (
	"context"
	"fmt"
	appErr "github.com/charmLd/token-generator-api/domain/error"
	"github.com/charmLd/token-generator-api/mocks/mock_adapters"
	"github.com/charmLd/token-generator-api/mocks/mock_repositories"
	"github.com/charmLd/token-generator-api/util/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//todo complete all test cases
func TestGenerateToken(t *testing.T) {
	ctx := context.Background()

	loc := time.Location{}
	db := mock_adapters.MockAdapter{}
	cfgToken := config.TokenConfig{}
	authUseCase := &AuthUseCase{
		UserRepository:  mock_repositories.NewMockUserRepository(&db, &loc),
		TokenRepository: mock_repositories.NewMockTokenRepository(&db, &loc),
		TokenAdapter:    mock_adapters.NewMockTokenAdapter(cfgToken),
	}
	testCases := []struct {
		userId string
		//expectedtoken string
		expectedError error
	}{
		/*{"1", "exampleJWT", nil},
		{"2", "", (&appErr.Error{}).New(ctx, appErr.DOMAIN, "token Issue Error",
			"CORE-2012",
			"could not issue a token")},*/

		{"1", nil},
		{"2", (&appErr.Error{}).New(ctx, appErr.DOMAIN, "token Issue Error",
			"CORE-2012",
			"could not issue a token")},
	}
	for _, test := range testCases {
		_, err := authUseCase.GenerateToken(ctx, test.userId)

		/*if assert.Equal(t, test.expectedtoken, token) {
			t.Logf(fmt.Sprintf("Response expected %v , received %v : success", test.expectedtoken, token))
		} else {
			t.Errorf(fmt.Sprintf("Response expected %v , received %v : failed", test.expectedtoken, token))
		}*/
		/*if err != nil {
			t.Errorf(fmt.Sprintf("Response expected %v , received %v : failed", test.expected, err))
		}*/
		if assert.Equal(t, test.expectedError, err) {
			t.Logf(fmt.Sprintf("Response expected %v , received %v : success", test.expectedError, err))
		} else {
			t.Errorf(fmt.Sprintf("Response expected %v , received %v : failed", test.expectedError, err))
		}
	}
}
