package snipendpoint

import (
	"io"
	"net/http"
	"net/url"
)

func (l *snipEndpoint) createShortURL(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	originalURL := string(body)

	id, err := l.service.SetURL(r.Context(), originalURL)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fullShortURL, err := url.JoinPath(l.baseURL, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fullShortURL))
}
