package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/shared/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type MetricAPI struct {
	server string
	client *resty.Client
}

func NewMetricAPI(server string) *MetricAPI {
	return &MetricAPI{
		server: server,
		client: resty.New().SetTimeout(3 * time.Second),
	}
}

func (c MetricAPI) Send(ctx context.Context, metric domain.Metric) error {
	request := mapToMetricRequest(metric)
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetContext(ctx).
		Post(fmt.Sprintf("http://%s/update", c.server))
	if err != nil {
		return err
	}

	contentEncoding := response.Header().Get("Content-Encoding")
	gzipUsed := strings.Contains(contentEncoding, "gzip")
	logger.Log.Info("Response received",
		zap.Int("status code", response.StatusCode()),
		zap.Bool("responsed with gzip", gzipUsed))

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("error on metric update: %d, %q, %v", response.StatusCode(), response.String(), metric)
	}

	return nil
}
