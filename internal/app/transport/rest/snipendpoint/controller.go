package snipendpoint

import (
	"context"

	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
	"github.com/go-chi/chi/v5"
)

const (
	endpointGetURL              = "/{id}"
	endpointCreateShortURL      = "/"
	endpointCreateShortURLJSON  = "/api/shorten"
	endpointCreateShortURLBatch = "/api/shorten/batch"
	endpointGetUserURLs         = "/api/user/urls"
	endpointDeleteURLs          = "/api/user/urls"
)

type config interface {
	GetPrefix() (string, error)
	GetBaseURL() string
}

//go:generate moq -out service_moq_test.go . service
type service interface {
	SetURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
	SetURLs(ctx context.Context, urls []*urlsnipper.SetURLsInput) (map[string]*urlsnipper.SetURLsOutput, error)
	GetURLs(ctx context.Context) ([]*urlsnipper.URL, error)
	DeleteURLs(ctx context.Context, ids []string)
}

type snipEndpoint struct {
	service service
	prefix  string
	baseURL string
}

func NewSnipEndpoint(service service, conf config) (*snipEndpoint, error) {
	prefix, err := conf.GetPrefix()
	if err != nil {
		return nil, err
	}
	return &snipEndpoint{
		service: service,
		prefix:  prefix,
		baseURL: conf.GetBaseURL(),
	}, nil
}

func (l *snipEndpoint) Register(r *chi.Mux) {
	r.Route(l.prefix, func(r chi.Router) {
		r.Post(endpointCreateShortURL, l.createShortURL)
		r.Get(endpointGetURL, l.getURL)
		r.Post(endpointCreateShortURLJSON, l.createShortURLJSON)
		r.Post(endpointCreateShortURLBatch, l.createShortURLBatch)
		r.Get(endpointGetUserURLs, l.getURLs)
		r.Delete(endpointDeleteURLs, l.deleteURLs)

	})
}
