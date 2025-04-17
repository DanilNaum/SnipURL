package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type logger interface {
	Infoln(args ...any)
}

type cookieManager interface {
	Set(w http.ResponseWriter, value string)
	Get(r *http.Request) (string, error)
}

type middleware struct {
	logger        logger
	cookieManager cookieManager
}

func NewMiddleware(logger logger, cookieManager cookieManager) *middleware {
	return &middleware{
		logger:        logger,
		cookieManager: cookieManager,
	}
}

func (m *middleware) Register(mux *chi.Mux) *chi.Mux {

	mux.Use(m.authentication)
	mux.Use(m.logging)

	mux.Use(m.gzipPack)

	mux.Use(m.gzipUnpack)

	return mux
}
