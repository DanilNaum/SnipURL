package snipendpoint

import (
	"errors"
	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
	"net/http"
)

// getURL обрабатывает HTTP-запрос для получения URL по его идентификатору.
//
// Этот метод извлекает идентификатор из пути запроса и использует сервис
// для получения соответствующего URL. Если URL был удален, метод возвращает
// статус 410 Gone. В случае других ошибок возвращается статус 500 Internal Server Error.
// Если URL успешно найден, происходит перенаправление на этот URL с кодом 302 Found.
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
