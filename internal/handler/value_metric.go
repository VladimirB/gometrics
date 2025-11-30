package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VladimirB/gometrics/internal/logger"
	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ValueMetricHandler struct {
	metricsService *service.MetricsService
}

func NewValueMetricHandler(service *service.MetricsService) *ValueMetricHandler {
	return &ValueMetricHandler{
		metricsService: service,
	}
}

func (h ValueMetricHandler) GetMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		if metricType != models.Counter && metricType != models.Gauge {
			http.Error(w, "unsupported metric type", http.StatusBadRequest)
			return
		}

		metricName := chi.URLParam(r, "metricName")
		if metric, err := h.metricsService.Get(metricName); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			if metric.MType == models.Counter {
				w.Write([]byte(fmt.Sprint(*metric.Delta)))
			} else {
				w.Write([]byte(fmt.Sprint(*metric.Value)))
			}
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (h ValueMetricHandler) PostValueMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buffer bytes.Buffer
		if _, err := buffer.ReadFrom(r.Body); err != nil {
			logger.Log.Error("cant read request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var askedMetric models.Metrics
		if err := json.Unmarshal(buffer.Bytes(), &askedMetric); err != nil {
			logger.Log.Error("cant unmarshal metric to JSON", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := askedMetric.Validate(); err != nil {
			logger.Log.Error("asked metric invalid", zap.Error(err), zap.String("metric", askedMetric.String()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metric, err := h.metricsService.Get(askedMetric.ID)
		if err != nil {
			metric = models.Metrics{
				ID: askedMetric.ID,
				MType: askedMetric.MType,
			}

			if askedMetric.MType == models.Counter {
				metric.Delta = new(int64)
			} else {
				metric.Value = new(float64)
			}
		}

		if resp, err := json.Marshal(metric); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(resp)
			w.WriteHeader(http.StatusOK)
		}
	}
}
