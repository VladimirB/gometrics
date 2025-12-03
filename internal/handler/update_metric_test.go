package handler_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/VladimirB/gometrics/internal/handler"
	models "github.com/VladimirB/gometrics/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostMetricHandler(t *testing.T) {
	testTable := []struct {
		name       string
		path       string
		statusCode int
	}{
		{"success counter metric", "/update/counter/PollCount/100", http.StatusOK},
		{"success gauge metric", "/update/gauge/Alloc/100.55", http.StatusOK},
		{"incorrect counter metric value", "/update/counter/name/str", http.StatusBadRequest},
		{"incorrect gauge metric value", "/update/gauge/name/str", http.StatusBadRequest},
		{"unknown metric type", "/update/unknown/name/100", http.StatusBadRequest},
		{"no counter metric name", "/update/counter/100", http.StatusNotFound},
		{"no gauge metric name", "/update/gauge/100", http.StatusNotFound},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			response, err := httpClient.R().Post(testServer.URL + tt.path)

			require.NoError(t, err)
			assert.Equal(t, tt.statusCode, response.StatusCode())
		})
	}
}

func TestUpdateMetricHandler_UpdateMetricJSONHandler(t *testing.T) {
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
			response, err := httpClient.R().
				SetBody(tt.body).
				Post(testServer.URL + "/update")

			require.NoError(t, err)
			assert.Equal(t, tt.statusCode, response.StatusCode())
			if response.StatusCode() == http.StatusOK {
				assert.JSONEq(t, tt.response, string(response.Body()))
			}
		})
	}
}

func TestUpdateMetricHandler_UpdateCounterAndCheckSum(t *testing.T) {
	t.Run("update counter and check sum", func(t *testing.T) {
		id := handler.GenTestID()

		value1 := rand.Int63n(256)
		_, err1 := httpClient.R().
			SetBody(fmt.Sprintf(`{"id": "%s", "type": "counter", "delta": %d}`, id, value1)).
			Post(testServer.URL + "/update")
		require.NoError(t, err1)

		value2 := rand.Int63n(256)
		_, err2 := httpClient.R().
			SetBody(fmt.Sprintf(`{"id": "%s", "type": "counter", "delta": %d}`, id, value2)).
			Post(testServer.URL + "/update")
		require.NoError(t, err2)

		var resultMetric models.Metrics
		resp, err := httpClient.R().
			SetBody(fmt.Sprintf(`{"id": "%s", "type": "counter"}`, id)).
			SetResult(&resultMetric).
			Post(testServer.URL + "/value")

		assert.Equal(t, http.StatusOK, resp.StatusCode())
		require.NoError(t, err)
		assert.Equal(t, value1+value2, *resultMetric.Delta)
	})
}
