package crypto

import (
	"crypto/rand"
	"io"
)

func GenerateRandomKey(length int) []byte {
	key := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil
	}
	return key
}
