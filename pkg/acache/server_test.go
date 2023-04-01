package acache

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	server Server
)

type MockStorage struct {
	routes Routes
}

func (ms *MockStorage) GetRoutes() (routes Routes, err error) {
	return routes, nil
}

func (ms *MockStorage) AddRoute(route Route) error {
	ms.routes = append(ms.routes, route)
	return nil
}

func (ms *MockStorage) Clear() error {
	return nil
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	gin.SetMode(gin.ReleaseMode) // Release mode to make gin less verbose
	router := gin.New()          // TODO: Figure out if its okay to use New here instead of Default

	testResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("Hello world!")),
		Header:     http.Header{},
	}

	testResponse.Header.Add("Content-Type", "application/json; charset=utf-8")

	route, err := NewRouteFromResponse("/test", "/alias/test", http.MethodGet, testResponse)
	if err != nil {
		t.Fail()
	}

	routes := Routes{route}
	mockStorage := &MockStorage{}

	server = NewServer(mockStorage, router)
	server.RegisterRoutes(routes)

	return func(t *testing.T) {
		t.Log("teardown test case")
	}
}

func TestOKStatusCode(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/alias/test", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}

func TestNotOKStatusCode(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/does-not-exist", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}

func TestHeaders(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/alias/test", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	val, ok := resp.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}

	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}
