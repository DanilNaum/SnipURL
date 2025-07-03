package db

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type dbConfig struct {
	DBDSN *string `env:"DATABASE_DSN"`
}

// DBConfigFromFlags creates a database configuration from command-line flags.
// It returns a dbConfig with a DSN (Data Source Name) parsed from the -d flag.
func DBConfigFromFlags() *dbConfig {
	dsn := flag.String("d", "", "dsn")

	return &dbConfig{
		DBDSN: dsn,
	}
}

// DBConfigFromEnv creates a database configuration by parsing environment variables.
// It uses the env package to populate the configuration and logs a fatal error if parsing fails.
func DBConfigFromEnv(log logger) *dbConfig {
	c := &dbConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

// MergeDBConfigs combines database configurations from environment and flags.
// If either configuration is nil, it logs a fatal error.
// It prioritizes the environment configuration, using flag configuration as a fallback for DSN.
func MergeDBConfigs(envConfig, flagsConfig *dbConfig, log logger) *dbConfig {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if envConfig.DBDSN == nil {
		envConfig.DBDSN = flagsConfig.DBDSN
	}

	return envConfig
}

// GetDSN returns the Data Source Name (DSN) from the database configuration.
// It returns an empty string if no DSN is set.
func (c *dbConfig) GetDSN() string {
	return *c.DBDSN
}
