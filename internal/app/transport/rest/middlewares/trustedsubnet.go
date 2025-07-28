package middlewares

import (
	"net/http"
)

const xRealIPHeader = "X-Real-IP"

// IsTrustedSubNet is a middleware that validates incoming requests against a trusted subnet.
// It checks the X-Real-IP header and returns HTTP 403 Forbidden if the IP is missing
// or doesn't match the configured trusted subnet. If validation passes, the request
// is forwarded to the next handler in the chain.
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
