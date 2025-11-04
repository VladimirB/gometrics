package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/stretchr/testify/require"
)

func CreateTestServer() *httptest.Server {
	memStorage := repository.NewMemStorage()
	updateHandler := NewUpdateMetricHandler(memStorage)
	valueHandler := NewValueMetricHandler(memStorage)
	return httptest.NewServer(NewRouter(updateHandler, valueHandler))
}

func CreateTestServerWithStorage(storage *repository.MemStorage) *httptest.Server {
	memStorage := storage
	updateHandler := NewUpdateMetricHandler(memStorage)
	valueHandler := NewValueMetricHandler(memStorage)
	return httptest.NewServer(NewRouter(updateHandler, valueHandler))
}

func MakeTestRequest(t *testing.T, ts *httptest.Server, method string, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
