package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/shared/logger"
	"github.com/go-resty/resty/v2"
)

type MetricAPI struct {
	server string
	client *resty.Client
}

func NewMetricAPI(server string) *MetricAPI {
	return &MetricAPI{
		server: server,
		client: resty.New().
			SetTimeout(3 * time.Second).
			SetLogger(logger.Log.Sugar()),
	}
}

func (m *MetricAPI) Send(ctx context.Context, metric domain.Metric) error {
	request := mapToMetricDTO(metric)
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := m.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetContext(ctx).
		Post(fmt.Sprintf("http://%s/update", m.server))
	if err != nil {
		return fmt.Errorf("POST /update failed: %w", err)
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("error on metric update: %d, %q, %v", response.StatusCode(), response.String(), metric)
	}

	return nil
}

func (m *MetricAPI) SendAll(ctx context.Context, metrics []domain.Metric) error {
	dtos := make([]metricDTO, len(metrics))
	for i, m := range metrics {
		dtos[i] = mapToMetricDTO(m)
	}

	body, err := json.Marshal(dtos)
	if err != nil {
		return fmt.Errorf("marshaling to JSON failed: %w", err)
	}

	response, err := m.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetContext(ctx).
		Post(fmt.Sprintf("http://%s/updates", m.server))
	if err != nil {
		return fmt.Errorf("POST /updates failed: %w", err)
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("error on metrics update: %d, %q", response.StatusCode(), response.String())
	}

	logger.Log.Info("")

	return nil
}
