package cookie

import "github.com/caarlos0/env/v6"

type logger interface {
	Fatalf(format string, v ...any)
}

type cookieConfig struct {
	Secret string `env:"COOKIE_SECRET" envDefault:"secret1234567890"`
}

func CookieConfigFromEnv(log logger) *cookieConfig {
	c := &cookieConfig{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("error parse config from Env: %s", err)
	}
	return c
}

func (c *cookieConfig) GetSecret() string {
	return c.Secret
}
