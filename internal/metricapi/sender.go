package metricapi

import model "github.com/VladimirB/gometrics/internal/model"

type Sender interface {
	PostMetric(server string, metric model.Metrics) (Response, error)
}

type Response struct {
	StatusCode int
	Body string
}