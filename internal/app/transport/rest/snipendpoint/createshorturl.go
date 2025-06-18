package snipendpoint

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
)

// createShortURL handles HTTP requests to create a shortened URL.
// It reads the original URL from the request body, generates a short ID using the URL snipper service,
// and returns the full shortened URL path.
//
// The response status codes are:
//   - 201 (Created) if the URL was successfully shortened
//   - 409 (Conflict) if the URL already exists
//   - 500 (Internal Server Error) if any internal error occurs
//
// The response body contains the full shortened URL if successful.
func (s *snipEndpoint) createShortURL(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	originalURL := string(body)

	id, err := s.service.SetURL(r.Context(), originalURL)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusCreated)
	case errors.Is(err, urlsnipper.ErrConflict):
		w.WriteHeader(http.StatusConflict)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fullShortURL, err := url.JoinPath(s.baseURL, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fullShortURL))
}
