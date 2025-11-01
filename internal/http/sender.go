package http

import model "github.com/VladimirB/gometrics/internal/model"

type Sender interface {
	PostCounter(model.Metrics) (Response, error)
	PostGauge(model.Metrics) (Response, error)
}

type Response struct {
	StatusCode int
}