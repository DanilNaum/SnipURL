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
	Path string `env:"FILE_STORAGE_PATH"`
}

func NewDumpConfig(log logger) *dumpConfig {
	envConfig := configFromEnv(log)
	flagsConfig := configFromFlags()

	c := mergeConfigs(envConfig, flagsConfig, log)

	return c
}

func configFromFlags() *dumpConfig {
	path := flag.String("f", "../../dump/storage.json", "path to dump file")
	flag.Parse()
	return &dumpConfig{
		Path: *path,
	}
}

func configFromEnv(log logger) *dumpConfig {
	c := &dumpConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

func mergeConfigs(envConfig, flagsConfig *dumpConfig, log logger) *dumpConfig {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if envConfig.Path == "" {
		envConfig.Path = flagsConfig.Path
	}

	return envConfig
}

func (c *dumpConfig) GetPath() string {
	return c.Path
}
