package cookie

import (
	"net/http/httptest"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	manager := NewCookieManager([]byte("secret"))

	w := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Set(w, "testValue")
	}
}

func BenchmarkCreateSignature(b *testing.B) {
	manager := NewCookieManager([]byte("secret"))
	data := []byte("testValue")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.createSignature(data)
	}
}

func BenchmarkVerifySignature(b *testing.B) {
	manager := NewCookieManager([]byte("secret"))
	data := []byte("testValue")
	signature := manager.createSignature(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.verifySignature(data, signature)
	}
}
