package rest

import (
	"context"
	"net/http"

	middlewares "github.com/DanilNaum/SnipURL/internal/app/transport/rest/middlwares"
	psqlping "github.com/DanilNaum/SnipURL/internal/app/transport/rest/psqlPing"
	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/snipendpoint"
	"github.com/go-chi/chi/v5"
)

type logger interface {
	Infoln(args ...any)
}

type config interface {
	GetPrefix() (string, error)
	GetBaseURL() string
}

type service interface {
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, url string) (string, error)
}

type psqlStoragePinger interface {
	Ping(context.Context) error
}

func NewController(mux *chi.Mux, conf config, service service, psqlStoragePinger psqlStoragePinger, logger logger) (http.Handler, error) {

	middlewares := middlewares.NewMiddleware(logger)

	muxWithMiddlewares := middlewares.Register(mux)

	snipEndpoint, err := snipendpoint.NewSnipEndpoint(service, conf)
	if err != nil {
		return nil, err
	}

	psqlPingEndpoint := psqlping.NewPsqlPingEndpoint(psqlStoragePinger)

	psqlPingEndpoint.Register(muxWithMiddlewares)
	snipEndpoint.Register(muxWithMiddlewares)

	return muxWithMiddlewares, nil
}
