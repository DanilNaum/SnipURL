package httpserver

import (
	"fmt"
	"strings"
	"time"
)

// Option represents a configuration function for customizing a server's settings.
// It allows modifying server properties through functional options pattern.
type Option func(s *server)

// WithShutdownTimeout sets the graceful shutdown timeout for the HTTP server.
// The provided duration determines how long the server will wait for existing connections
// to complete before forcefully shutting down.
func WithShutdownTimeout(t time.Duration) Option {
	return func(s *server) {
		s.shutdownTimeout = t
	}
}

// WithPort configures the HTTP server to listen on the specified port.
// The port is set as part of the server's address, defaulting to the format ":<port>".
func WithPort(port int) Option {
	return func(s *server) {
		s.server.Addr = fmt.Sprintf(":%d", port)
	}
}

// WithAddr sets the server's listening address, extracting the port from the provided address.
// If no port is specified, it defaults to port 80. The address is parsed to ensure
// a valid port is used for the server configuration.
func WithAddr(addr string) Option {
	return func(s *server) {
		parts := strings.Split(addr, ":")
		if len(parts) < 2 {
			s.server.Addr = ":80" // Если порт не указан, возвращаем пустую строку
			return
		}
		s.server.Addr = ":" + parts[len(parts)-1]
	}
}

// WithTLS returns an Option that enables or disables TLS for the server.
// Parameters:
//
//	on - a boolean indicating whether to enable (true) or disable (false) TLS.
func WithTLS(on bool) Option {
	return func(s *server) {
		s.tls = on
	}
}
