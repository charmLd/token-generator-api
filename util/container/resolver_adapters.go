package container

import (
	"github.com/charmLd/token-generator-api/externals/adapters"
	"github.com/charmLd/token-generator-api/externals/adapters/mysql"
	"github.com/charmLd/token-generator-api/util/config"
)

var resolvedAdapters Adapters

// Resolve all adapters.
func resolveAdapters(cfg *config.Config) Adapters {

	resolveDBAdapter(cfg.DBConf)
	resolveTokenAdapter(cfg.TokenConf)

	return resolvedAdapters
}

//Resolve the database adapter.
func resolveDBAdapter(cfg config.DBConfig) {
	pg := mysql.Adapter{}
	db, _ := pg.New(cfg)
	resolvedAdapters.MySQL = db
}

//Resolve JWT adapter
func resolveTokenAdapter(cfg config.TokenConfig) {
	//tokenAdapter := adapters.NewTokenAdapter(cfg)
	tokenAdapter := adapters.NewTokenAdapter(cfg)
	resolvedAdapters.Token = tokenAdapter
}
