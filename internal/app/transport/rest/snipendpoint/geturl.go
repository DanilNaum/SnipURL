package snipendpoint

import (
	"errors"
	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
	"net/http"
)

func (s *snipEndpoint) getURL(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	url, err := s.service.GetURL(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, urlsnipper.ErrDeleted):
			http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
			return
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
