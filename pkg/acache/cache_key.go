package acache

import (
	"encoding/base64"
)

func NewCacheKey(alias string) (string, error) {
	encoded := encodeBase64String([]byte(alias))

	return encoded, nil
}

func CacheKeyFromKey(key string) (string, error) {
	decodedString, err := decodeBase64String(key)
	if err != nil {
		return "", err
	}

	return string(decodedString), nil
}

func encodeBase64String(str []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(str)
}

func decodeBase64String(str string) ([]byte, error) {
	return base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(str)
}
