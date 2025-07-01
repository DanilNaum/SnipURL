package pprof

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const endpointPProf = "/debug/pprof/*"

type pprofEndpoint struct {
}

// NewPProfEndpoint creates a new pprofEndpoint instance.
// It returns a pointer to the created pprofEndpoint.
func NewPProfEndpoint() *pprofEndpoint {
	return &pprofEndpoint{}
}

// Register adds the pprof HTTP handler to the provided chi router.
// It maps the default pprof endpoints to the default HTTP ServeMux.
func (p *pprofEndpoint) Register(r *chi.Mux) {
	r.Handle(endpointPProf, http.DefaultServeMux)
}
