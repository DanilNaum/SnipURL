package url

// URLRecord represents a shortened URL entry in the database.
// It contains information about the original URL, its shortened version, and associated metadata.
type URLRecord struct {
	ID          int
	ShortURL    string
	OriginalURL string
	UserID      string
	Deleted     bool
}
