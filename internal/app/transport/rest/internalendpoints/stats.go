package internalendpoints

import (
	"encoding/json"
	"net/http"
)

func (ie *internalEndpoints) getStats(w http.ResponseWriter, r *http.Request) {
	stats, err := ie.service.GetState(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	resp, err := json.Marshal(&State{
		UrlsNum:  stats.UrlsNum,
		UsersNum: stats.UsersNum,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
