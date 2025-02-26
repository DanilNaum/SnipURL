package rest

import (
	"context"
	"net/http"

	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/snipendpoint"
	"github.com/go-chi/chi/v5"
)

type service interface {
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, url string) (string, error)
}

func NewController(mux *chi.Mux, service service) http.Handler {
	snipEndpoint := snipendpoint.NewSnipEndpoint(service)
	snipEndpoint.Register(mux)

	return mux
}
