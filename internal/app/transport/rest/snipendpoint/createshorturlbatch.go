package snipendpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
)

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
		shortURL, err := url.JoinPath(s.baseURL, r.ShortURLID)
		if err != nil {
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
