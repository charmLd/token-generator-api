package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Parse parses all configuration to a single Config object.
func Parse() *Config {

	return &Config{
		AppConf:   parseAppConfig(),
		DBConf:    parseDBConfig(),
		TokenConf: parseTokenConfig(),
	
	}
}

// Parse application configurations.
func parseAppConfig() AppConfig {

	content := read("app.yaml")

	cfg := AppConfig{}

	err := yaml.Unmarshal(content, &cfg)
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return cfg
}

// Parse key value store configurations.
func parseTokenConfig() TokenConfig {

	content := read("token.yaml")

	cfg := TokenConfig{}

	err := yaml.Unmarshal(content, &cfg)
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return cfg
}

// Parse database configurations.
func parseDBConfig() DBConfig {

	content := read("database.yaml")

	cfg := DBConfig{}

	err := yaml.Unmarshal(content, &cfg)
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return cfg
}
