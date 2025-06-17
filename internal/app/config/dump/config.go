package dump

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type dumpConfig struct {
	Path *string `env:"FILE_STORAGE_PATH"`
}

func DumpConfigFromFlags() *dumpConfig {
	path := flag.String("f", "storage.json", "path to dump file")
	return &dumpConfig{
		Path: path,
	}
}

func DumpConfigFromEnv(log logger) *dumpConfig {
	c := &dumpConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

func MergeDumpConfigs(envConfig, flagsConfig *dumpConfig, log logger) *dumpConfig {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if envConfig.Path == nil {
		envConfig.Path = flagsConfig.Path
	}

	return envConfig
}

func (c *dumpConfig) GetPath() string {
	return *c.Path
}
