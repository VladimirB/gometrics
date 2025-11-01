package metricapi

import (
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

func (c Client) PostMetric(url string, metric models.Metrics) (Response, error) {
	response, err := c.client.Post(url, "text/plain", http.NoBody)
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()

	return Response{
		StatusCode: response.StatusCode,
	}, nil
}
