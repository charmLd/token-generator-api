package container

import (
	"time"

	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
	"github.com/charmLd/token-generator-api/domain/boundary/repositories"
	"github.com/charmLd/token-generator-api/util/config"
)

type Container struct {
	Configs       Configs
	Adapters      Adapters
	Repositories  Repositories
	Location      *time.Location
	ThrottledTime int
}

type Adapters struct {
	MySQL adapters.DBAdapterInterface
	Token adapters.TokenAdapterInterface
}

type Repositories struct {
	UserRepository repositories.UserRepositoryInterface

	TokenRepository repositories.TokenRepositoryInterface
}

type Configs struct {
	TokenConfig config.TokenConfig
	DBConfig    config.DBConfig
	AppConfig   config.AppConfig
}
