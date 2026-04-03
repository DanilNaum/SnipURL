package internalendpoints

import (
	"net/http"

	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/utils/responder"
)

func (ie *internalEndpoints) getStats(w http.ResponseWriter, r *http.Request) {
	stats, err := ie.service.GetState(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	responder.RespondWithJSON(w, &State{
		UrlsNum:  stats.UrlsNum,
		UsersNum: stats.UsersNum,
	})
}
