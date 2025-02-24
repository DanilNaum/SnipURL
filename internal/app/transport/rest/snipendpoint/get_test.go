package snipendpoint

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnipEndpoint_get(t *testing.T) {
	type input struct {
		method string
		id     string
	}
	type mocks struct {
		getURLFunc              func(ctx context.Context, id string) (string, error)
		getURLFuncNumberOfCalls int
	}
	type want struct {
		code   int
		body   string
		header http.Header
	}
	tests := []struct {
		name  string
		input input
		mocks mocks
		want  want
	}{
		{
			name: "happy_path",

			input: input{
				method: http.MethodGet,
				id:     "123",
			},

			mocks: mocks{
				getURLFunc: func(ctx context.Context, id string) (string, error) {
					return "https://example.com", nil
				},
				getURLFuncNumberOfCalls: 1,
			},

			want: want{
				code:   http.StatusTemporaryRedirect,
				header: http.Header{"Location": []string{"https://example.com"}},
			},
		},
		{
			name: "wrong method",
			input: input{method: http.MethodPost,
				id: "123",
			},
			want: want{
				code: http.StatusMethodNotAllowed,
				body: "Only GET requests are allowed!",
			},
		},
		{
			name: "service error",
			input: input{
				method: http.MethodGet,
				id:     "123",
			},
			mocks: mocks{
				getURLFunc: func(ctx context.Context, id string) (string, error) {
					return "", errors.New("service error")
				},
				getURLFuncNumberOfCalls: 1,
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
				GetURLFunc: tt.mocks.getURLFunc,
			}

			endpoint := &snipEndpoint{
				service: mockService,
			}

			req := httptest.NewRequest(tt.input.method, "/"+tt.input.id, nil)
			w := httptest.NewRecorder()

			endpoint.get(w, req)

			require.Equal(t, tt.want.code, w.Code, "Expected status code %d, got %d", tt.want.code, w.Code)

			switch tt.want.code {
			case http.StatusTemporaryRedirect:
				for k, v := range tt.want.header {
					require.Equal(t, v, w.Header().Values(k), "Expected header %v, got %v", v, w.Header().Values(k))
				}
			default:
				require.Equal(t, strings.TrimSpace(tt.want.body), strings.TrimSpace(w.Body.String()), "Expected body %s, got %s", tt.want.body, w.Body.String())
			}

			require.Equal(t, tt.mocks.getURLFuncNumberOfCalls, len(mockService.GetURLCalls()))

		})
	}
}
