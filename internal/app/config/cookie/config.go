package cookie

import "github.com/caarlos0/env/v6"

type logger interface {
	Fatalf(format string, v ...any)
}

type cookieConfig struct {
	Secret string `env:"COOKIE_SECRET" envDefault:"secret1234567890"`
}

// CookieConfigFromEnv parses cookie configuration from environment variables.
// It uses the env package to load configuration and logs a fatal error if parsing fails.
// Returns a configured cookieConfig with default or environment-specified values.
func CookieConfigFromEnv(log logger) *cookieConfig {
	c := &cookieConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

// GetSecret returns the cookie secret used for configuration.
func (c *cookieConfig) GetSecret() string {
	return c.Secret
}
