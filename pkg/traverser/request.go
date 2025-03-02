package traverser

import (
	"io"
	"net/http"
)

type StorableRequest struct {
	Method     string            `json:"method,omitempty"`
	URL        string            `json:"url,omitempty"`
	Proto      string            `json:"proto,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       []byte            `json:"body,omitempty"`
	Host       string            `json:"host,omitempty"`
	RemoteAddr string            `json:"remote_addr,omitempty"`
}

func NewStorableRequestWithResponse(req *http.Request, resp http.Response) (StorableRequest, error) {
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
