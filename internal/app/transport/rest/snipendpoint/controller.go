package snipendpoint

import (
	"context"

	"github.com/go-chi/chi/v5"
)

const (
	endpointGet  = "/{id}"
	endpointPost = "/"
)

//go:generate moq -out service_moq_test.go . service
type service interface {
	SetURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
}

type snipEndpoint struct {
	service service
}

func NewSnipEndpoint(service service) *snipEndpoint {
	return &snipEndpoint{
		service: service,
	}
}

func (l *snipEndpoint) Register(r *chi.Mux) {
	r.Post(endpointPost, l.post)
	r.Get(endpointGet, l.get)
}
