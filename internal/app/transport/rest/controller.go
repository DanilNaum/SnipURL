package rest

import (
	"context"
	"net/http"

	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/snipendpoint"
)

type service interface {
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, url string) (string, error)
}

func NewController(mux *http.ServeMux, service service) http.Handler {
	snipEndpoint := snipendpoint.NewSnipEndpoint(service)
	snipEndpoint.Register(mux)

	return mux
}
