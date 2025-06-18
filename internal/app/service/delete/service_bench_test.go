package deleteurl

import (
	"context"
	"testing"
)

func BenchmarkDelete(b *testing.B) {
	service := NewDeleteService(context.Background(), &urlStorageMock{
		DeleteURLsFunc: func(userID string, ids []string) error {
			return nil
		},
	})
	userID := "user1"
	ids := []string{"id1", "id2", "id3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Delete(userID, ids)
	}
}
