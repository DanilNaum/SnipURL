package server

import (
	"flag"

	"net/url"
	"strings"

	"regexp"
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type config struct {
	host    string
	baseURL string
}

func NewConfigFromFlags(log logger) *config {
	host := flag.String("a", "localhost:8080", "host:port")

	baseURL := flag.String("b", "http://localhost:8080", "base url")

	flag.Parse()

	c := &config{
		host:    *host,
		baseURL: *baseURL,
	}

	c.validate(log)

	return c
}

func (c *config) validate(log logger) {
	hostPattern := `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])(:[0-9]{1,5})?$`
	re := regexp.MustCompile(hostPattern)
	if !re.MatchString(c.host) {
		log.Fatalf("invalid host: %s", c.host)
		return
	}

	baseURLPattern := `^(http|https):\/\/(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])(:[0-9]{1,5})?(\/.*)?$`
	re = regexp.MustCompile(baseURLPattern)
	if !re.MatchString(c.baseURL) {
		log.Fatalf("invalid base url: %s", c.baseURL)
		return
	}

	if !strings.Contains(c.baseURL, c.host) {
		log.Fatalf("base url %s must contain host %s", c.baseURL, c.host)
	}

}

func (c *config) HTTPServerHost() string {
	return c.host
}

func (c *config) BaseURL() string {
	return c.baseURL
}

func (c *config) GetPrefix() (string, error) {
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}
	path := parsedURL.Path
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		return "/", nil
	}
	// Возвращаем путь
	return path, nil
}
