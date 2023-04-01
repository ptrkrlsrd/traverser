package acache

import (
	"io"
	"net/http"
)

type StorableResponse struct {
	Status           string
	StatusCode       int
	Headers          map[string]string
	Body             string
	ContentLength    int64
	TransferEncoding []string
}

// NewStorableResponse maps a http.Response to a StorableResponse
func NewStorableResponse(httpResponse *http.Response) (StorableResponse, error) {
	defer httpResponse.Body.Close()
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return StorableResponse{}, err
	}

	headers := ToHeaderMap(httpResponse.Header)

	return StorableResponse{
		Status:           httpResponse.Status,
		StatusCode:       httpResponse.StatusCode,
		Headers:          headers,
		Body:             string(body),
		ContentLength:    httpResponse.ContentLength,
		TransferEncoding: httpResponse.TransferEncoding,
	}, nil
}

func ToHeaderMap(headers http.Header) map[string]string {
	headerMap := make(map[string]string)
	for k, v := range headers {
		headerMap[k] = v[0]
	}
	return headerMap
}
