package traverser

import (
	"io"
	"net/http"
)

type StorableResponse struct {
	Status           string            `json:"status,omitempty"`
	StatusCode       int               `json:"status_code,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
	Body             string            `json:"body,omitempty"`
	ContentLength    int64             `json:"content_length,omitempty"`
	TransferEncoding []string          `json:"transfer_encoding,omitempty"`
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
