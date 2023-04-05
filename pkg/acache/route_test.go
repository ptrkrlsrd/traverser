package acache

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const wantEncoded = "YWxpYXNodHRwOi8vZXhhbXBsZS5jb20vdGVzdA"

func Test_Id_ToKey(t *testing.T) {
	type fields struct {
		Alias string
		Request http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Creates base64 encoded string",
            fields: fields{
                Alias: "alias",
                Request: http.Request{
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
                },
            },
			want:    wantEncoded,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewCacheKey(tt.fields.Alias, &tt.fields.Request)
			if err != nil {
				t.Fail()
			}

			got, err := id.ToKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("id.ToKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("id.ToKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeBase64String(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fail()
	}

	id, err := NewCacheKey("alias", req)
	if err != nil {
		t.Fail()
	}

	testData, _ := json.Marshal(id)

	tests := []struct {
		name string
		args []byte
		want string
	}{
		{
			name: "Encodes base64 string",
			args: testData,
			want: "eyJhbGlhcyI6ImFsaWFzIiwicmVxdWVzdCI6eyJtZXRob2QiOiJHRVQiLCJ1cmwiOiJodHRwczovL2V4YW1wbGUuY29tIiwicHJvdG8iOiJIVFRQLzEuMSIsImhvc3QiOiJleGFtcGxlLmNvbSJ9fQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeBase64String(tt.args); got != tt.want {
				t.Errorf("encodeBase64String() = %v, want %v", got, tt.want)
			}
		})
	}
}
