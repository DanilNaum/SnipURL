package snipendpoint

import (
	"fmt"
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	url := string(body)

	id, err := l.service.SetURL(r.Context(), url)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	scheme := "http"  // По умолчанию используем HTTP
	if r.TLS != nil { // Если запрос пришел по HTTPS, используем HTTPS
		scheme = "https"
	}
	host := r.Host // Используем хост из запроса

	// Формируем полный URL
	fullURL := fmt.Sprintf("%s://%s/%s", scheme, host, id)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fullURL))
}
