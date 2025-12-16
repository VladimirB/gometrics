package handler

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladimirB/gometrics/internal/service"
	"github.com/stretchr/testify/require"
)

func CreateTestServer(service *service.MetricsService) *httptest.Server {
	db, _, err := sqlmock.New()
    if err != nil {
        log.Fatalf("Error on sqlmock creation: %s", err)
    }
    defer db.Close()

	metricsService := service
	mainPageHandler := NewMainPageHandler(metricsService, nil)
	updateHandler := NewUpdateMetricHandler(metricsService)
	valueHandler := NewValueMetricHandler(metricsService)
	dbPingHandler := NewDatabasePingHandler(db)
	return httptest.NewServer(NewRouter(mainPageHandler, updateHandler, valueHandler, dbPingHandler))
}

func MakeTestRequest(t *testing.T, ts *httptest.Server, method string, path string, body string) (*http.Response, string) {
	var br io.Reader
	if body != "" {
		br = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, ts.URL+path, br)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func GenTestID() string {
	return fmt.Sprintf("TestID%d", rand.Intn(1024))
}
