package container

import (
	"github.com/charmLd/token-generator-api/util/config"
	"time"
)

// Resolve resolves the entire container.
// The order of resolution is very important. Low level dependencies need to be resolved before high level dependencies.
// It generally happens in this order.
// 		- Adapters
// 		- Repositories
// 		- Services
func Resolve(cfg *config.Config, location *time.Location) *Container {

	return &Container{
		Adapters:      resolveAdapters(cfg),
		Repositories:  resolveRepositories(location),
		Configs:       resolveConfigs(cfg),
		Location:      location,
		ThrottledTime: cfg.TokenConf.ThrottledTime,
	}
}
