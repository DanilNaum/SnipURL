package snipendpoint

import (
	"net/http"
)

func (s *snipEndpoint) getURL(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	url, err := s.service.GetURL(r.Context(), id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
