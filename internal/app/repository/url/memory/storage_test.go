package memory

import (
	"context"
	"testing"

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
		startStorageState map[string]*urlstorage.URLRecord
		args              args
		storageStateAfter map[string]*urlstorage.URLRecord
		wantErr           error
	}{
		{
			name: "success_new_url",

			startStorageState: map[string]*urlstorage.URLRecord{},

			args: args{
				id:  "abc123",
				url: "https://example.com",
			},

			storageStateAfter: map[string]*urlstorage.URLRecord{
				"abc123": {OriginalURL: "https://example.com"},
			},

			wantErr: nil,
		},
		{
			name: "success_same_url",

			startStorageState: map[string]*urlstorage.URLRecord{
				"abc123": {OriginalURL: "https://example.com"},
			},

			args: args{
				id:  "abc123",
				url: "https://example.com",
			},

			storageStateAfter: map[string]*urlstorage.URLRecord{
				"abc123": {OriginalURL: "https://example.com"},
			},
			wantErr: urlstorage.ErrConflict,
		},
		{
			name: "error_id_busy",

			startStorageState: map[string]*urlstorage.URLRecord{
				"abc123": {OriginalURL: "https://example.com"},
			},

			args: args{
				id:  "abc123",
				url: "https://different.com",
			},

			storageStateAfter: map[string]*urlstorage.URLRecord{
				"abc123": {OriginalURL: "https://example.com"},
			},
			wantErr: urlstorage.ErrIDIsBusy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				urls: tt.startStorageState,
			}

			length, err := s.SetURL(context.Background(), tt.args.id, tt.args.url)
			require.ErrorIs(t, err, tt.wantErr)
			if err == nil {
				require.Equal(t, length, len(tt.storageStateAfter))
			}

		})
	}
}
func TestStorage_GetURL(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name              string
		startStorageState map[string]*urlstorage.URLRecord
		id                string
		want              string
		wantErr           error
	}{
		{
			name: "success_get_existing_url",
			startStorageState: map[string]*urlstorage.URLRecord{
				"abc123": {OriginalURL: "https://example.com"},
			},

			id: "abc123",

			want:    "https://example.com",
			wantErr: nil,
		},
		{
			name:              "error_not_found",
			startStorageState: map[string]*urlstorage.URLRecord{},

			id: "nonexistent",

			want:    "",
			wantErr: urlstorage.ErrNotFound,
		},
		{
			name: "success_get_with_multiple_urls",
			startStorageState: map[string]*urlstorage.URLRecord{
				"abc123": {OriginalURL: "https://example.com"},
				"def456": {OriginalURL: "https://test.com"},
				"ghi789": {OriginalURL: "https://sample.com"},
			},

			id: "def456",

			want:    "https://test.com",
			wantErr: nil,
		},
		{
			name: "error_empty_id",
			startStorageState: map[string]*urlstorage.URLRecord{
				"abc123": {OriginalURL: "https://example.com"},
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
			}

			got, err := s.GetURL(context.Background(), tt.id)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}
