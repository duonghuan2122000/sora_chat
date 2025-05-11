package util

import (
	"crypto/rand"
	"encoding/base64"
)

// Tạo token reference
func GenerateReferenceToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
