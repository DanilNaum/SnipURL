package urlsnipper

import (
	"context"
	"errors"
	"testing"

	dump "github.com/DanilNaum/SnipURL/pkg/utils/dumper"
	"github.com/stretchr/testify/require"
)

func TestUrlSnipperService_SetURL(t *testing.T) {
	tests := []struct {
		name                string
		url                 string
		hashResults         []string
		hashFuncGenerator   func() func(s string) string
		setURLFuncGenerator func() func(ctx context.Context, id string, url string) (int, error)
		setURLResults       []error

		want    string
		wantErr error
	}{
		{
			name: "successful first attempt",
			url:  "http://example.com",
			hashFuncGenerator: func() func(s string) string {
				i := 0
				return func(s string) string {
					switch i {
					case 0:
						i++
						return "abc123"
					default:
						t.Error("unexpected call to hash function")
						return ""

					}
				}
			},
			setURLFuncGenerator: func() func(ctx context.Context, id string, url string) (int, error) {
				i := 0
				return func(ctx context.Context, id string, url string) (int, error) {
					switch i {
					case 0:
						i++
						return 1, nil
					default:
						t.Error("unexpected call to setURL function")
						return -1, nil
					}
				}
			},
			want: "abc123",
		},
		{
			name: "success after collision",
			url:  "http://example.com",

			hashFuncGenerator: func() func(s string) string {
				i := 0
				return func(s string) string {
					switch i {
					case 0:
						i++
						return "abc123"
					case 1:
						i++
						return "def456"
					default:
						t.Error("unexpected call to hash function")
						return ""

					}
				}
			},
			setURLFuncGenerator: func() func(ctx context.Context, id string, url string) (int, error) {
				i := 0
				return func(ctx context.Context, id string, url string) (int, error) {
					switch i {
					case 0:
						i++
						return -1, errors.New("collision")
					case 1:
						i++
						return 1, nil
					default:
						t.Error("unexpected call to setURL function")
						return -1, nil
					}
				}
			},
			want: "def456",
		},
		{
			name: "max attempts reached",
			url:  "http://example.com",
			hashFuncGenerator: func() func(s string) string {
				i := 0
				return func(s string) string {
					if i < _maxAttempts {
						i++
						return "abc123"
					} else {
						t.Error("unexpected call to hash function")
						return ""

					}
				}
			},
			setURLFuncGenerator: func() func(ctx context.Context, id string, url string) (int, error) {
				i := 0
				return func(ctx context.Context, id string, url string) (int, error) {
					if i < _maxAttempts {
						i++
						return -1, errors.New("collision")
					} else {
						t.Error("unexpected call to setURL function")
						return -1, nil
					}
				}
			},

			wantErr: ErrFailedToGenerateID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHasher := &hasherMock{
				HashFunc: tt.hashFuncGenerator(),
			}

			mockStorage := &urlStorageMock{
				SetURLFunc: tt.setURLFuncGenerator(),
			}

			mockDumper := &dumperMock{
				AddFunc: func(record *dump.URLRecord) error {
					return nil
				},
			}

			s := &urlSnipperService{
				hasher:  mockHasher,
				storage: mockStorage,
				dumper:  mockDumper,
			}

			got, err := s.SetURL(context.Background(), tt.url)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)

		})
	}
}
func TestUrlSnipperService_GetURL(t *testing.T) {
	tests := []struct {
		name                    string
		id                      string
		getURLFunc              func(ctx context.Context, id string) (string, error)
		getURLFuncNumberOfCalls int
		want                    string
		wantErr                 error
	}{
		{
			name: "successful url retrieval",
			id:   "abc123",
			getURLFunc: func(ctx context.Context, id string) (string, error) {
				return "http://example.com", nil

			},
			getURLFuncNumberOfCalls: 1,
			want:                    "http://example.com",
		},
		{
			name: "storage error",
			id:   "def456",
			getURLFunc: func(ctx context.Context, id string) (string, error) {
				return "", errors.New("storage error")

			},
			getURLFuncNumberOfCalls: 1,
			wantErr:                 ErrFailedToGetURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &urlStorageMock{
				GetURLFunc: tt.getURLFunc,
			}

			s := &urlSnipperService{
				storage: mockStorage,
			}

			got, err := s.GetURL(context.Background(), tt.id)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.getURLFuncNumberOfCalls, len(mockStorage.GetURLCalls()))
		})
	}
}
