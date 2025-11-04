package handler_test

import (
	"net/http"
	"testing"

	"github.com/VladimirB/gometrics/internal/handler"
	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestValueMetricHandler_GetMetricHandler(t *testing.T) {
	storage := repository.NewMemStorage()
	storage.Save(models.PollCount, createMetric(models.Counter, models.PollCount, 100))
	storage.Save(models.Alloc, createMetric(models.Gauge, models.Alloc, 3.14))

	server := handler.CreateTestServerWithStorage(storage)
	defer server.Close()

	testTable := []struct {
		name       string
		url        string
		statusCode int
		body       string
	}{
		{"success counter metric get", "/value/counter/PollCount", http.StatusOK, "100"},
		{"success gauge metric get", "/value/gauge/Alloc", http.StatusOK, "3.14"},
		{"metric not found", "/value/counter/Unknown", http.StatusNotFound, ""},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			response, body := handler.MakeTestRequest(t, server, http.MethodGet, tt.url)
			assert.Equal(t, tt.statusCode, response.StatusCode)
			if response.StatusCode == http.StatusOK {
				assert.Equal(t, tt.body, body)
			}
		})
	}
}

func createMetric(metricType string, metricName string, value float64) models.Metrics {
	metric := models.Metrics{
		ID:    metricName,
		MType: metricType,
		Value: new(float64),
	}
	metric.Value = &value
	return metric
}
