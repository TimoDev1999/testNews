package random

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString генерирует случайную строку заданной длины
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
