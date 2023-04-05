package acache

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type cacheKey struct {
	Alias   string          `json:"alias"`
	Request StorableRequest `json:"request"`
}

func NewCacheKey(alias string, request *http.Request) (cacheKey, error) {
	storableRequest, err := NewStorableRequest(request)
	if err != nil {
		return cacheKey{}, err
	}

	return cacheKey{
		Alias:   alias,
		Request: storableRequest,
	}, nil
}

func (self cacheKey) ToKey() (string, error) {
	encoded := encodeBase64String([]byte(self.Alias + self.Request.URL))

	return encoded, nil
}

func CacheKeyFromKey(key string) (cacheKey, error) {
	decodedString, err := decodeBase64String(key)
	if err != nil {
		return cacheKey{}, err
	}

    var c cacheKey

	err = json.Unmarshal(decodedString, &c)
	if err != nil {
		return cacheKey{}, err
	}

	return c, nil
}

func encodeBase64String(str []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(str)
}

func decodeBase64String(str string) ([]byte, error) {
	return base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(str)
}
