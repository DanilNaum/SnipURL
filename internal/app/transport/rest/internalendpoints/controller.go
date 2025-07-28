package internalendpoints

import (
	"context"

	"github.com/DanilNaum/SnipURL/internal/app/service/private"
	"github.com/go-chi/chi/v5"
)

const (
	endpointStats = "/api/internal/stats"
)

type config interface {
	GetPrefix() (string, error)
	GetBaseURL() string
}

type service interface {
	GetState(ctx context.Context) (*private.State, error)
}

type internalEndpoints struct {
	service service
}

func NewInternalEndpoint(service service) (*internalEndpoints, error) {

	return &internalEndpoints{
		service: service,
	}, nil
}

func (s *internalEndpoints) Register(r chi.Router) {

	r.Get(endpointStats, s.getStats)

}
