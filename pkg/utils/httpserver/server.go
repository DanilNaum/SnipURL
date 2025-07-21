package httpserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const (
	_defaultShutdownTimeout = 10 * time.Second
	_defaultAddr            = ":8080"

	cacheDir      = "cache-dir"
	hostWhitelist = "*"
)

type server struct {
	tls             bool
	server          *http.Server
	notify          chan error
	idleConnsClosed chan struct{}
	shutdownTimeout time.Duration
}

// NewHTTPServer creates a new HTTP server with the given handler and optional configuration.
// It sets up a server with default address and shutdown timeout, and allows customization
// through optional configuration functions. The server is started immediately upon creation.
//
// Parameters:
//   - ctx: context to close server than ctx.Done
//   - handler: The HTTP handler to serve requests
//   - opts: Optional configuration functions to modify server settings
//
// Returns a configured and started server instance.
func NewHTTPServer(ctx context.Context, handler http.Handler, opts ...Option) *server {
	httpServer := &http.Server{
		Addr:    _defaultAddr,
		Handler: handler,
	}

	s := server{
		server:          httpServer,
		idleConnsClosed: make(chan struct{}, 1),
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(&s)
	}

	s.start()

	go func() {
		select {
		case <-ctx.Done():
		case err := <-s.notify:
			if err != nil {
				log.Printf("server return error: %s, shutdown ...", err.Error())
			}
		}
		err := s.shutdown()
		if err != nil {
			log.Fatalf("server shutdown with error %s", err.Error())
		}

		close(s.idleConnsClosed)

	}()

	return &s

}

func (s *server) start() {
	go func() {
		var err error
		if s.tls {
			manager := &autocert.Manager{
				// директория для хранения сертификатов
				Cache: autocert.DirCache(cacheDir),
				// функция, принимающая Terms of Service издателя сертификатов
				Prompt: autocert.AcceptTOS,
				// перечень доменов, для которых будут поддерживаться сертификаты
				HostPolicy: autocert.HostWhitelist(hostWhitelist),
			}
			s.server.TLSConfig = manager.TLSConfig()
			err = s.server.ListenAndServeTLS("", "")
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil {
			s.notify <- err
		}
		close(s.notify)
	}()
}

// IdleConnsClosed returns the empty channel to signal that server is shutdown
func (s *server) IdleConnsClosed() chan struct{} {
	return s.idleConnsClosed
}

// shutdown gracefully stops the HTTP server, allowing active connections to complete
// within the configured shutdown timeout. Returns an error if the shutdown fails.
func (s *server) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
