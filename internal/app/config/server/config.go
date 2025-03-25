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

type config struct {
	Host    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func NewConfig(log logger) *config {

	envConfig := configFromEnv(log)
	flagsConfig := configFromFlags()

	c := mergeConfigs(envConfig, flagsConfig, log)

	c.validate(log)

	return c
}

func configFromFlags() *config {
	host := flag.String("a", "localhost:8080", "host:port")

	baseURL := flag.String("b", "http://localhost:8080", "base url")

	flag.Parse()
	return &config{
		Host:    *host,
		BaseURL: *baseURL,
	}
}

func configFromEnv(log logger) *config {
	c := &config{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

func mergeConfigs(envConfig, flagsConfig *config, log logger) *config {
	if envConfig == nil {
		log.Fatalf("error env config is nil")
		return nil
	}

	if flagsConfig == nil {
		log.Fatalf("error flags config is nil")
		return nil
	}

	if envConfig.Host == "" {
		envConfig.Host = flagsConfig.Host
	}

	if envConfig.BaseURL == "" {
		envConfig.BaseURL = flagsConfig.BaseURL
	}

	return envConfig
}

func (c *config) validate(log logger) {

	_, err := url.Parse(c.Host)
	if err != nil {
		log.Fatalf("invalid host: %s", c.Host)
		return
	}

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		log.Fatalf("invalid base url: %s", c.BaseURL)
		return
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		log.Fatalf("invalid base url: %s", c.BaseURL)
		return
	}

	if !strings.Contains(c.BaseURL, c.Host) {
		log.Fatalf("base url %s must contain host %s", c.BaseURL, c.Host)
	}

}

func (c *config) HTTPServerHost() string {
	return c.Host
}

func (c *config) GetBaseURL() string {
	return c.BaseURL
}

func (c *config) GetPrefix() (string, error) {
	parsedURL, err := url.Parse(c.BaseURL)
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
