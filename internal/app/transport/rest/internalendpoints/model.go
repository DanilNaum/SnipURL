package internalendpoints

// State represents the current state with counts of URLs and users.
// It includes JSON tags for serialization.
type State struct {
	UrlsNum  int `json:"urls"`
	UsersNum int `json:"user"`
}
