package grpc

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/DanilNaum/SnipURL/internal/app/service/private"
	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
	"github.com/DanilNaum/SnipURL/pkg/protobuf"
	"google.golang.org/protobuf/types/known/emptypb"
)

type config interface {
	GetPrefix() (string, error)
	GetBaseURL() string
}

type service interface {
	SetURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
	SetURLs(ctx context.Context, urls []*urlsnipper.SetURLsInput) (map[string]*urlsnipper.SetURLsOutput, error)
	GetURLs(ctx context.Context) ([]*urlsnipper.URL, error)
	DeleteURLs(ctx context.Context, ids []string)
}

type internalService interface {
	GetState(ctx context.Context) (*private.State, error)
}

type psqlStoragePinger interface {
	Ping(context.Context) error
}

// Server представляет gRPC сервер для SnipURL
type Server struct {
	protobuf.UnimplementedSnipURLServiceServer
	service           service
	internalService   internalService
	psqlStoragePinger psqlStoragePinger
	baseURL           string
}

// NewServer создает новый экземпляр gRPC сервера
func NewServer(
	service service,
	internalService internalService,
	psqlStoragePinger psqlStoragePinger,
	conf config,
) (*Server, error) {
	return &Server{
		service:           service,
		internalService:   internalService,
		psqlStoragePinger: psqlStoragePinger,
		baseURL:           conf.GetBaseURL(),
	}, nil
}

// CreateShortURL создает короткую ссылку из переданного URL
func (s *Server) CreateShortURL(ctx context.Context, req *protobuf.ShortURLRequest) (*protobuf.ShortURLResponse, error) {
	id, err := s.service.SetURL(ctx, req.Url)
	if err != nil {
		if errors.Is(err, urlsnipper.ErrConflict) {
			// Для конфликта все равно возвращаем успешный ответ с кодом 409
			fullShortURL, urlErr := url.JoinPath(s.baseURL, id)
			if urlErr != nil {
				return &protobuf.ShortURLResponse{
					Response: &protobuf.ShortURLResponse_Error{
						Error: &protobuf.Error{
							Status: &protobuf.Status{
								Code:    http.StatusInternalServerError,
								Message: "Failed to construct URL",
							},
						},
					},
				}, nil
			}
			return &protobuf.ShortURLResponse{
				Response: &protobuf.ShortURLResponse_Success{
					Success: &protobuf.SuccessShortURL{
						Status: &protobuf.Status{
							Code:    http.StatusConflict,
							Message: "URL already exists",
						},
						ShortUrl: fullShortURL,
					},
				},
			}, nil
		}
		return &protobuf.ShortURLResponse{
			Response: &protobuf.ShortURLResponse_Error{
				Error: &protobuf.Error{
					Status: &protobuf.Status{
						Code:    http.StatusInternalServerError,
						Message: "Internal server error",
					},
				},
			},
		}, nil
	}

	fullShortURL, err := url.JoinPath(s.baseURL, id)
	if err != nil {
		return &protobuf.ShortURLResponse{
			Response: &protobuf.ShortURLResponse_Error{
				Error: &protobuf.Error{
					Status: &protobuf.Status{
						Code:    http.StatusInternalServerError,
						Message: "Failed to construct URL",
					},
				},
			},
		}, nil
	}

	return &protobuf.ShortURLResponse{
		Response: &protobuf.ShortURLResponse_Success{
			Success: &protobuf.SuccessShortURL{
				Status: &protobuf.Status{
					Code:    http.StatusCreated,
					Message: "URL created successfully",
				},
				ShortUrl: fullShortURL,
			},
		},
	}, nil
}

// GetOriginalURL получает оригинальный URL по короткому ID
func (s *Server) GetOriginalURL(ctx context.Context, req *protobuf.ShortURLID) (*protobuf.OriginalURLResponse, error) {
	originalURL, err := s.service.GetURL(ctx, req.Id)
	if err != nil {
		if errors.Is(err, urlsnipper.ErrDeleted) {
			return &protobuf.OriginalURLResponse{
				Response: &protobuf.OriginalURLResponse_Error{
					Error: &protobuf.Error{
						Status: &protobuf.Status{
							Code:    http.StatusGone,
							Message: "URL has been deleted",
						},
					},
				},
			}, nil
		}
		return &protobuf.OriginalURLResponse{
			Response: &protobuf.OriginalURLResponse_Error{
				Error: &protobuf.Error{
					Status: &protobuf.Status{
						Code:    http.StatusInternalServerError,
						Message: "Internal server error",
					},
				},
			},
		}, nil
	}

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
	}, nil
}

// CreateShortURLJson создает короткую ссылку из JSON запроса
func (s *Server) CreateShortURLJson(ctx context.Context, req *protobuf.JsonShortURLRequest) (*protobuf.JsonShortURLResponse, error) {
	id, err := s.service.SetURL(ctx, req.Url)
	if err != nil {
		if errors.Is(err, urlsnipper.ErrConflict) {
			fullShortURL, urlErr := url.JoinPath(s.baseURL, id)
			if urlErr != nil {
				return &protobuf.JsonShortURLResponse{
					Response: &protobuf.JsonShortURLResponse_Error{
						Error: &protobuf.Error{
							Status: &protobuf.Status{
								Code:    http.StatusInternalServerError,
								Message: "Failed to construct URL",
							},
						},
					},
				}, nil
			}
			return &protobuf.JsonShortURLResponse{
				Response: &protobuf.JsonShortURLResponse_Success{
					Success: &protobuf.SuccessJsonShortURL{
						Status: &protobuf.Status{
							Code:    http.StatusConflict,
							Message: "URL already exists",
						},
						ShortUrl: fullShortURL,
					},
				},
			}, nil
		}
		return &protobuf.JsonShortURLResponse{
			Response: &protobuf.JsonShortURLResponse_Error{
				Error: &protobuf.Error{
					Status: &protobuf.Status{
						Code:    http.StatusInternalServerError,
						Message: "Internal server error",
					},
				},
			},
		}, nil
	}

	fullShortURL, err := url.JoinPath(s.baseURL, id)
	if err != nil {
		return &protobuf.JsonShortURLResponse{
			Response: &protobuf.JsonShortURLResponse_Error{
				Error: &protobuf.Error{
					Status: &protobuf.Status{
						Code:    http.StatusInternalServerError,
						Message: "Failed to construct URL",
					},
				},
			},
		}, nil
	}

	return &protobuf.JsonShortURLResponse{
		Response: &protobuf.JsonShortURLResponse_Success{
			Success: &protobuf.SuccessJsonShortURL{
				Status: &protobuf.Status{
					Code:    http.StatusCreated,
					Message: "URL created successfully",
				},
				ShortUrl: fullShortURL,
			},
		},
	}, nil
}

// BatchCreateShortURLs создает несколько коротких ссылок за один запрос
func (s *Server) BatchCreateShortURLs(ctx context.Context, req *protobuf.BatchCreateRequest) (*protobuf.BatchCreateResponse, error) {
	urls := make([]*urlsnipper.SetURLsInput, 0, len(req.Items))
	for _, item := range req.Items {
		urls = append(urls, &urlsnipper.SetURLsInput{
			CorrelationID: item.CorrelationId,
			OriginalURL:   item.OriginalUrl,
		})
	}

	result, err := s.service.SetURLs(ctx, urls)
	if err != nil {
		return &protobuf.BatchCreateResponse{
			Response: &protobuf.BatchCreateResponse_Error{
				Error: &protobuf.Error{
					Status: &protobuf.Status{
						Code:    http.StatusInternalServerError,
						Message: "Internal server error",
					},
				},
			},
		}, nil
	}

	items := make([]*protobuf.BatchCreateResponseItem, 0, len(result))
	for _, res := range result {
		shortURL, err := url.JoinPath(s.baseURL, res.ShortURLID)
		if err != nil {
			return &protobuf.BatchCreateResponse{
				Response: &protobuf.BatchCreateResponse_Error{
					Error: &protobuf.Error{
						Status: &protobuf.Status{
							Code:    http.StatusInternalServerError,
							Message: "Failed to construct URL",
						},
					},
				},
			}, nil
		}
		items = append(items, &protobuf.BatchCreateResponseItem{
			CorrelationId: res.CorrelationID,
			ShortUrl:      shortURL,
		})
	}

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
	}, nil
}

// GetUserURLs получает все URL пользователя
func (s *Server) GetUserURLs(ctx context.Context, req *emptypb.Empty) (*protobuf.UserURLsResponse, error) {
	urls, err := s.service.GetURLs(ctx)
	if err != nil {
		return &protobuf.UserURLsResponse{
			Response: &protobuf.UserURLsResponse_Error{
				Error: &protobuf.Error{
					Status: &protobuf.Status{
						Code:    http.StatusInternalServerError,
						Message: "Internal server error",
					},
				},
			},
		}, nil
	}

	if len(urls) == 0 {
		return &protobuf.UserURLsResponse{
			Response: &protobuf.UserURLsResponse_Success{
				Success: &protobuf.SuccessUserURLs{
					Status: &protobuf.Status{
						Code:    http.StatusNoContent,
						Message: "No URLs found",
					},
					Items: []*protobuf.UserURLItem{},
				},
			},
		}, nil
	}

	items := make([]*protobuf.UserURLItem, 0, len(urls))
	for _, urlItem := range urls {
		shortURL, err := url.JoinPath(s.baseURL, urlItem.ShortURL)
		if err != nil {
			return &protobuf.UserURLsResponse{
				Response: &protobuf.UserURLsResponse_Error{
					Error: &protobuf.Error{
						Status: &protobuf.Status{
							Code:    http.StatusInternalServerError,
							Message: "Failed to construct URL",
						},
					},
				},
			}, nil
		}
		items = append(items, &protobuf.UserURLItem{
			ShortUrl:    shortURL,
			OriginalUrl: urlItem.OriginalURL,
		})
	}

	return &protobuf.UserURLsResponse{
		Response: &protobuf.UserURLsResponse_Success{
			Success: &protobuf.SuccessUserURLs{
				Status: &protobuf.Status{
					Code:    http.StatusOK,
					Message: "URLs retrieved successfully",
				},
				Items: items,
			},
		},
	}, nil
}

// DeleteUserURLs удаляет URL пользователя
func (s *Server) DeleteUserURLs(ctx context.Context, req *protobuf.DeleteUserURLsRequest) (*protobuf.DeleteResponse, error) {
	s.service.DeleteURLs(ctx, req.UrlIds)

	return &protobuf.DeleteResponse{
		Status: &protobuf.Status{
			Code:    http.StatusAccepted,
			Message: "URLs deletion accepted",
		},
	}, nil
}

// Ping проверяет состояние базы данных
func (s *Server) Ping(ctx context.Context, req *emptypb.Empty) (*protobuf.PingResponse, error) {
	err := s.psqlStoragePinger.Ping(ctx)
	if err != nil {
		return &protobuf.PingResponse{
			Status: &protobuf.Status{
				Code:    http.StatusInternalServerError,
				Message: "Database unavailable",
			},
		}, nil
	}

	return &protobuf.PingResponse{
		Status: &protobuf.Status{
			Code:    http.StatusOK,
			Message: "Database available",
		},
	}, nil
}

// GetStats получает статистику сервиса
func (s *Server) GetStats(ctx context.Context, req *emptypb.Empty) (*protobuf.StatsResponse, error) {
	stats, err := s.internalService.GetState(ctx)
	if err != nil {
		return &protobuf.StatsResponse{
			Response: &protobuf.StatsResponse_Error{
				Error: &protobuf.Error{
					Status: &protobuf.Status{
						Code:    http.StatusInternalServerError,
						Message: "Internal server error",
					},
				},
			},
		}, nil
	}

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
	}, nil
}
