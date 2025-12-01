package handler_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"testing"

	"github.com/VladimirB/gometrics/internal/handler"
	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/go-resty/resty/v2"

	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/VladimirB/gometrics/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestUpdateMetricHandler_UpdateMetricJSONHandler(t *testing.T) {
	storage := repository.NewMemStorage()
	service := service.NewMetricsService(storage)
	handler.PrepareTestData(service)

	server := handler.CreateTestServerWithStorage(storage)
	defer server.Close()

	tests := []struct {
		name       string
		body       string
		statusCode int
		response   string
	}{
		{"update unknown counter", `{"id": "UnknownCounter", "type": "counter", "delta": 100}`, http.StatusOK, `{"id": "UnknownCounter", "type": "counter", "delta": 100}`},
		{"update unknown gauge", `{"id": "UnknownGauge", "type": "gauge", "value": 3.14}`, http.StatusOK, `{"id": "UnknownGauge", "type": "gauge", "value": 3.14}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, body := handler.MakeTestRequest(t, server, http.MethodPost, "/update", tt.body)
			defer response.Body.Close()

			assert.Equal(t, tt.statusCode, response.StatusCode)
			if response.StatusCode == http.StatusOK {
				assert.JSONEq(t, tt.response, body)
			}
		})
	}
}

func TestUpdateMetricHandler_UpdateCounterAndCheckSum(t *testing.T) {
	server := handler.CreateTestServer()
	defer server.Close()

	t.Run("update counter and check sum", func(t *testing.T) {
		id := "TestID" + strconv.Itoa(rand.Intn(256))

		value1 := rand.Int63n(256)
		response1, _ := handler.MakeTestRequest(t, server, http.MethodPost, "/update", fmt.Sprintf(`{"id": "%s", "type": "counter", "delta": %d}`, id, value1))
		defer response1.Body.Close()

		value2 := rand.Int63n(256)
		response2, _ := handler.MakeTestRequest(t, server, http.MethodPost, "/update", fmt.Sprintf(`{"id": "%s", "type": "counter", "delta": %d}`, id, value2))
		defer response2.Body.Close()

		var resultMetric models.Metrics
		client := resty.New().SetTimeout(3 * time.Second)
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"id": "%s", "type": "counter"}`, id)).
			SetResult(&resultMetric).
			Post(server.URL + "/value")

		assert.Equal(t, http.StatusOK, resp.StatusCode())
		require.NoError(t, err)
		assert.Equal(t, value1+value2, *resultMetric.Delta)
	})
}
