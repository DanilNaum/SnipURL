package snipendpoint

type postJsonRequest struct {
	URL string `json:"url"`
}

type postJsonResponse struct {
	Result string `json:"result"`
}
