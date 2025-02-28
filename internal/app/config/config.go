package config

import (
	"github.com/DanilNaum/SnipURL/internal/app/config/server"
)

type logger interface {
	Fatalf(format string, v ...any)
}

type serverConfig interface {
	HTTPServerHost() string
	BaseURL() string
	GetPrefix() (string, error)
}

type Config struct {
	serverConfig serverConfig
}

func GetConfig(log logger) *Config {
	return &Config{
		serverConfig: server.NewConfigFromFlags(log),
	}
}

func (c *Config) ServerConfig() serverConfig {
	return c.serverConfig
}
