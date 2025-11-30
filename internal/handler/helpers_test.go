package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/VladimirB/gometrics/internal/service"
	"github.com/stretchr/testify/require"
)

func CreateTestServer() *httptest.Server {
	memStorage := repository.NewMemStorage()
	metricsService := service.NewMetricsService(memStorage)
	mainPageHandler := NewMainPageHandler(memStorage, nil)
	updateHandler := NewUpdateMetricHandler(metricsService)
	valueHandler := NewValueMetricHandler(metricsService)
	return httptest.NewServer(NewRouter(mainPageHandler, updateHandler, valueHandler))
}

func CreateTestServerWithStorage(storage *repository.MemStorage) *httptest.Server {
	memStorage := storage
	metricsService := service.NewMetricsService(memStorage)
	mainPageHandler := NewMainPageHandler(memStorage, nil)
	updateHandler := NewUpdateMetricHandler(metricsService)
	valueHandler := NewValueMetricHandler(metricsService)
	return httptest.NewServer(NewRouter(mainPageHandler, updateHandler, valueHandler))
}

func MakeTestRequest(t *testing.T, ts *httptest.Server, method string, path string, body string) (*http.Response, string) {
	var br io.Reader
	if body != "" {
		br = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, ts.URL+path, br)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
