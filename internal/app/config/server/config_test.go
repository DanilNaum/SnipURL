package server

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer_NewConfig(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		env      map[string]string
		expected *config
	}{
		{
			name: "default_values",
			args: []string{},
			expected: &config{
				Host:    "localhost:8080",
				BaseURL: "http://localhost:8080",
			},
		},
		{
			name: "custom_host_and_base_url_from_flags",
			args: []string{"-a", "example.com:8443", "-b", "https://example.com:8443"},
			expected: &config{
				Host:    "example.com:8443",
				BaseURL: "https://example.com:8443",
			},
		},
		{
			name: "custom_host_and_base_url_from_env",
			env: map[string]string{
				"SERVER_ADDRESS": "example.com:8443",
				"BASE_URL":       "https://example.com:8443",
			},
			expected: &config{
				Host:    "example.com:8443",
				BaseURL: "https://example.com:8443",
			},
		},
		{
			name: "custom_host_and_base_url_from_flags_and_env",
			args: []string{"-a", "example_args.com:8443", "-b", "https://example_args.com:8443"},
			env: map[string]string{
				"SERVER_ADDRESS": "example.com:8443",
				"BASE_URL":       "https://example.com:8443",
			},
			expected: &config{
				Host:    "example.com:8443",
				BaseURL: "https://example.com:8443",
			},
		},
		{
			name: "custom_host_in_flags_and_base_url_in_env",
			args: []string{"-a", "example.com:8443"},
			env: map[string]string{
				"BASE_URL": "https://example.com:8443",
			},
			expected: &config{
				Host:    "example.com:8443",
				BaseURL: "https://example.com:8443",
			},
		},
		{
			name: "custom_host_in_env_and_base_url_in_flags",
			args: []string{"-b", "https://example.com:8443"},
			env: map[string]string{
				"SERVER_ADDRESS": "example.com:8443",
			},
			expected: &config{
				Host:    "example.com:8443",
				BaseURL: "https://example.com:8443",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ExitOnError)

			os.Args = append([]string{"cmd"}, tt.args...)

			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			result := NewConfig(&loggerMock{})

			require.Equal(t, tt.expected.Host, result.Host)
			require.Equal(t, tt.expected.BaseURL, result.BaseURL)
		})
	}
}
func TestServer_ConfigValidate(t *testing.T) {
	tests := []struct {
		name          string
		config        *config
		expectedError string
	}{
		{
			name: "valid_host_and_baseurl",
			config: &config{
				Host:    "example.com:8080",
				BaseURL: "https://example.com:8080",
			},
		},
		{
			name: "invalid_host_no_domain",
			config: &config{
				Host:    ":8080",
				BaseURL: "https://example.com",
			},

			expectedError: "invalid host: :8080",
		},
		{
			name: "invalid_host_special_chars",
			config: &config{
				Host:    "example@.com:8080",
				BaseURL: "https://example.com",
			},
			expectedError: "invalid host: example@.com:8080",
		},
		{
			name: "invalid_baseurl_no_protocol",
			config: &config{
				Host:    "example.com:8080",
				BaseURL: "example.com",
			},
			expectedError: "invalid base url: example.com",
		},
		{
			name: "invalid_baseurl_wrong_protocol",
			config: &config{
				Host:    "example.com:8080",
				BaseURL: "ftp://example.com",
			},
			expectedError: "invalid base url: ftp://example.com",
		},
		{
			name: "valid_host_with_subdomain",
			config: &config{
				Host:    "sub.example.com:8080",
				BaseURL: "http://sub.example.com:8080",
			},
		},
		{
			name: "valid_baseurl_with_path",
			config: &config{
				Host:    "example.com:8080",
				BaseURL: "https://example.com:8080/path/to/resource",
			},
		},
		{
			name: "host_and_baseurl_do_not_match",
			config: &config{
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

			tt.config.validate(loggerMock)

		})
	}
}
func TestServer_ConfigGetPrefix(t *testing.T) {
	tests := []struct {
		name         string
		config       *config
		expectedPath string
		expectError  bool
	}{
		{
			name: "empty_path",
			config: &config{
				BaseURL: "http://example.com",
			},
			expectedPath: "/",
			expectError:  false,
		},
		{
			name: "root_path",
			config: &config{
				BaseURL: "http://example.com/",
			},
			expectedPath: "/",
			expectError:  false,
		},
		{
			name: "single_level_path",
			config: &config{
				BaseURL: "http://example.com/api",
			},
			expectedPath: "/api",
			expectError:  false,
		},
		{
			name: "multi_level_path",
			config: &config{
				BaseURL: "http://example.com/api/v1/users",
			},
			expectedPath: "/api/v1/users",
			expectError:  false,
		},
		{
			name: "path_with_query_params",
			config: &config{
				BaseURL: "http://example.com/api?version=1",
			},
			expectedPath: "/api",
			expectError:  false,
		},
		{
			name: "path_with_fragment",
			config: &config{
				BaseURL: "http://example.com/api#section1",
			},
			expectedPath: "/api",
			expectError:  false,
		},
		{
			name: "invalid_url",
			config: &config{
				BaseURL: "://invalid-url",
			},
			expectedPath: "",
			expectError:  true,
		},
		{
			name: "encoded_path_segments",
			config: &config{
				BaseURL: "http://example.com/api/user%20space/data",
			},
			expectedPath: "/api/user space/data",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := tt.config.GetPrefix()

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPath, path)
			}
		})
	}
}
