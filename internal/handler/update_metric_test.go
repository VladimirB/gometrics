package handler_test

import (
	"net/http"
	"testing"

	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestPostMetricHandler(t *testing.T) {
	server := handler.CreateTestServer()
	defer server.Close()

	testTable := []struct {
		name       string
		url        string
		statusCode int
	}{
		{"success counter metric", "/update/counter/PollCount/100", http.StatusOK},
		{"success gauge metric", "/update/gauge/Alloc/100.55", http.StatusOK},
		{"incorrect counter metric value", "/update/counter/name/str", http.StatusBadRequest},
		{"incorrect gauge metric value", "/update/gauge/name/str", http.StatusBadRequest},
		{"unknown metric type", "/update/unknown/name/100", http.StatusBadRequest},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := handler.MakeTestRequest(t, server, http.MethodPost, tt.url, "")
			response.Body.Close()
			assert.Equal(t, tt.statusCode, response.StatusCode)
		})
	}
}

func TestPostMetricNoNameHandler(t *testing.T) {
	server := handler.CreateTestServer()
	defer server.Close()

	testTable := []struct {
		name       string
		url        string
		statusCode int
	}{
		{"no counter metric name", "/update/counter/100", http.StatusNotFound},
		{"no gauge metric name", "/update/gauge/100", http.StatusNotFound},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := handler.MakeTestRequest(t, server, http.MethodPost, tt.url, "")
			response.Body.Close()
			assert.Equal(t, tt.statusCode, response.StatusCode)
		})
	}
}
