package db

import (
	"flag"

	"github.com/DanilNaum/SnipURL/internal/app/config/utils"
	"github.com/caarlos0/env/v6"
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type dbConfig struct {
	DBDSN *string `json:"database_dsn" env:"DATABASE_DSN"`
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

// DBConfigFromJSONFile loads database configuration from a JSON file.
// It takes the path to a JSON file and a logger, loads the configuration into a dbConfig struct,
// and logs a fatal error if loading fails. Returns the loaded dbConfig.
func DBConfigFromJSONFile(jsonFileName string, log logger) *dbConfig {
	var config dbConfig
	if jsonFileName != "" {
		if err := utils.LoadConfigFromFile(jsonFileName, &config); err != nil {
			log.Fatalf(err.Error())
		}
	}
	return &config
}

// MergeDBConfigs combines database configurations from environment and flags.
// If either configuration is nil, it logs a fatal error.
// It prioritizes the environment configuration, using flag configuration as a fallback for DSN.
func MergeDBConfigs(envConfig, flagsConfig, fileConfig *dbConfig, log logger) *dbConfig {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if *flagsConfig.DBDSN == "" {
		flagsConfig.DBDSN = nil
	}

	if fileConfig == nil {
		return &dbConfig{
			DBDSN: utils.Merge(envConfig.DBDSN, flagsConfig.DBDSN),
		}
	}

	return &dbConfig{
		DBDSN: utils.Merge(envConfig.DBDSN, flagsConfig.DBDSN, fileConfig.DBDSN),
	}
}

// GetDSN returns the Data Source Name (DSN) from the database configuration.
// It returns an empty string if no DSN is set.
func (c *dbConfig) GetDSN() string {
	if c.DBDSN == nil {
		return ""
	}
	return *c.DBDSN
}
