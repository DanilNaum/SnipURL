package server

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/require"
)

type conf struct {
	Host    string
	BaseURL string
}

func TestServer_ConfigValidate(t *testing.T) {
	tests := []struct {
		name          string
		config        *conf
		expectedError string
	}{
		{
			name: "valid_host_and_baseurl",
			config: &conf{
				Host:    "example.com:8080",
				BaseURL: "https://example.com:8080",
			},
		},
		{
			name: "invalid_host_no_domain",
			config: &conf{
				Host:    ":8080",
				BaseURL: "https://example.com",
			},

			expectedError: "invalid host: :8080",
		},
		{
			name: "invalid_host_special_chars",
			config: &conf{
				Host:    "example@.com:8080",
				BaseURL: "https://example.com",
			},
			expectedError: "invalid host: example@.com:8080",
		},
		{
			name: "invalid_baseurl_no_protocol",
			config: &conf{
				Host:    "example.com:8080",
				BaseURL: "example.com",
			},
			expectedError: "invalid base url: example.com",
		},
		{
			name: "invalid_baseurl_wrong_protocol",
			config: &conf{
				Host:    "example.com:8080",
				BaseURL: "ftp://example.com",
			},
			expectedError: "invalid base url: ftp://example.com",
		},
		{
			name: "valid_host_with_subdomain",
			config: &conf{
				Host:    "sub.example.com:8080",
				BaseURL: "http://sub.example.com:8080",
			},
		},
		{
			name: "valid_baseurl_with_path",
			config: &conf{
				Host:    "example.com:8080",
				BaseURL: "https://example.com:8080/path/to/resource",
			},
		},
		{
			name: "host_and_baseurl_do_not_match",
			config: &conf{
				Host:    "example.com:8080",
				BaseURL: "https://example.com:8081",
			},
			expectedError: "base url https://example.com:8081 must contain host example.com:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerMock := &loggerMock{
				FatalfFunc: func(format string, v ...any) {
					err := fmt.Sprintf(format, v...)
					require.Equal(t, tt.expectedError, err)
				},
			}
			config := &serverConfig{
				Host:    &tt.config.Host,
				BaseURL: &tt.config.BaseURL,
			}
			config.ValidateServerConfig(loggerMock)

		})
	}
}
func TestServer_ConfigGetPrefix(t *testing.T) {
	tests := []struct {
		name         string
		config       *conf
		expectedPath string
		expectError  bool
	}{
		{
			name: "empty_path",
			config: &conf{
				BaseURL: "http://example.com",
			},
			expectedPath: "/",
			expectError:  false,
		},
		{
			name: "root_path",
			config: &conf{
				BaseURL: "http://example.com/",
			},
			expectedPath: "/",
			expectError:  false,
		},
		{
			name: "single_level_path",
			config: &conf{
				BaseURL: "http://example.com/api",
			},
			expectedPath: "/api",
			expectError:  false,
		},
		{
			name: "multi_level_path",
			config: &conf{
				BaseURL: "http://example.com/api/v1/users",
			},
			expectedPath: "/api/v1/users",
			expectError:  false,
		},
		{
			name: "path_with_query_params",
			config: &conf{
				BaseURL: "http://example.com/api?version=1",
			},
			expectedPath: "/api",
			expectError:  false,
		},
		{
			name: "path_with_fragment",
			config: &conf{
				BaseURL: "http://example.com/api#section1",
			},
			expectedPath: "/api",
			expectError:  false,
		},
		{
			name: "invalid_url",
			config: &conf{
				BaseURL: "://invalid-url",
			},
			expectedPath: "",
			expectError:  true,
		},
		{
			name: "encoded_path_segments",
			config: &conf{
				BaseURL: "http://example.com/api/user%20space/data",
			},
			expectedPath: "/api/user space/data",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &serverConfig{
				BaseURL: &tt.config.BaseURL,
			}
			path, err := config.GetPrefix()

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPath, path)
			}
		})
	}
}
