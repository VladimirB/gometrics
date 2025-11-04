package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
    req, err := http.NewRequest(method, ts.URL+path, nil)
    require.NoError(t, err)
						
    resp, err := ts.Client().Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    require.NoError(t, err)

    return resp, string(respBody)
}

func TestPostMetricHandler(t *testing.T) {
	server := httptest.NewServer(handler.NewRouter())
	defer server.Close()

	testTable := []struct {
		name string
		url string
		statusCode int
	}{
		{"success counter metric", "/update/counter/name/100", http.StatusOK},
		{"success gauge metric", "/update/gauge/name/100.55", http.StatusOK},
		{"incorrect counter metric value", "/update/counter/name/str", http.StatusBadRequest},
		{"incorrect gauge metric value", "/update/gauge/name/str", http.StatusBadRequest},
		{"unknown metric type", "/update/unknown/name/100", http.StatusBadRequest},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := testRequest(t, server, http.MethodPost, tt.url)
			assert.Equal(t, tt.statusCode, response.StatusCode)
		})
	}
}

func TestPostMetricNoNameHandler(t *testing.T) {
	server := httptest.NewServer(handler.NewRouter())
	defer server.Close()

	testTable := []struct {
		name string
		url string
		statusCode int
	}{
		{"no counter metric name", "/update/counter/100", http.StatusNotFound},
		{"no gauge metric name", "/update/gauge/100", http.StatusNotFound},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := testRequest(t, server, http.MethodPost, tt.url)
			assert.Equal(t, tt.statusCode, response.StatusCode)
		})
	}
}
