package httpserver

import (
	"fmt"
	"time"
)

type Option func(s *server)

func WithShutdownTimeout(t time.Duration) Option {
	return func(s *server) {
		s.shutdownTimeout = t
	}
}

func WithPort(port int) Option {
	return func(s *server) {
		s.server.Addr = fmt.Sprintf(":%d", port)
	}
}

func WithAddr(addr string) Option {
	return func(s *server) {
		s.server.Addr = addr
	}

}
