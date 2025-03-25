package config

import (
	"github.com/DanilNaum/SnipURL/internal/app/config/dump"
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

type dumpConfig interface {
	GetPath() string
}

type config struct {
	serverConfig serverConfig
	dumpConfig   dumpConfig
}

func NewConfig(log logger) *config {
	return &config{
		serverConfig: server.NewConfig(log),
		dumpConfig:   dump.NewDumpConfig(log),
	}
}

func (c *config) ServerConfig() serverConfig {
	return c.serverConfig
}

func (c *config) DumpConfig() dumpConfig {
	return c.dumpConfig
}
