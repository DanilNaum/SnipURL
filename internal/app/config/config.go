package config

import (
	"flag"
	"os"

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
	GetEnableHTTPS() bool
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

	var configFile string
	flag.StringVar(&configFile, "c", "", "config file name")
	flag.StringVar(&configFile, "config", "", "config file name (long form)")

	flag.Parse()

	dbConfigEnv := db.DBConfigFromEnv(log)
	dumpConfigEnv := dump.DumpConfigFromEnv(log)
	serverConfigEnv := server.ServerConfigFromEnv(log)
	cookieConfigEnv := cookie.CookieConfigFromEnv(log)

	if configFile == "" {
		configFile = os.Getenv("CONFIG")
	}

	dbConfigFile := db.DBConfigFromJSONFile(configFile, log)
	dumpConfigFile := dump.DumpConfigFromJSONFile(configFile, log)
	serverConfigFile := server.ServerConfigFromJSONFile(configFile, log)

	serverConfig := server.MergeServerConfigs(serverConfigEnv, serverConfigFlags, serverConfigFile, log)
	dumpConfig := dump.MergeDumpConfigs(dumpConfigEnv, dumpConfigFlags, dumpConfigFile, log)
	dbConfig := db.MergeDBConfigs(dbConfigEnv, dbConfigFlag, dbConfigFile, log)

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
