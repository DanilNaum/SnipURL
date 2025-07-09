package httpserver

import (
	"context"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"time"
)

const (
	_defaultShutdownTimeout = 10 * time.Second
	_defaultAddr            = ":8080"
)

type server struct {
	tls             bool
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// NewHTTPServer creates a new HTTP server with the given handler and optional configuration.
// It sets up a server with default address and shutdown timeout, and allows customization
// through optional configuration functions. The server is started immediately upon creation.
//
// Parameters:
//   - handler: The HTTP handler to serve requests
//   - opts: Optional configuration functions to modify server settings
//
// Returns a configured and started server instance.
func NewHTTPServer(handler http.Handler, opts ...Option) *server {
	httpServer := &http.Server{
		Addr:    _defaultAddr,
		Handler: handler,
	}

	s := server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(&s)
	}

	s.start()

	return &s

}

func (s *server) start() {
	go func() {
		var err error
		if s.tls {
			manager := &autocert.Manager{
				// директория для хранения сертификатов
				Cache: autocert.DirCache("cache-dir"),
				// функция, принимающая Terms of Service издателя сертификатов
				Prompt: autocert.AcceptTOS,
				// перечень доменов, для которых будут поддерживаться сертификаты
				HostPolicy: autocert.HostWhitelist("*"),
			}
			s.server.TLSConfig = manager.TLSConfig()
			err = s.server.ListenAndServeTLS("", "")
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil {
			if err != http.ErrServerClosed {
				s.notify <- err
			}
		}
		close(s.notify)
	}()
}

// Notify returns the error notification channel for the server.
// This channel receives any non-ErrServerClosed errors that occur during server startup or operation.
func (s *server) Notify() chan error {
	return s.notify
}

// Shutdown gracefully stops the HTTP server, allowing active connections to complete
// within the configured shutdown timeout. Returns an error if the shutdown fails.
func (s *server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
