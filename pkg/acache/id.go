package acache

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type id struct {
	Alias   string          `json:"alias"`
	Request StorableRequest `json:"request"`
}

func NewID(alias string, request *http.Request) (id, error) {
	storableRequest, err := NewStorableRequest(request)
	if err != nil {
		return id{}, err
	}

	return id{
		Alias:   alias,
		Request: storableRequest,
	}, nil
}

func (id id) ToKey() (string, error) {
	jsonString, err := json.Marshal(id)
	if err != nil {
		return "", err
	}

	encoded := encodeBase64String(string(jsonString))

	return encoded, nil
}

func encodeBase64String(str string) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(str))
}
