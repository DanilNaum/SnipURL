package snipendpoint

import (
	"context"
	"net/http"
)

const (
	endpointGet  = "/{id}"
	endpointPost = "/"
)

type service interface {
	SetUrl(ctx context.Context, url string) (string, error)
	GetUrl(ctx context.Context, id string) (string, error)
}
type snipEndpoint struct {
	service service
}

func NewSnipEndpoint(service service) *snipEndpoint {
	return &snipEndpoint{
		service: service,
	}
}

func (l *snipEndpoint) Register(handler *http.ServeMux) {

	handler.HandleFunc(endpointPost, l.post)
	handler.HandleFunc(endpointGet, l.get)
}
