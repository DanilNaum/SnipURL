package config

import (
	"github.com/DanilNaum/SnipURL/internal/app/config/server"
)

type logger interface {
	Fatalf(format string, v ...any)
}

type serverConfig interface {
	HTTPServerHost() string
	GetBaseURL() string
	GetPrefix() (string, error)
}

type Config struct {
	serverConfig serverConfig
}

func NewConfig(log logger) *Config {
	return &Config{
		serverConfig: server.NewConfig(log),
	}
}

func (c *Config) ServerConfig() serverConfig {
	return c.serverConfig
}
