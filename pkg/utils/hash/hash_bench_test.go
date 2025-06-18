package hash

import (
	"testing"
)

func BenchmarkHasher_Hash(b *testing.B) {
	h := NewHasher(10)
	for i := 0; i < b.N; i++ {
		h.Hash("qwertyuioasdfghjkklzxcvbnm")
	}
}
