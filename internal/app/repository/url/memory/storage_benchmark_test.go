package memory

import (
	"context"
	"strconv"
	"testing"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	dump "github.com/DanilNaum/SnipURL/pkg/utils/dumper"
)

func BenchmarkStorage_SetURL(b *testing.B) {
	s := NewStorage()
	for i := 0; i < b.N; i++ {
		_, err := s.SetURL(context.Background(), "url"+strconv.Itoa(i), "https://example.com")
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkStorage_GetURL(b *testing.B) {
	s := NewStorage()
	for i := 0; i < 1000; i++ {
		s.SetURL(context.Background(), "url"+strconv.Itoa(i), "https://example.com")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.GetURL(context.Background(), "url"+strconv.Itoa(i%1000))
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkStorage_SetURLs(b *testing.B) {
	s := NewStorage()
	urls := make([]*urlstorage.URLRecord, 0, b.N)
	for i := 0; i < b.N; i++ {
		urls = append(urls, &urlstorage.URLRecord{ShortURL: "url" + strconv.Itoa(i), OriginalURL: "https://example.com"})
	}

	b.ResetTimer()
	_, err := s.SetURLs(context.Background(), urls)
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}
}

func BenchmarkStorage_GetURLs(b *testing.B) {
	s := NewStorage()
	ctx := context.WithValue(context.Background(), key, "userID")
	for i := 0; i < 1000; i++ {
		s.SetURL(ctx, "url"+strconv.Itoa(i), "https://example.com")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.GetURLs(ctx)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkStorage_DeleteURLs(b *testing.B) {
	s := NewStorage()
	for i := 0; i < 1000; i++ {
		s.SetURL(context.Background(), "url"+strconv.Itoa(i), "https://example.com")
	}

	ids := make([]string, 0, 1000)
	for i := 0; i < 1000; i++ {
		ids = append(ids, "url"+strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := s.DeleteURLs("userID", ids)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
func BenchmarkStorage_RestoreStorage(b *testing.B) {
	s := NewStorage()
	dumper := &mockDumper{
		records: make(chan dump.URLRecord, b.N),
	}

	for i := 0; i < b.N; i++ {
		dumper.records <- dump.URLRecord{ShortURL: "url" + strconv.Itoa(i), OriginalURL: "https://example.com"}
	}
	close(dumper.records)

	b.ResetTimer()
	err := s.RestoreStorage(dumper)
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}
}

type mockDumper struct {
	records chan dump.URLRecord
}

func (m *mockDumper) ReadAll() (chan dump.URLRecord, error) {
	return m.records, nil
}
