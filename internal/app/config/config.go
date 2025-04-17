package config

import (
	"flag"

	"github.com/DanilNaum/SnipURL/internal/app/config/cookie"
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
type cookieConfig interface {
	GetSecret() string
}

type config struct {
	serverConfig serverConfig
	dumpConfig   dumpConfig
	dbConfig     dbConfig
	cookieConfig cookieConfig
}

func NewConfig(log logger) *config {
	dbConfigFlag := db.DBConfigFromFlags()
	dumpConfigFlags := dump.DumpConfigFromFlags()
	serverConfigFlags := server.ServerConfigFromFlags()

	flag.Parse()

	dbConfigEnv := db.DBConfigFromEnv(log)
	dumpConfigEnv := dump.DumpConfigFromEnv(log)
	serverConfigEnv := server.ServerConfigFromEnv(log)
	cookieConfigEnv := cookie.CookieConfigFromEnv(log)

	serverConfig := server.MergeServerConfigs(serverConfigEnv, serverConfigFlags, log)
	dumpConfig := dump.MergeDumpConfigs(dumpConfigEnv, dumpConfigFlags, log)
	dbConfig := db.MergeDBConfigs(dbConfigEnv, dbConfigFlag, log)

	return &config{
		serverConfig: serverConfig,
		dumpConfig:   dumpConfig,
		dbConfig:     dbConfig,
		cookieConfig: cookieConfigEnv,
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

func (c *config) CookieConfig() cookieConfig {
	return c.cookieConfig
}
