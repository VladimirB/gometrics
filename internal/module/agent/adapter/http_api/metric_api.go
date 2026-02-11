package http_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/shared/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type MetricApi struct {
	server string
	client *resty.Client
}

func NewMetricApi(server string) *MetricApi {
	return &MetricApi{
		server: server,
		client: resty.New().SetTimeout(3 * time.Second),
	}
}

func (c MetricApi) Send(ctx context.Context, metric models.Metrics) error {
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
