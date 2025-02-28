package snipendpoint

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}
func TestSnipEndpoint_post(t *testing.T) {
	type input struct {
		method    string
		body      string
		bodyError bool
		isHTTPS   bool
		host      string
	}
	type mocks struct {
		setURLFunc              func(ctx context.Context, url string) (string, error)
		setURLFuncNumberOfCalls int
	}
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name  string
		input input
		mocks mocks
		want  want
	}{
		{
			name: "happy_path_http",
			input: input{
				method:  http.MethodPost,
				body:    "https://example.com",
				isHTTPS: false,
				host:    "localhost:8080",
			},
			mocks: mocks{
				setURLFunc: func(ctx context.Context, url string) (string, error) {
					return "abc123", nil
				},
				setURLFuncNumberOfCalls: 1,
			},
			want: want{
				code: http.StatusCreated,
				body: "http://localhost:8080/abc123",
			},
		},
		{
			name: "happy_path_https",
			input: input{
				method:  http.MethodPost,
				body:    "https://example.com",
				isHTTPS: true,
				host:    "localhost:8080",
			},
			mocks: mocks{
				setURLFunc: func(ctx context.Context, url string) (string, error) {
					return "abc123", nil
				},
				setURLFuncNumberOfCalls: 1,
			},
			want: want{
				code: http.StatusCreated,
				body: "https://localhost:8080/abc123",
			},
		},
		{
			name: "wrong_method",
			input: input{
				method: http.MethodGet,
				body:   "https://example.com",
			},
			want: want{
				code: http.StatusMethodNotAllowed,
				body: "Only POST requests are allowed!",
			},
		},
		{
			name: "service_error",
			input: input{
				method: http.MethodPost,
				body:   "https://example.com",
				host:   "localhost:8080",
			},
			mocks: mocks{
				setURLFunc: func(ctx context.Context, url string) (string, error) {
					return "", errors.New("service error")
				},
				setURLFuncNumberOfCalls: 1,
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "Internal Server Error",
			},
		},
		{
			name: "body_error",
			input: input{
				method:    http.MethodPost,
				bodyError: true,
				host:      "localhost:8080",
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "Internal Server Error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &serviceMock{
				SetURLFunc: tt.mocks.setURLFunc,
			}
			baseURL := "http://"
			if tt.input.isHTTPS {
				baseURL = "https://"
			}
			baseURL += tt.input.host
			endpoint := &snipEndpoint{
				service: mockService,
				prefix:  "/",
				baseURL: baseURL,
			}
			var bodyReader io.Reader
			bodyReader = strings.NewReader(tt.input.body)
			if tt.input.bodyError {
				bodyReader = errorReader{}
			}

			req := httptest.NewRequest(tt.input.method, "/", bodyReader)
			if tt.input.isHTTPS {
				req.TLS = &tls.ConnectionState{}
			}
			if tt.input.host != "" {
				req.Host = tt.input.host
			}
			w := httptest.NewRecorder()

			endpoint.post(w, req)

			require.Equal(t, tt.want.code, w.Code)
			require.Equal(t, strings.TrimSpace(tt.want.body), strings.TrimSpace(w.Body.String()))
			require.Equal(t, tt.mocks.setURLFuncNumberOfCalls, len(mockService.SetURLCalls()))
		})
	}
}
