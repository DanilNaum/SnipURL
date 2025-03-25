package hash

import (
	"crypto/md5"
	"encoding/hex"
)

type hasher struct {
	length int
}

func NewHasher(length int) *hasher {
	return &hasher{
		length: length,
	}
}

func (h *hasher) Hash(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:h.length]
}
