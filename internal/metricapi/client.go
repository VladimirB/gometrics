package metricapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewClient() *Client {
	return &Client{
		client: resty.New().SetTimeout(3 * time.Second),
	}
}

func (c Client) Send(server string, metric models.Metrics) error {
	response, err := c.client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"metricType":  metric.MType,
			"metricName":  metric.ID,
			"metricValue": fmt.Sprintf("%f", *metric.Value),
		}).
		Post(fmt.Sprintf("http://%s/update/{metricType}/{metricName}/{metricValue}", server))

	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("error on metric send: %d, %q, %v", response.StatusCode(), response.String(), metric)
	}

	return nil
}

func (c Client) SendAsJSON(server string, metric models.Metrics) error {
	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	response, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://%s/update", server))
	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("error on metric update: %d, %q, %v", response.StatusCode(), response.String(), metric)
	}

	return nil
}
