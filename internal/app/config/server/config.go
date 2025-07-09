package server

import (
	"flag"

	"net/url"
	"strings"

	"github.com/caarlos0/env/v6"
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type serverConfig struct {
	Host        *string `env:"SERVER_ADDRESS"`
	BaseURL     *string `env:"BASE_URL"`
	EnableHTTPS *bool   `env:"ENABLE_HTTPS"`
}

// ServerConfigFromFlags parses command-line flags to configure server settings.
// It defines flags for host address and base URL with default values, and returns a serverConfig.
// The host flag defaults to "localhost:8080" and the base URL flag defaults to "http://localhost:8080".
func ServerConfigFromFlags() *serverConfig {
	host := flag.String("a", "localhost:8080", "host:port")

	baseURL := flag.String("b", "http://localhost:8080", "base url")

	enableHTTPs := flag.Bool("s", false, "enable HTTPs")

	return &serverConfig{
		Host:        host,
		BaseURL:     baseURL,
		EnableHTTPS: enableHTTPs,
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

// MergeServerConfigs merges two serverConfig objects, prioritizing non-nil values from envConfig.
// Logs a fatal error if either config is nil.
// Returns the merged serverConfig.
func MergeServerConfigs(envConfig, flagsConfig *serverConfig, log logger) *serverConfig {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if envConfig.Host == nil {
		envConfig.Host = flagsConfig.Host
	}

	if envConfig.BaseURL == nil {
		envConfig.BaseURL = flagsConfig.BaseURL
	}

	if envConfig.EnableHTTPS == nil {
		envConfig.EnableHTTPS = flagsConfig.EnableHTTPS
	}
	return envConfig
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
