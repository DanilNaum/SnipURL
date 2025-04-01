package config

import (
	"flag"

	"github.com/DanilNaum/SnipURL/internal/app/config/db"
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

type dbConfig interface {
	GetDSN() string
}

type config struct {
	serverConfig serverConfig
	dumpConfig   dumpConfig
	dbConfig     dbConfig
}

func NewConfig(log logger) *config {
	dbConfigFlag := db.DBConfigFromFlags()
	dumpConfigFlags := dump.DumpConfigFromFlags()
	serverConfigFlags := server.ServerConfigFromFlags()

	flag.Parse()

	dbConfigEnv := db.DBConfigFromEnv(log)
	dumpConfigEnv := dump.DumpConfigFromEnv(log)
	serverConfigEnv := server.ServerConfigFromEnv(log)

	serverConfig := server.MergeServerConfigs(serverConfigEnv, serverConfigFlags, log)
	dumpConfig := dump.MergeDumpConfigs(dumpConfigEnv, dumpConfigFlags, log)
	dbConfig := db.MergeDBConfigs(dbConfigEnv, dbConfigFlag, log)

	return &config{
		serverConfig: serverConfig,
		dumpConfig:   dumpConfig,
		dbConfig:     dbConfig,
	}
}

func (c *config) ServerConfig() serverConfig {
	return c.serverConfig
}

func (c *config) DumpConfig() dumpConfig {
	return c.dumpConfig
}

func (c *config) DBConfig() dbConfig {
	return c.dbConfig
}
