package util

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandString(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.RawURLEncoding.EncodeToString(b)

	return state, nil
}
