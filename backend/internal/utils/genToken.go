package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateNewToken(lenght int) (string, error) {
	b := make([]byte, lenght)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
