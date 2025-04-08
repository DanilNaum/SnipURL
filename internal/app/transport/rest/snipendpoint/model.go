package snipendpoint

type createShortURLJSONRequest struct {
	URL string `json:"url"`
}

type createShortURLJSONResponse struct {
	Result string `json:"result"`
}

type createShortURLBatchJSONRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type createShortURLBatchJSONResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
