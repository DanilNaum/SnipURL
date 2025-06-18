package snipendpoint

import (
	"encoding/json"
	"net/http"
)

// getURLs handles HTTP requests to retrieve a list of URLs.
// It returns a JSON response containing the URLs or appropriate HTTP status codes:
// - 200 OK with JSON payload if URLs are found
// - 204 No Content if no URLs exist
// - 500 Internal Server Error if any error occurs during processing
func (s *snipEndpoint) getURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	urls, err := s.service.GetURLs(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	urlsResp, err := getURLsJSONResponseFromServiceModel(s.baseURL, urls)
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
