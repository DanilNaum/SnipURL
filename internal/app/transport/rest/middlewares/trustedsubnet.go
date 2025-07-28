package middlewares

import (
	"net/http"
)

const xRealIPHeader = "X-Real-IP"

func (m *middleware) IsTrustedSubNet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRealIP := r.Header.Get(xRealIPHeader)
		if xRealIP == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if xRealIP != m.trustedSubnet {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
