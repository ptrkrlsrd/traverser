package acache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const wantEncoded = "eyJhbGlhcyI6ImFsaWFzIiwicmVxdWVzdCI6eyJNZXRob2QiOiJHRVQiLCJVUkwiOnsiU2NoZW1lIjoiaHR0cHMiLCJPcGFxdWUiOiIiLCJVc2VyIjpudWxsLCJIb3N0IjoiZXhhbXBsZS5jb20iLCJQYXRoIjoiIiwiUmF3UGF0aCI6IiIsIkZvcmNlUXVlcnkiOmZhbHNlLCJSYXdRdWVyeSI6IiIsIkZyYWdtZW50IjoiIiwiUmF3RnJhZ21lbnQiOiIifSwiUHJvdG8iOiJIVFRQLzEuMSIsIlByb3RvTWFqb3IiOjEsIlByb3RvTWlub3IiOjEsIkhlYWRlciI6e30sIkJvZHkiOm51bGwsIkNvbnRlbnRMZW5ndGgiOjAsIlRyYW5zZmVyRW5jb2RpbmciOm51bGwsIkNsb3NlIjpmYWxzZSwiSG9zdCI6ImV4YW1wbGUuY29tIiwiRm9ybSI6bnVsbCwiUG9zdEZvcm0iOm51bGwsIk11bHRpcGFydEZvcm0iOm51bGwsIlRyYWlsZXIiOm51bGx9fQ"

func Test_id_ToKey(t *testing.T) {
	type fields struct {
		Alias string
		Body  []byte
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
				Body:  []byte("body"),
			},
			want:    wantEncoded,
			wantErr: false,
		},
		{
			name: "Creates base64 encoded string",
			fields: fields{
				Alias: "alias",
				Body:  []byte("body"),
			},
			want:    wantEncoded,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "https://example.com", nil)
			if err != nil {
				t.Fail()
			}

			id, err := NewID(tt.fields.Alias, req)
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
			fmt.Println(got)
		})
	}
}

func Test_encodeBase64String(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fail()
	}

	id, err := NewID("alias", req)
	if err != nil {
		t.Fail()
	}

	testData, _ := json.Marshal(id)

	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "Encodes base64 string",
			args: string(testData),
			want: wantEncoded,
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
