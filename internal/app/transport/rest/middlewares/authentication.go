package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/DanilNaum/SnipURL/pkg/cookie"
	"github.com/google/uuid"
)

// Key is a custom type used to store and retrieve a user identifier in a context.
// It provides a type-safe way to set and access a user ID value in request contexts.
type Key struct {
	Key string
}

var key = Key{Key: "userID"}

func (m *middleware) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := m.cookieManager.Get(r)
		if err != nil {
			if !errors.Is(err, cookie.ErrNoCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			userID = uuid.NewString()
			m.cookieManager.Set(w, userID)

		}

		newCtx := context.WithValue(r.Context(), key, userID)

		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
