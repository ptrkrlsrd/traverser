package acache

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStorableResponse(t *testing.T) {
	testBody := "Hello world!"
	status := "200 OK"
	testResponse := &http.Response{
		Status:     status,
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(testBody)),
		Header:     http.Header{},
	}

	testResponse.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := NewStorableResponse(testResponse)
	if err != nil {
		t.Fail()
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	if resp.Status != status {
		t.Fatalf("Expected status '200 ok' got %v", resp.Status)
	}

	if string(resp.Body) != testBody {
		t.Fatalf("Expected status '%s' got %v", testBody, resp.Body)
	}
}

func TestIfCanWriteStorableResponse(t *testing.T) {
	status := "200 OK"
	testResponse := &http.Response{
		Status:     status,
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("Hello world!")),
		Header:     http.Header{},
	}

	testResponse.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := NewStorableResponse(testResponse)
	if err != nil {
		t.Fail()
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		t.Fail()
	}

	var readResp StorableResponse
	err = json.Unmarshal(respData, &readResp)
	if err != nil {
		t.Fail()
	}
}
