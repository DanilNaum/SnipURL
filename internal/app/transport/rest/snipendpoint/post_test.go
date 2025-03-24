package snipendpoint

import (
	"context"
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
		body      string
		bodyError bool
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
				body: "https://example.com",
				host: "http://localhost:8080",
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
				body: "https://example.com",
				host: "https://localhost:8080",
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
			name: "service_error",
			input: input{
				body: "https://example.com",
				host: "https://localhost:8080",
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
				bodyError: true,
				host:      "https://localhost:8080",
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "Internal Server Error",
			},
		},
		{
			name: "host_error",
			input: input{
				body: "https://example.com",
				host: "asd",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &serviceMock{
				SetURLFunc: tt.mocks.setURLFunc,
			}

			baseURL := tt.input.host

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

			req := httptest.NewRequest(http.MethodPost, "/", bodyReader)

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
