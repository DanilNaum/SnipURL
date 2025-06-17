package snipendpoint

import (
	"net/url"

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

func getURLsJSONResponseFromServiceModel(baseURL string, resp []*urlsnipper.URL) ([]*getURLsJSONResponse, error) {
	urls := make([]*getURLsJSONResponse, 0, len(resp))
	for _, u := range resp {
		fullShortURL, err := url.JoinPath(baseURL, u.ShortURL)
		if err != nil {
			return nil, err
		}
		urls = append(urls, &getURLsJSONResponse{
			ShortURL:    fullShortURL,
			OriginalURL: u.OriginalURL,
		})
	}
	return urls, nil
}
