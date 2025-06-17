package snipendpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (s *snipEndpoint) deleteURLs(w http.ResponseWriter, r *http.Request) {
	var req []string
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

	s.service.DeleteURLs(r.Context(), req)

	w.WriteHeader(http.StatusAccepted)
}
