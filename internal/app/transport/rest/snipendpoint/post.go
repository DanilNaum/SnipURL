package snipendpoint

import (
	"io"
	"net/http"
)

func (l *snipEndpoint) post(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
	}

	url := string(body)

	id, err := l.service.SetUrl(r.Context(), url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(id))

	w.WriteHeader(http.StatusCreated)

}
