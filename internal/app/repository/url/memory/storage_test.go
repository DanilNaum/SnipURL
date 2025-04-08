package memory

import (
	"context"
	"sync"
	"testing"

	"github.com/DanilNaum/SnipURL/internal/app/repository/url"
	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"

	"github.com/stretchr/testify/require"
)

func TestStorage_SetURL(t *testing.T) {

	type args struct {
		id  string
		url string
	}
	tests := []struct {
		name              string
		startStorageState map[string]string
		args              args
		storageStateAfter map[string]string
		wantErr           error
	}{
		{
			name: "success_new_url",

			startStorageState: map[string]string{},

			args: args{
				id:  "abc123",
				url: "https://example.com",
			},

			storageStateAfter: map[string]string{
				"abc123": "https://example.com",
			},

			wantErr: nil,
		},
		{
			name: "success_same_url",

			startStorageState: map[string]string{
				"abc123": "https://example.com",
			},

			args: args{
				id:  "abc123",
				url: "https://example.com",
			},

			storageStateAfter: map[string]string{
				"abc123": "https://example.com",
			},
			wantErr: url.ErrConflict,
		},
		{
			name: "error_id_busy",

			startStorageState: map[string]string{
				"abc123": "https://example.com",
			},

			args: args{
				id:  "abc123",
				url: "https://different.com",
			},

			storageStateAfter: map[string]string{
				"abc123": "https://example.com",
			},
			wantErr: urlstorage.ErrIDIsBusy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				urls: tt.startStorageState,
				mu:   sync.RWMutex{},
			}

			length, err := s.SetURL(context.Background(), tt.args.id, tt.args.url)
			require.ErrorIs(t, err, tt.wantErr)
			if err == nil {
				require.Equal(t, length, len(tt.storageStateAfter))
			}
			require.Equal(t, tt.storageStateAfter, s.urls)

		})
	}
}
func TestStorage_GetURL(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name              string
		startStorageState map[string]string
		id                string
		want              string
		wantErr           error
	}{
		{
			name: "success_get_existing_url",
			startStorageState: map[string]string{
				"abc123": "https://example.com",
			},

			id: "abc123",

			want:    "https://example.com",
			wantErr: nil,
		},
		{
			name:              "error_not_found",
			startStorageState: map[string]string{},

			id: "nonexistent",

			want:    "",
			wantErr: urlstorage.ErrNotFound,
		},
		{
			name: "success_get_with_multiple_urls",
			startStorageState: map[string]string{
				"abc123": "https://example.com",
				"def456": "https://test.com",
				"ghi789": "https://sample.com",
			},

			id: "def456",

			want:    "https://test.com",
			wantErr: nil,
		},
		{
			name: "error_empty_id",
			startStorageState: map[string]string{
				"abc123": "https://example.com",
			},

			id: "",

			want:    "",
			wantErr: urlstorage.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				urls: tt.startStorageState,
				mu:   sync.RWMutex{},
			}

			got, err := s.GetURL(context.Background(), tt.id)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}
