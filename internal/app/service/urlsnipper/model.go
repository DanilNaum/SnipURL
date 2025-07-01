package urlsnipper

// SetURLsInput represents the input parameters for setting URLs in the URL snipper service.
type SetURLsInput struct {
	CorrelationID string
	OriginalURL   string
}

// SetURLsOutput represents the output returned after setting a URL in the URL snipper service.
type SetURLsOutput struct {
	CorrelationID string
	ShortURLID    string
}

// URL represents a mapping between a short URL and its original long URL.
type URL struct {
	ShortURL    string
	OriginalURL string
}
