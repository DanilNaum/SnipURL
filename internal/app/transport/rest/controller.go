package rest

import (
	"context"
	"net/http"

	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/internalendpoints"
	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/pprof"

	"github.com/DanilNaum/SnipURL/internal/app/service/private"
	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
	middlewares "github.com/DanilNaum/SnipURL/internal/app/transport/rest/middlewares"
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
	GetTrustedSubNet() string
}

type service interface {
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, url string) (string, error)
	SetURLs(ctx context.Context, urls []*urlsnipper.SetURLsInput) (map[string]*urlsnipper.SetURLsOutput, error)
	GetURLs(ctx context.Context) ([]*urlsnipper.URL, error)
	DeleteURLs(ctx context.Context, ids []string)
}

type internalService interface {
	GetState(ctx context.Context) (*private.State, error)
}

type psqlStoragePinger interface {
	Ping(context.Context) error
}

type cookieManager interface {
	Set(w http.ResponseWriter, value string)
	Get(r *http.Request) (string, error)
}

// NewController creates and configures a new HTTP handler with middleware, endpoints for URL shortening,
// PostgreSQL ping, and pprof profiling. It sets up routes using the provided chi router and various
// dependencies like configuration, service, storage pinger, cookie manager, and logger.
//
// Parameters:
//   - mux: The base chi router to be configured
//   - conf: Configuration interface for retrieving application settings
//   - service: Service interface for URL shortening operations
//   - psqlStoragePinger: Interface for checking PostgreSQL storage connectivity
//   - cookieManager: Interface for managing HTTP cookies
//   - logger: Logger interface for logging information
//
// Returns an configured HTTP handler and an error if initialization fails.
func NewController(mux *chi.Mux, conf config, service service, internalService internalService, psqlStoragePinger psqlStoragePinger, cookieManager cookieManager, logger logger) (http.Handler, error) {

	middlewares := middlewares.NewMiddleware(logger, cookieManager, conf.GetTrustedSubNet())

	muxWithMiddlewares := middlewares.Register(mux)
	// muxWithInternalMiddlewares := middlewares.RegisterForInternalReq(mux)

	snipEndpoint, err := snipendpoint.NewSnipEndpoint(service, conf)
	if err != nil {
		return nil, err
	}

	psqlPingEndpoint := psqlping.NewPsqlPingEndpoint(psqlStoragePinger)

	psqlPingEndpoint.Register(muxWithMiddlewares)

	snipEndpoint.Register(muxWithMiddlewares)

	pprofEndpoint := pprof.NewPProfEndpoint()
	pprofEndpoint.Register(muxWithMiddlewares)

	internalEndpoints, err := internalendpoints.NewInternalEndpoint(internalService)
	if err != nil {
		return nil, err
	}

	internalEndpoints.Register(mux.With(middlewares.IsTrustedSubNet))

	return mux, nil
}
