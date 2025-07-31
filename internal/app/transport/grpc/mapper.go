package grpc

import (
	"net/http"

	"github.com/DanilNaum/SnipURL/internal/app/service/private"
	// "github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
	"github.com/DanilNaum/SnipURL/pkg/protobuf"
)

// ShortURL Response Mappers

func shortURLSuccessResponse(shortURL string, statusCode int32, message string) *protobuf.ShortURLResponse {
	return &protobuf.ShortURLResponse{
		Response: &protobuf.ShortURLResponse_Success{
			Success: &protobuf.SuccessShortURL{
				Status: &protobuf.Status{
					Code:    statusCode,
					Message: message,
				},
				ShortUrl: shortURL,
			},
		},
	}
}

func shortURLErrorResponse(statusCode int32, message string) *protobuf.ShortURLResponse {
	return &protobuf.ShortURLResponse{
		Response: &protobuf.ShortURLResponse_Error{
			Error: &protobuf.Error{
				Status: &protobuf.Status{
					Code:    statusCode,
					Message: message,
				},
			},
		},
	}
}

func shortURLConflictResponse(shortURL string) *protobuf.ShortURLResponse {
	return shortURLSuccessResponse(shortURL, http.StatusConflict, "URL already exists")
}

func shortURLCreatedResponse(shortURL string) *protobuf.ShortURLResponse {
	return shortURLSuccessResponse(shortURL, http.StatusCreated, "URL created successfully")
}

func shortURLInternalErrorResponse() *protobuf.ShortURLResponse {
	return shortURLErrorResponse(http.StatusInternalServerError, "Internal server error")
}

func shortURLConstructErrorResponse() *protobuf.ShortURLResponse {
	return shortURLErrorResponse(http.StatusInternalServerError, "Failed to construct URL")
}

// OriginalURL Response Mappers

func originalURLSuccessResponse(originalURL string) *protobuf.OriginalURLResponse {
	return &protobuf.OriginalURLResponse{
		Response: &protobuf.OriginalURLResponse_Success{
			Success: &protobuf.SuccessOriginalURL{
				Status: &protobuf.Status{
					Code:    http.StatusOK,
					Message: "URL found",
				},
				OriginalUrl: originalURL,
			},
		},
	}
}

func originalURLErrorResponse(statusCode int32, message string) *protobuf.OriginalURLResponse {
	return &protobuf.OriginalURLResponse{
		Response: &protobuf.OriginalURLResponse_Error{
			Error: &protobuf.Error{
				Status: &protobuf.Status{
					Code:    statusCode,
					Message: message,
				},
			},
		},
	}
}

func originalURLDeletedResponse() *protobuf.OriginalURLResponse {
	return originalURLErrorResponse(http.StatusGone, "URL has been deleted")
}

func originalURLInternalErrorResponse() *protobuf.OriginalURLResponse {
	return originalURLErrorResponse(http.StatusInternalServerError, "Internal server error")
}

// JsonShortURL Response Mappers

func jsonShortURLSuccessResponse(shortURL string, statusCode int32, message string) *protobuf.JsonShortURLResponse {
	return &protobuf.JsonShortURLResponse{
		Response: &protobuf.JsonShortURLResponse_Success{
			Success: &protobuf.SuccessJsonShortURL{
				Status: &protobuf.Status{
					Code:    statusCode,
					Message: message,
				},
				ShortUrl: shortURL,
			},
		},
	}
}

func jsonShortURLErrorResponse(statusCode int32, message string) *protobuf.JsonShortURLResponse {
	return &protobuf.JsonShortURLResponse{
		Response: &protobuf.JsonShortURLResponse_Error{
			Error: &protobuf.Error{
				Status: &protobuf.Status{
					Code:    statusCode,
					Message: message,
				},
			},
		},
	}
}

func jsonShortURLConflictResponse(shortURL string) *protobuf.JsonShortURLResponse {
	return jsonShortURLSuccessResponse(shortURL, http.StatusConflict, "URL already exists")
}

func jsonShortURLCreatedResponse(shortURL string) *protobuf.JsonShortURLResponse {
	return jsonShortURLSuccessResponse(shortURL, http.StatusCreated, "URL created successfully")
}

func jsonShortURLInternalErrorResponse() *protobuf.JsonShortURLResponse {
	return jsonShortURLErrorResponse(http.StatusInternalServerError, "Internal server error")
}

func jsonShortURLConstructErrorResponse() *protobuf.JsonShortURLResponse {
	return jsonShortURLErrorResponse(http.StatusInternalServerError, "Failed to construct URL")
}

// BatchCreate Response Mappers

func batchCreateSuccessResponse(items []*protobuf.BatchCreateResponseItem) *protobuf.BatchCreateResponse {
	return &protobuf.BatchCreateResponse{
		Response: &protobuf.BatchCreateResponse_Success{
			Success: &protobuf.SuccessBatchCreate{
				Status: &protobuf.Status{
					Code:    http.StatusCreated,
					Message: "URLs created successfully",
				},
				Items: items,
			},
		},
	}
}

func batchCreateErrorResponse(statusCode int32, message string) *protobuf.BatchCreateResponse {
	return &protobuf.BatchCreateResponse{
		Response: &protobuf.BatchCreateResponse_Error{
			Error: &protobuf.Error{
				Status: &protobuf.Status{
					Code:    statusCode,
					Message: message,
				},
			},
		},
	}
}

func batchCreateInternalErrorResponse() *protobuf.BatchCreateResponse {
	return batchCreateErrorResponse(http.StatusInternalServerError, "Internal server error")
}

func batchCreateConstructErrorResponse() *protobuf.BatchCreateResponse {
	return batchCreateErrorResponse(http.StatusInternalServerError, "Failed to construct URL")
}

func batchCreateResponseItem(correlationID, shortURL string) *protobuf.BatchCreateResponseItem {
	return &protobuf.BatchCreateResponseItem{
		CorrelationId: correlationID,
		ShortUrl:      shortURL,
	}
}

// UserURLs Response Mappers

func userURLsSuccessResponse(items []*protobuf.UserURLItem, statusCode int32, message string) *protobuf.UserURLsResponse {
	return &protobuf.UserURLsResponse{
		Response: &protobuf.UserURLsResponse_Success{
			Success: &protobuf.SuccessUserURLs{
				Status: &protobuf.Status{
					Code:    statusCode,
					Message: message,
				},
				Items: items,
			},
		},
	}
}

func userURLsErrorResponse(statusCode int32, message string) *protobuf.UserURLsResponse {
	return &protobuf.UserURLsResponse{
		Response: &protobuf.UserURLsResponse_Error{
			Error: &protobuf.Error{
				Status: &protobuf.Status{
					Code:    statusCode,
					Message: message,
				},
			},
		},
	}
}

func userURLsFoundResponse(items []*protobuf.UserURLItem) *protobuf.UserURLsResponse {
	return userURLsSuccessResponse(items, http.StatusOK, "URLs retrieved successfully")
}

func userURLsNoContentResponse() *protobuf.UserURLsResponse {
	return userURLsSuccessResponse([]*protobuf.UserURLItem{}, http.StatusNoContent, "No URLs found")
}

func userURLsInternalErrorResponse() *protobuf.UserURLsResponse {
	return userURLsErrorResponse(http.StatusInternalServerError, "Internal server error")
}

func userURLsConstructErrorResponse() *protobuf.UserURLsResponse {
	return userURLsErrorResponse(http.StatusInternalServerError, "Failed to construct URL")
}

func userURLItem(shortURL, originalURL string) *protobuf.UserURLItem {
	return &protobuf.UserURLItem{
		ShortUrl:    shortURL,
		OriginalUrl: originalURL,
	}
}

// Delete Response Mappers

func deleteAcceptedResponse() *protobuf.DeleteResponse {
	return &protobuf.DeleteResponse{
		Status: &protobuf.Status{
			Code:    http.StatusAccepted,
			Message: "URLs deletion accepted",
		},
	}
}

// Ping Response Mappers

func pingResponse(statusCode int32, message string) *protobuf.PingResponse {
	return &protobuf.PingResponse{
		Status: &protobuf.Status{
			Code:    statusCode,
			Message: message,
		},
	}
}

func pingSuccessResponse() *protobuf.PingResponse {
	return pingResponse(http.StatusOK, "Database available")
}

func pingErrorResponse() *protobuf.PingResponse {
	return pingResponse(http.StatusInternalServerError, "Database unavailable")
}

// Stats Response Mappers

func statsSuccessResponse(stats *private.State) *protobuf.StatsResponse {
	return &protobuf.StatsResponse{
		Response: &protobuf.StatsResponse_Success{
			Success: &protobuf.SuccessStats{
				Status: &protobuf.Status{
					Code:    http.StatusOK,
					Message: "Stats retrieved successfully",
				},
				Data: &protobuf.StatsData{
					Urls:  int32(stats.UrlsNum),
					Users: int32(stats.UsersNum),
				},
			},
		},
	}
}

func statsErrorResponse(statusCode int32, message string) *protobuf.StatsResponse {
	return &protobuf.StatsResponse{
		Response: &protobuf.StatsResponse_Error{
			Error: &protobuf.Error{
				Status: &protobuf.Status{
					Code:    statusCode,
					Message: message,
				},
			},
		},
	}
}

func statsInternalErrorResponse() *protobuf.StatsResponse {
	return statsErrorResponse(http.StatusInternalServerError, "Internal server error")
}
