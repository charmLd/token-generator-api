package controllers

import (
	"github.com/charmLd/token-generator-api/domain/usecases"

	"github.com/charmLd/token-generator-api/util/container"
)

// BaseController contains controller logic for endpoints
type BaseController struct {
	Container   *container.Container
	AuthUseCase usecases.AuthUseCase
}

// NewBaseController returns a base type for this controller
func NewBaseController(container *container.Container) *BaseController {

	authUseCase := usecases.AuthUseCase{
		Config: usecases.AuthConfig{
			//RTExpiry: container.Configs.TokenConfig.RTExpiry * 60 * 60 * 24, //converting days into seconds
			LoginTokenExpiry:     container.Configs.TokenConfig.LoginTokenExpiry * 60 * 60 * 24, //converting days into seconds
			GeneratedTokenExpiry: container.Configs.TokenConfig.TokenExpiry * 60 * 60 * 24,      //converting days into seconds

			//GracePeriod: container.Configs.TokenConfig.GracePeriod,
		},

		TokenAdapter:    container.Adapters.Token,
		TokenRepository: container.Repositories.TokenRepository,
		UserRepository:  container.Repositories.UserRepository,
	}

	return &BaseController{
		Container:   container,
		AuthUseCase: authUseCase,
	}
}
