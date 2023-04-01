package acache

import (
	"io"
	"net/http"
)

type StorableRequest struct {
	Method     string            `yaml:"method"`
	URL        string            `yaml:"url"`
	Proto      string            `yaml:"proto"`
	Headers    map[string]string `yaml:"headers"`
	Body       []byte            `yaml:"body"`
	Host       string            `yaml:"host"`
	RemoteAddr string            `yaml:"remoteAddr"`
	Response   StorableResponse  `yaml:"response"`
}

func NewStorableRequestWithResponse(req *http.Request, resp http.Response) (StorableRequest, error) {
	headers := make(map[string]string)
	for k, v := range req.Header {
		headers[k] = v[0]
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return StorableRequest{}, err
	}

	storableRequest := StorableRequest{
		Method:     req.Method,
		URL:        req.URL.String(),
		Proto:      req.Proto,
		Headers:    headers,
		Host:       req.Host,
		RemoteAddr: req.RemoteAddr,
		Response: StorableResponse{
			StatusCode: resp.StatusCode,
			Headers:    headers,
			Body:       string(body),
		},
	}

	if req.Body != nil {
		defer req.Body.Close()
		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			return StorableRequest{}, err
		}
		storableRequest.Body = requestBody
	}

	return storableRequest, nil
}

// NewStorableRequest maps a http.Request to a StorableRequest
func NewStorableRequest(req *http.Request) (StorableRequest, error) {
	headers := make(map[string]string)
	for k, v := range req.Header {
		headers[k] = v[0]
	}

	storableRequest := StorableRequest{
		Method:     req.Method,
		URL:        req.URL.String(),
		Proto:      req.Proto,
		Headers:    headers,
		Host:       req.Host,
		RemoteAddr: req.RemoteAddr,
	}

	if req.Body != nil {
		defer req.Body.Close()
		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			return StorableRequest{}, err
		}
		storableRequest.Body = requestBody
	}

	return storableRequest, nil
}
