package pprof

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const endpointPProf = "/debug/pprof/*"

type pprofEndpoint struct {
}

func NewPProfEndpoint() *pprofEndpoint {
	return &pprofEndpoint{}
}

func (p *pprofEndpoint) Register(r *chi.Mux) {
	r.Handle(endpointPProf, http.DefaultServeMux)
}
