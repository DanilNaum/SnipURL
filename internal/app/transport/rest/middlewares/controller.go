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
	trustedSubnet string
}

// NewMiddleware creates a new middleware instance with the provided logger and cookie manager.
func NewMiddleware(logger logger, cookieManager cookieManager, trustedSubnet string) *middleware {
	return &middleware{
		logger:        logger,
		cookieManager: cookieManager,
		trustedSubnet: trustedSubnet,
	}
}

// Register configures and applies middleware to the given chi router.
// It adds authentication, logging, gzip compression, and decompression middleware.
func (m *middleware) Register(mux *chi.Mux) *chi.Mux {
	mux.Use(m.authentication)
	mux.Use(m.logging)

	mux.Use(m.gzipPack)

	mux.Use(m.gzipUnpack)

	return mux
}

// func (m *middleware) RegisterForInternalReq(mux *chi.Mux) *chi.Mux {

// 	mux.Use(m.isTrustedSubNet)
// 	return mux
// }
