package snipendpoint

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
)

func (l *snipEndpoint) createShortURL(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	originalURL := string(body)

	id, err := l.service.SetURL(r.Context(), originalURL)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusCreated)
	case errors.Is(err, urlsnipper.ErrConflict):
		w.WriteHeader(http.StatusConflict)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fullShortURL, err := url.JoinPath(l.baseURL, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fullShortURL))
}
