package httpserver

import (
	"fmt"
	"strings"
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
		parts := strings.Split(addr, ":")
		if len(parts) < 2 {
			s.server.Addr = ":80" // Если порт не указан, возвращаем пустую строку
			return
		}
		s.server.Addr = ":" + parts[len(parts)-1]
	}

}
