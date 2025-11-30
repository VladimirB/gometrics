package handler_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/VladimirB/gometrics/internal/handler"
	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/VladimirB/gometrics/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestValueMetricHandler_GetMetricHandler(t *testing.T) {
	storage := repository.NewMemStorage()
	service := service.NewMetricsService(storage)
	service.SaveByFields(models.Counter, models.PollCount, 100)
	service.SaveByFields(models.Gauge, models.Alloc, 3.14)

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
			response, body := handler.MakeTestRequest(t, server, http.MethodGet, tt.url, "")
			defer response.Body.Close()
			assert.Equal(t, tt.statusCode, response.StatusCode)
			if response.StatusCode == http.StatusOK {
				assert.Equal(t, tt.body, body)
			}
		})
	}
}

func TestValueMetricHandler_PostValueMetricHandler(t *testing.T) {
	storage := repository.NewMemStorage()
	service := service.NewMetricsService(storage)
	data := prepareTestData(service)

	server := handler.CreateTestServerWithStorage(storage)
	defer server.Close()

	tests := []struct {
		name       string
		body       string
		statusCode int
		response   string
	}{
		{"ask no value metric", `{"id":"AnyID", "type":"counter"}`, http.StatusOK, `{"id":"AnyID", "type":"counter"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, body := handler.MakeTestRequest(t, server, http.MethodPost, "/value", tt.body)
			defer response.Body.Close()

			assert.Equal(t, tt.statusCode, response.StatusCode)
			if response.StatusCode == http.StatusOK {
				assert.JSONEq(t, tt.response, body)
			}
		})
	}

	for metricID := range models.AllowedMetrics {
		testName := fmt.Sprintf("POST value %s", metricID)
		t.Run(testName, func(t *testing.T) {
			response, body := handler.MakeTestRequest(t, server, http.MethodPost, "/value", requestBody(metricID))
			defer response.Body.Close()

			assert.Equal(t, http.StatusOK, response.StatusCode)
			if response.StatusCode == http.StatusOK {
				assert.JSONEq(t, expectedResponse(metricID, data), body)
			}
		})
	}
}

func prepareTestData(service *service.MetricsService) map[string]float64 {
	var result = make(map[string]float64)

	for metricID := range models.AllowedMetrics {
		if metricID == models.PollCount {
			result[metricID] = float64(rand.Int())
			service.SaveByFields(models.Counter, metricID, result[metricID])
		} else {
			result[metricID] = rand.Float64()
			service.SaveByFields(models.Gauge, metricID, result[metricID])
		}
	}

	return result
}

func requestBody(metricID string) string {
	if metricID == models.PollCount {
		return fmt.Sprintf(`{"id":"%s", "type":"counter"}`, metricID)
	} else {
		return fmt.Sprintf(`{"id":"%s", "type":"gauge"}`, metricID)
	}
}

func expectedResponse(metricID string, data map[string]float64) string {
	if metricID == models.PollCount {
		return fmt.Sprintf(`{"id":"%s", "type":"counter", "delta":%d}`, metricID, int64(data[metricID]))
	} else {
		return fmt.Sprintf(`{"id":"%s", "type":"gauge", "value":%v}`, metricID, data[metricID])
	}
}
