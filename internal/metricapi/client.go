package metricapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/VladimirB/gometrics/internal/compress"
	"github.com/VladimirB/gometrics/internal/logger"
	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
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

	// compressedBody, err := compress.Zip(body)
	// if err != nil {
	// 	return err
	// }

	response, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(body).
		Post(fmt.Sprintf("http://%s/update", server))
	if err != nil {
		return err
	}

	contentEncoding := response.Header().Get("Content-Encoding")
	gzipUsed := strings.Contains(contentEncoding, "gzip")
	if gzipUsed {
		compress.Unzip(response.Body())
	}
	
	logger.Log.Info("Response received", 
		zap.String("body", string(body)),
		zap.String("response", response.String()), 
		zap.Int("status code", response.StatusCode()), 
		zap.Bool("responsed with gzip", gzipUsed))

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("error on metric update: %d, %q, %v", response.StatusCode(), response.String(), metric)
	}

	return nil
}
