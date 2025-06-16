package urlsnipper

import (
	"context"
	"testing"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	dump "github.com/DanilNaum/SnipURL/pkg/utils/dumper"
)

func BenchmarkSetURL(b *testing.B) {
	hasher := &hasherMock{
		HashFunc: func(s string) string {
			return "mockedHash"
		},
	}
	dumper := &dumperMock{
		AddFunc: func(record *dump.URLRecord) error {
			return nil
		},
	}
	storage := &urlStorageMock{
		SetURLFunc: func(ctx context.Context, id string, url string) (int, error) {
			return 1, nil
		},
	}
	service := NewURLSnipperService(storage, hasher, dumper, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.SetURL(context.Background(), "http://example.com")
	}
}

func BenchmarkGetURL(b *testing.B) {
	hasher := &hasherMock{
		HashFunc: func(s string) string {
			return "mockedHash"
		},
	}
	storage := &urlStorageMock{
		GetURLFunc: func(ctx context.Context, id string) (string, error) {
			return "http://example.com", nil
		},
	}
	service := NewURLSnipperService(storage, hasher, nil, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetURL(context.Background(), "mockedHash")
	}
}

func BenchmarkSetURLs(b *testing.B) {
	hasher := &hasherMock{
		HashFunc: func(s string) string {
			return "mockedHash"
		},
	}
	dumper := &dumperMock{
		AddFunc: func(record *dump.URLRecord) error {
			return nil
		},
	}
	storage := &urlStorageMock{
		SetURLsFunc: func(ctx context.Context, urls []*urlstorage.URLRecord) ([]*urlstorage.URLRecord, error) {
			return urls, nil
		},
	}
	service := NewURLSnipperService(storage, hasher, dumper, nil, nil)

	urls := []*SetURLsInput{
		{CorrelationID: "1", OriginalURL: "http://example.com"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.SetURLs(context.Background(), urls)
	}
}

func BenchmarkGetURLs(b *testing.B) {
	storage := &urlStorageMock{
		GetURLsFunc: func(ctx context.Context) ([]*urlstorage.URLRecord, error) {
			return []*urlstorage.URLRecord{
				{ShortURL: "mockedHash", OriginalURL: "http://example.com"},
			}, nil
		},
	}
	service := NewURLSnipperService(storage, nil, nil, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetURLs(context.Background())
	}
}

// func BenchmarkDeleteURLs(b *testing.B) {
// 	deleteService := &deleteServiceMock{
// 		DeleteFunc: func(userID string, ids []string) {
// 			// Mock implementation does nothing
// 		},
// 	}
// 	service := NewURLSnipperService(nil, nil, nil, deleteService, nil)

// 	ids := []string{"id1", "id2", "id3"}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		service.DeleteURLs(context.Background(), ids)
// 	}
// }
