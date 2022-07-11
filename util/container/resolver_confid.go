package container

import "github.com/charmLd/token-generator-api/util/config"

var configs Configs

func resolveConfigs(cfg *config.Config) Configs {
	//resolveRoleConfig(cfg.Roles)
	resolveTokenConfig(cfg.TokenConf)
	resolveAppConfig(cfg.AppConf)
	resolveDBConfig(cfg.DBConf)
	return configs
}

func resolveTokenConfig(cfg config.TokenConfig) {
	configs.TokenConfig = cfg
}
func resolveDBConfig(cfg config.DBConfig) {
	configs.DBConfig = cfg
}
func resolveAppConfig(cfg config.AppConfig) {
	configs.AppConfig = cfg
}
