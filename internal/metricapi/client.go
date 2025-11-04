package metricapi

import (
	"fmt"
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

func (c Client) PostMetric(metric models.Metrics) (Response, error) {
	response, err := c.client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"metricType": metric.MType,
			"metricName": metric.ID,
			"metricValue": fmt.Sprintf("%f", *metric.Value),
		}).
		Post("http://localhost:8080/update/{metricType}/{metricName}/{metricValue}")

	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: response.StatusCode(),
		Body: response.String(),
	}, nil
}
