package snipendpoint

import (
	"context"

	"github.com/go-chi/chi/v5"
)

const (
	endpointGet  = "/{id}"
	endpointPost = "/"
)

type config interface {
	GetPrefix() (string, error)
	BaseURL() string
}

//go:generate moq -out service_moq_test.go . service
type service interface {
	SetURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
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
		baseURL: conf.BaseURL(),
	}, nil
}

func (l *snipEndpoint) Register(r *chi.Mux) {
	r.Route(l.prefix, func(r chi.Router) {
		r.Post(endpointPost, l.post)
		r.Get(endpointGet, l.get)
	})
}
