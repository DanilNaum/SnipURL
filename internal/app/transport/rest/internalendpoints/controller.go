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

// NewInternalEndpoint creates a new instance of internalEndpoints with the provided service.
// Returns the initialized internalEndpoints struct and nil error.
func NewInternalEndpoint(service service) (*internalEndpoints, error) {

	return &internalEndpoints{
		service: service,
	}, nil
}

// Register registers internal endpoint routes with the provided chi router.
// It sets up the stats endpoint for retrieving internal statistics.
func (ie *internalEndpoints) Register(r chi.Router) {

	r.Get(endpointStats, ie.getStats)

}
