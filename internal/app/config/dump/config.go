package dump

import (
	"flag"

	"github.com/DanilNaum/SnipURL/internal/app/config/utils"
	"github.com/caarlos0/env/v6"
)

var defaultVal = "storage.json"

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type dumpConfig struct {
	Path *string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`
}

// DumpConfigFromFlags creates a dumpConfig with a default storage path from command-line flags.
// It sets the default dump file path to "storage.json" if not specified.
func DumpConfigFromFlags() *dumpConfig {
	path := flag.String("f", "", "path to dump file")
	return &dumpConfig{
		Path: path,
	}
}

// DumpConfigFromEnv parses environment variables to configure the dump configuration.
// It logs a fatal error if parsing the environment configuration fails.
func DumpConfigFromEnv(log logger) *dumpConfig {
	c := &dumpConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

// DumpConfigFromJSONFile loads configuration from a JSON file into a dumpConfig struct.
// Logs a fatal error and exits if loading fails.
func DumpConfigFromJSONFile(jsonFileName string, log logger) *dumpConfig {
	var config dumpConfig
	if jsonFileName != "" {
		if err := utils.LoadConfigFromFile(jsonFileName, &config); err != nil {
			log.Fatalf(err.Error())
		}
	}
	return &config
}

// MergeDumpConfigs combines environment and flag-based configurations for dump settings.
// It prioritizes environment configuration and uses flag configuration as a fallback.
// Logs a fatal error if either configuration is nil.
func MergeDumpConfigs(envConfig, flagsConfig, fileConfig *dumpConfig, log logger) *dumpConfig {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if *flagsConfig.Path == "" {
		flagsConfig.Path = nil
	}

	if fileConfig == nil {
		return &dumpConfig{
			Path: utils.Merge(envConfig.Path, flagsConfig.Path, &defaultVal),
		}

	}

	return &dumpConfig{
		Path: utils.Merge(envConfig.Path, flagsConfig.Path, fileConfig.Path, &defaultVal),
	}
}

// GetPath returns the configured file storage path.
// Returns an empty string if no path is set.
func (c *dumpConfig) GetPath() string {
	return *c.Path
}
