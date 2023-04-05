package acache

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var num = 10000

func BenchmarkReverseString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewCacheKey("/abc", &http.Request{
            Method: http.MethodGet,
            URL: &url.URL{
                Scheme: "http",
                Host:   "example.com",
                Path:   "/test",
            },
            GetBody: func() (io.ReadCloser, error) {
                return io.NopCloser(strings.NewReader("Hello, World!")), nil
            },
            ContentLength:    13,
            TransferEncoding: []string{"identity"},
            Close:            false,
            Host:             "example.com",
            Form: map[string][]string{
                "key1": {"value1"},
                "key2": {"value2"},
            },
            PostForm: map[string][]string{
                "postKey1": {"postValue1"},
                "postKey2": {"postValue2"},
            },
            MultipartForm: nil,
            Trailer: map[string][]string{
                "Trailer-Key": {"Trailer-Value"},
            },
            RemoteAddr: "192.0.2.1:12345",
            RequestURI: "/test?param=value",
            TLS:        nil,
            Cancel:     nil,
            Response: &http.Response{
                Status:           "200 OK",
                StatusCode:       200,
                Proto:            "HTTP/1.1",
                ProtoMajor:       1,
                ProtoMinor:       1,
                Header:           http.Header{"Content-Type": {"application/json"}},
                Body:             ioutil.NopCloser(strings.NewReader(`{"status": "success"}`)),
                ContentLength:    21,
                TransferEncoding: []string{"identity"},
                Close:            false,
                Uncompressed:     false,
                Trailer:          http.Header{"Content-Encoding": {"gzip"}},
                Request:          nil,
                TLS:              nil,
            },
        })
    }
}
