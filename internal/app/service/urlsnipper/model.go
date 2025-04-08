package urlsnipper

type SetURLsInput struct {
	CorrelationID string
	OriginalURL   string
}

type SetURLsOutput struct {
	CorrelationID string
	ShortURLID    string
}
