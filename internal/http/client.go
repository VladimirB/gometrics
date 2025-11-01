package http

import (
	"fmt"
	"net/http"
	"time"

	models "github.com/VladimirB/gometrics/internal/model"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c Client) PostCounter(metric models.Metrics) (Response, error) {
	url := fmt.Sprintf("http://localhost:8080/update/counter/%s/%d", metric.ID, int(*metric.Value))

	response, err := c.client.Post(url, "text/plain", http.NoBody)
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()

	return Response{
		StatusCode: response.StatusCode,
	}, nil
}

func (c Client) PostGauge(metric models.Metrics) (Response, error) {
	url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%f", metric.ID, *metric.Value)

	response, err := c.client.Post(url, "text/plain", http.NoBody)
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()

	return Response{
		StatusCode: response.StatusCode,
	}, nil
}