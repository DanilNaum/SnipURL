package rest

import (
	"context"
	"net/http"

	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/snipendpoint"
	"github.com/go-chi/chi/v5"
)

type config interface {
	GetPrefix() (string, error)
	BaseURL() string
}

type service interface {
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, url string) (string, error)
}

func NewController(mux *chi.Mux, conf config, service service) (http.Handler, error) {

	snipEndpoint, err := snipendpoint.NewSnipEndpoint(service, conf)
	if err != nil {
		return nil, err
	}
	snipEndpoint.Register(mux)

	return mux, nil
}
