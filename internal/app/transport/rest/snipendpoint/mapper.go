package snipendpoint

import (
	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
)

func createShortURLBatchJSONRequestToServiceModel(req *createShortURLBatchJSONRequest) *urlsnipper.SetURLsInput {
	return &urlsnipper.SetURLsInput{
		CorrelationID: req.CorrelationID,
		OriginalURL:   req.OriginalURL,
	}
}

func createShortURLBatchJSONResponseFromServiceModel(resp *urlsnipper.SetURLsOutput) *createShortURLBatchJSONResponse {
	return &createShortURLBatchJSONResponse{
		CorrelationID: resp.CorrelationID,
		ShortURL:      resp.ShortURLID,
	}
}
