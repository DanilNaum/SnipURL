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
	Host    *string `env:"SERVER_ADDRESS"`
	BaseURL *string `env:"BASE_URL"`
}

func ServerConfigFromFlags() *serverConfig {
	host := flag.String("a", "localhost:8080", "host:port")

	baseURL := flag.String("b", "http://localhost:8080", "base url")

	return &serverConfig{
		Host:    host,
		BaseURL: baseURL,
	}
}

func ServerConfigFromEnv(log logger) *serverConfig {
	c := &serverConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

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

	return envConfig
}

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

func (c *serverConfig) HTTPServerHost() string {
	return *c.Host
}

func (c *serverConfig) GetBaseURL() string {
	return *c.BaseURL
}

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
