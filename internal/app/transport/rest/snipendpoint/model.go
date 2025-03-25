package snipendpoint

type postJSONRequest struct {
	URL string `json:"url"`
}

type postJSONResponse struct {
	Result string `json:"result"`
}
