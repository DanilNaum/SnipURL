package snipendpoint

import (
	"bytes"
	"encoding/json"
	"errors"

	"net/http"
	"net/url"

	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
)

func (s *snipEndpoint) createShortURLJSON(w http.ResponseWriter, r *http.Request) {
	var req createShortURLJSONRequest
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	originalURL := req.URL

	id, err := s.service.SetURL(r.Context(), originalURL)
	switch {
	case err == nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

	case errors.Is(err, urlsnipper.ErrConflict):
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fullURL, err := url.JoinPath(s.baseURL, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&createShortURLJSONResponse{
		Result: fullURL,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
