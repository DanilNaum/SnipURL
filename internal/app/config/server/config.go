package server

import (
	"flag"

	"net/url"
	"strings"

	"github.com/DanilNaum/SnipURL/internal/app/config/utils"
	"github.com/caarlos0/env/v6"
)

var (
	defaultHost        = "localhost:8080"
	defaultBaseURL     = "http://localhost:8080"
	defaultEnableHTTPS = false
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type serverConfig struct {
	Host        *string `json:"server_address" env:"SERVER_ADDRESS"`
	BaseURL     *string `json:"base_url" env:"BASE_URL"`
	EnableHTTPS *bool   `json:"enable_https" env:"ENABLE_HTTPS"`
}

// ServerConfigFromFlags parses command-line flags to configure server settings.
// It defines flags for host address and base URL with default values, and returns a serverConfig.
// The host flag defaults to "localhost:8080" and the base URL flag defaults to "http://localhost:8080".
func ServerConfigFromFlags() *serverConfig {
	host := flag.String("a", "", "host:port")

	baseURL := flag.String("b", "", "base url")

	enableHTTPS := flag.Bool("s", false, "enable HTTPS")

	return &serverConfig{
		Host:        host,
		BaseURL:     baseURL,
		EnableHTTPS: enableHTTPS,
	}
}

// ServerConfigFromEnv parses server configuration from environment variables.
// It uses the env package to populate a serverConfig struct with values from the environment.
// If parsing fails, it logs a fatal error with the parsing error details.
// Returns the parsed server configuration or nil if parsing fails.
func ServerConfigFromEnv(log logger) *serverConfig {
	c := &serverConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

// ServerConfigFromJsonFile reads a JSON configuration file and returns a pointer
// to a serverConfig object. It takes the name of the JSON file and a logger
// instance as parameters.
func ServerConfigFromJsonFile(jsonFileName string, log logger) *serverConfig {
	var config serverConfig
	if jsonFileName != "" {
		if err := utils.LoadConfigFromFile(jsonFileName, &config); err != nil {
			log.Fatalf(err.Error())
		}
	}
	return &config
}

// MergeServerConfigs merges  serverConfig objects, prioritizing non-nil values from envConfig.
// Logs a fatal error if either config is nil.
// Returns the merged serverConfig.
func MergeServerConfigs(envConfig, flagsConfig, fileConfig *serverConfig, log logger) *serverConfig {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if *flagsConfig.Host == "" {
		flagsConfig.Host = nil
	}

	if *flagsConfig.BaseURL == "" {
		flagsConfig.BaseURL = nil
	}

	if *flagsConfig.EnableHTTPS == false {
		flagsConfig.EnableHTTPS = nil
	}

	if fileConfig == nil {
		return &serverConfig{
			Host:        utils.Merge(envConfig.Host, flagsConfig.Host, &defaultHost),
			BaseURL:     utils.Merge(envConfig.BaseURL, flagsConfig.BaseURL, &defaultBaseURL),
			EnableHTTPS: utils.Merge(envConfig.EnableHTTPS, flagsConfig.EnableHTTPS, &defaultEnableHTTPS),
		}
	}
	return &serverConfig{
		Host:        utils.Merge(envConfig.Host, flagsConfig.Host, fileConfig.Host, &defaultHost),
		BaseURL:     utils.Merge(envConfig.BaseURL, flagsConfig.BaseURL, fileConfig.BaseURL, &defaultBaseURL),
		EnableHTTPS: utils.Merge(envConfig.EnableHTTPS, flagsConfig.EnableHTTPS, fileConfig.EnableHTTPS, &defaultEnableHTTPS),
	}
}

// ValidateServerConfig checks the server configuration for valid host and base URL values.
// It ensures the host and base URL are properly formatted, the base URL uses http or https,
// and that the base URL contains the host. Logs a fatal error if any validation fails.
func (c *serverConfig) ValidateServerConfig(log logger) {

	_, err := url.Parse(*c.Host)
	if err != nil {
		log.Fatalf("invalid host: %s", *c.Host)
		return
	}

	u, err := url.Parse(*c.BaseURL)
	if err != nil {
		log.Fatalf("invalid base url: %s", *c.BaseURL)
		return
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		log.Fatalf("invalid base url: %s", *c.BaseURL)
		return
	}

	if !strings.Contains(*c.BaseURL, *c.Host) {
		log.Fatalf("base url %s must contain host %s", *c.BaseURL, *c.Host)
	}

}

// HTTPServerHost returns the HTTP server host address from the serverConfig.
// It dereferences the Host field pointer and returns its value.
func (c *serverConfig) HTTPServerHost() string {
	return *c.Host
}

// GetBaseURL returns the base URL from the serverConfig.
// It dereferences the BaseURL field pointer and returns its value.
func (c *serverConfig) GetBaseURL() string {
	return *c.BaseURL
}

// GetEnableHTTPs returns the value of the EnableHTTPs field indicating if HTTPS is enabled.
func (c *serverConfig) GetEnableHTTPS() bool {
	return *c.EnableHTTPS
}

// GetPrefix extracts and returns the path prefix from the base URL.
// If the base URL has no path, it returns "/". If parsing fails, it returns an error.
// The returned path is trimmed of any trailing slash.
func (c *serverConfig) GetPrefix() (string, error) {
	parsedURL, err := url.Parse(*c.BaseURL)
	if err != nil {
		return "", err
	}

	path := parsedURL.Path

	path = strings.TrimSuffix(path, "/")

	if path == "" {
		return "/", nil
	}

	return path, nil
}
