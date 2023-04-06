package acache

import (
	"encoding/json"
	"testing"
)

const wantEncoded = "YWxpYXM"

func Test_Id_ToKey(t *testing.T) {
	type fields struct {
		Alias string
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
            },
			want:    wantEncoded,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewCacheKey(tt.fields.Alias)
			if err != nil {
				t.Fail()
			}

			got := id
			if got != tt.want {
				t.Errorf("id = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeBase64String(t *testing.T) {
	id, err := NewCacheKey("alias")
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
			want: "IllXeHBZWE0i",
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
