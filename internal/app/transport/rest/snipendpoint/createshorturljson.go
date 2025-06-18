package snipendpoint

import (
	"bytes"
	"encoding/json"
	"errors"

	"net/http"
	"net/url"

	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
)

// createShortURLJSON handles HTTP POST requests to create a shortened URL.
// It accepts a JSON request containing the original URL and returns a JSON response
// with the shortened URL.
//
// The method:
// 1. Reads and unmarshals the JSON request body
// 2. Calls the URL shortening service to generate a unique ID
// 3. Constructs the full shortened URL using the base URL and generated ID
// 4. Returns the result as a JSON response
//
// Response status codes:
//   - 201 Created: URL successfully shortened
//   - 400 Bad Request: Invalid JSON request
//   - 409 Conflict: URL already exists
//   - 500 Internal Server Error: Server-side error
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
