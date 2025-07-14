package hash

import (
	"crypto/md5"
	"encoding/hex"
)

type hasher struct {
	length int
}

// NewHasher creates a new hasher with a specified length for truncating hash values.
// The length parameter determines how many characters of the MD5 hash will be returned.
func NewHasher(length int) *hasher {
	return &hasher{
		length: length,
	}
}

// Hash generates an MD5 hash of the input string and returns a truncated version
// based on the predefined length. It converts the input to bytes, computes the MD5 hash,
// and returns the first 'length' characters of the hexadecimal representation.
func (h *hasher) Hash(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:h.length]
}
