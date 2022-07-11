package config

import ()

// Config holds all other config structs.
type Config struct {
	AppConf   AppConfig
	DBConf    DBConfig
	TokenConf TokenConfig
}

// AppConfig holds application configurations.
type AppConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	InternalPort     int    `yaml:"internal-port"`
	PublicPort     int    `yaml:"public-port"`
	Timezone string `yaml:"timezone"`
}

// PostgresConfig holds postgresql database configurations.
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	MaxOpenCon int  `yaml:"max_open_connection"`
	MaxIdleCon int `yaml:"max_idle_connection"`
}

type TokenConfig struct {
	LoginTokenExpiry         int `yaml:"login-token-expiry" `
	LoginSecretKey        string  `yaml:"login-secret-key" `

	TokenExpiry int  `yaml:"token-expiry" `
	TokenSecretKey            string  `yaml:"token-secret-key" `

	HashAlgorithm            string  `yaml:"hash-algorithm " `
	ThrottledTime int `yaml:"throttled-time"`
}
