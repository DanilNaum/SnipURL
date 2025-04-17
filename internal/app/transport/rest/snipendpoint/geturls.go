package snipendpoint

import (
	"encoding/json"
	"net/http"
)

func (l *snipEndpoint) getURLs(w http.ResponseWriter, r *http.Request) {

	urls, err := l.service.GetURLs(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	urlsResp, err := getURLsJSONResponseFromServiceModel(l.baseURL, urls)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(urlsResp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
