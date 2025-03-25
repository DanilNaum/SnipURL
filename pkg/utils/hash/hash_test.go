package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasher_Hash(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		input    string
		expected int
	}{
		{
			name:     "empty string",
			length:   8,
			input:    "",
			expected: 8,
		},
		{
			name:     "normal string",
			length:   16,
			input:    "test string",
			expected: 16,
		},
		{
			name:     "special characters",
			length:   12,
			input:    "!@#$%^&*()",
			expected: 12,
		},
		{
			name:     "very long string",
			length:   10,
			input:    "this is a very long string that needs to be hashed properly",
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHasher(tt.length)
			result := h.Hash(tt.input)

			require.Equal(t, tt.expected, len(result))

			// Test idempotency
			result2 := h.Hash(tt.input)
			require.Equal(t, result, result2)
		})
	}
}
