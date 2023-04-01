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

func (ck *cacheKey) UnmarshalJSON(data []byte) error {
	var raw map[string]*json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	if alias, ok := raw["alias"]; ok {
		err = json.Unmarshal(*alias, &ck.Alias)
		if err != nil {
			return err
		}
	}

	if request, ok := raw["request"]; ok {
		err = json.Unmarshal(*request, &ck.Request)
		if err != nil {
			return err
		}
	}

	return nil
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

func (key cacheKey) ToKey() (string, error) {
	jsonString, err := json.Marshal(key)
	if err != nil {
		return "", err
	}

	encoded := encodeBase64String(string(jsonString))

	return encoded, nil
}

func encodeBase64String(str string) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(str))
}
