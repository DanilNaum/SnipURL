package psqlping

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	endpointPing = "/ping"
)

type psqlStoragePinger interface {
	Ping(context.Context) error
}

type psqlPingEndpoint struct {
	psqlStoragePinger psqlStoragePinger
}

func NewPsqlPingEndpoint(psqlStoragePinger psqlStoragePinger) *psqlPingEndpoint {
	return &psqlPingEndpoint{
		psqlStoragePinger: psqlStoragePinger,
	}
}

func (l *psqlPingEndpoint) Register(r *chi.Mux) {
	r.Get(endpointPing, l.ping)
}

func (l *psqlPingEndpoint) ping(w http.ResponseWriter, r *http.Request) {
	err := l.psqlStoragePinger.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
