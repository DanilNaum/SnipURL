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

// NewConfig creates a new configuration by merging configuration values from flags, environment variables, and applying default settings.
// It takes a logger as a parameter to handle potential configuration errors.
// The function parses command-line flags and combines configurations for server, dump, database, and cookie settings.
// Returns a fully initialized config struct with merged configuration values.
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

// ServerConfig returns the server configuration for the current config instance.
// It provides access to the serverConfig field, which contains server-related settings.
func (c *config) ServerConfig() serverConfig {
	return c.serverConfig
}

// DumpConfig returns the dump configuration for the current config instance.
// It provides access to the dumpConfig field, which contains dump-related settings.
func (c *config) DumpConfig() dumpConfig {
	return c.dumpConfig
}

// DBConfig returns the database configuration for the current config instance.
// It provides access to the dbConfig field, which contains database-related settings.
func (c *config) DBConfig() dbConfig {
	return c.dbConfig
}

// CookieConfig returns the cookie configuration for the current config instance.
// It provides access to the cookieConfig field, which contains cookie-related settings.
func (c *config) CookieConfig() cookieConfig {
	return c.cookieConfig
}
