package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/VladimirB/gometrics/internal/shared/compress"
	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/module/server/adapter/handler"
	"github.com/VladimirB/gometrics/internal/module/server/port"
	"github.com/VladimirB/gometrics/internal/module/server/service"
	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var httpClient *resty.Client
var metricsService port.MetricService
var testServer *httptest.Server

func setup() {
	fmt.Println("Setup tests")

	storage := repository.NewMemStorage()
	metricsService = service.NewMetricsService(storage)

	testServer = handler.CreateTestServer(metricsService)

	httpClient = resty.New().
		SetTimeout(3*time.Second).
		SetBaseURL(testServer.URL).
		SetHeader("Content-Type", "application/json")
}

func tearDown() {
	fmt.Println("Teardown tests")
	defer testServer.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func TestValueMetricHandler_GetMetricHandler(t *testing.T) {
	counterID := handler.GenTestID()
	metricsService.SaveByFields(context.Background(), models.Counter, counterID, 100)
	gaugeID := handler.GenTestID()
	metricsService.SaveByFields(context.Background(), models.Gauge, gaugeID, 3.14)

	testTable := []struct {
		name       string
		path       string
		statusCode int
		body       string
	}{
		{"success counter metric get", "/value/counter/" + counterID, http.StatusOK, "100"},
		{"success gauge metric get", "/value/gauge/" + gaugeID, http.StatusOK, "3.14"},
		{"metric not found", "/value/counter/Unknown", http.StatusNotFound, ""},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := httpClient.R().Get(tt.path)

			assert.Equal(t, tt.statusCode, response.StatusCode())
			if response.StatusCode() == http.StatusOK {
				assert.Equal(t, tt.body, string(response.Body()))
			}
		})
	}
}

func TestValueMetricHandler_PostValueMetricHandler(t *testing.T) {
	counterID := handler.GenTestID()
	metricsService.SaveByFields(context.Background(), models.Counter, counterID, 1000)
	gaugeID := handler.GenTestID()
	metricsService.SaveByFields(context.Background(), models.Gauge, gaugeID, 3.1415)

	tests := []struct {
		name       string
		body       string
		statusCode int
		response   string
	}{
		{"ask no value counter", `{"id":"AnyID", "type":"counter"}`, http.StatusOK, `{"id":"AnyID", "type":"counter", "delta":0}`},
		{"ask no value gauge", `{"id":"AnyID", "type":"gauge"}`, http.StatusOK, `{"id":"AnyID", "type":"gauge", "value":0}`},
		{"ask existed counter", fmt.Sprintf(`{"id":"%s", "type":"counter"}`, counterID), http.StatusOK, fmt.Sprintf(`{"id":"%s", "type":"counter", "delta":%d}`, counterID, 1000)},
		{"ask existed gauge", fmt.Sprintf(`{"id":"%s", "type":"gauge"}`, gaugeID), http.StatusOK, fmt.Sprintf(`{"id":"%s", "type":"gauge", "value":%f}`, gaugeID, 3.1415)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := httpClient.R().
				SetBody(tt.body).
				Post("/value")

			require.NoError(t, err)
			assert.Equal(t, tt.statusCode, response.StatusCode())
			if response.StatusCode() == http.StatusOK {
				assert.JSONEq(t, tt.response, string(response.Body()))
			}
		})
	}
}

func TestValueMetricHandler_PostValueMetricHandlerGzip(t *testing.T) {
	counterID := handler.GenTestID()
	metricsService.SaveByFields(context.Background(), models.Counter, counterID, 1000)

	requestBody := fmt.Sprintf(`{"id":"%s", "type":"counter"}`, counterID)
	expectedResponse := fmt.Sprintf(`{"id":"%s", "type":"counter", "delta":%d}`, counterID, 1000)

	t.Run("send gzip", func(t *testing.T) {
		compressedRequest, err := compress.Zip([]byte(requestBody))
		require.NoError(t, err)

		response, err := httpClient.R().
			SetHeader("Content-Encoding", "gzip").
			SetBody(compressedRequest).
			Post("/value")
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.StatusCode())
		if response.StatusCode() == http.StatusOK {
			assert.JSONEq(t, expectedResponse, string(response.Body()))
		}
	})

	t.Run("accept gzip", func(t *testing.T) {
		response, err := httpClient.R().
			SetHeader("Accept-Encoding", "gzip").
			SetBody(requestBody).
			Post("/value")
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.StatusCode())
		assert.Equal(t, "gzip", response.Header().Get("Content-Encoding"))
		if response.StatusCode() == http.StatusOK {
			assert.JSONEq(t, expectedResponse, string(response.Body()))

		}
	})
}
