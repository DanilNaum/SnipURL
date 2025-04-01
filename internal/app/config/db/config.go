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

func DBConfigFromFlags() *dbConfig {
	dsn := flag.String("d", "", "dsn")

	return &dbConfig{
		DBDSN: dsn,
	}
}

func DBConfigFromEnv(log logger) *dbConfig {
	c := &dbConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

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

func (c *dbConfig) GetDSN() string {
	return *c.DBDSN
}
