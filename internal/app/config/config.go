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

type config struct {
	serverConfig serverConfig
}

func NewConfig(log logger) *config {
	return &config{
		serverConfig: server.NewConfig(log),
	}
}

func (c *config) ServerConfig() serverConfig {
	return c.serverConfig
}
