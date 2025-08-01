package snipendpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
)

// createShortURLBatch handles batch creation of short URLs.
// It accepts a JSON array of URL requests in the request body, where each request contains
// a correlation ID and original URL. The method processes these URLs in batch,
// creates short URLs for each one, and returns a JSON array response.
//
// The response contains corresponding correlation IDs and generated short URLs.
// If successful, it returns HTTP 201 (Created) with the JSON response.
// If there are any errors during processing, it returns appropriate HTTP error codes:
// - 400 Bad Request for invalid JSON input
// - 500 Internal Server Error for server-side processing errors
func (s *snipEndpoint) createShortURLBatch(w http.ResponseWriter, r *http.Request) {
	var req []*createShortURLBatchJSONRequest
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

	urls := make([]*urlsnipper.SetURLsInput, 0, len(req))
	for _, r := range req {
		urls = append(urls, createShortURLBatchJSONRequestToServiceModel(r))
	}

	res, err := s.service.SetURLs(r.Context(), urls)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := make([]*createShortURLBatchJSONResponse, 0, len(res))
	for _, r := range res {
		shortURL, errJoinPath := url.JoinPath(s.baseURL, r.ShortURLID)
		if errJoinPath != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resp = append(resp, &createShortURLBatchJSONResponse{
			CorrelationID: r.CorrelationID,
			ShortURL:      shortURL})

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	marshalResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(marshalResp)

}
