package rest

import (
	"context"
	"net/http"

	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
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
	SetURLs(ctx context.Context, urls []*urlsnipper.SetURLsInput) (map[string]*urlsnipper.SetURLsOutput, error)
	GetURLs(ctx context.Context) ([]*urlsnipper.Url, error)
}

type psqlStoragePinger interface {
	Ping(context.Context) error
}

type cookieManager interface {
	Set(w http.ResponseWriter, value string)
	Get(r *http.Request) (string, error)
}

func NewController(mux *chi.Mux, conf config, service service, psqlStoragePinger psqlStoragePinger, cookieManager cookieManager, logger logger) (http.Handler, error) {

	middlewares := middlewares.NewMiddleware(logger, cookieManager)

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
