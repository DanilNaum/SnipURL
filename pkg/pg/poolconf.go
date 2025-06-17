package pg

import (
	"fmt"
)

type connConfig struct {
	host      string
	port      string
	username  string
	password  string
	dbName    string
	dncString string
}

func NewConnConfig(host, port, username, password, dbName string) *connConfig {
	return &connConfig{
		host:     host,
		port:     port,
		username: username,
		password: password,
		dbName:   dbName,
	}
}

func NewConnConfigFromDsnString(dnsString string) *connConfig {
	return &connConfig{
		dncString: dnsString,
	}
}

func (c *connConfig) getDsn() string {
	if c.dncString != "" {
		return c.dncString
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.host,
		c.port,
		c.username,
		c.password,
		c.dbName,
	)
}
