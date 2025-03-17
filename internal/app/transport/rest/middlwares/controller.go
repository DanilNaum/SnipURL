package middlewares

import "github.com/go-chi/chi/v5"

type logger interface {
	Infoln(args ...any)
}

type middleware struct {
	logger logger
}

func NewMiddleware(logger logger) *middleware {
	return &middleware{
		logger: logger,
	}
}

func (m *middleware) Register(mux *chi.Mux) *chi.Mux {
	mux.Use(m.logging)
	return mux
}
