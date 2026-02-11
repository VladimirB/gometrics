package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/module/server/port"
	"github.com/VladimirB/gometrics/internal/shared/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ValueMetricHandler struct {
	metricService port.MetricService
}

func NewValueMetricHandler(service port.MetricService) *ValueMetricHandler {
	return &ValueMetricHandler{
		metricService: service,
	}
}

func (h ValueMetricHandler) GetMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		if metricType != domain.Counter && metricType != domain.Gauge {
			http.Error(w, "unsupported metric type", http.StatusBadRequest)
			return
		}

		metricName := chi.URLParam(r, "metricName")
		metric, err := h.metricService.Get(r.Context(), metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		switch metric.MType {
		case domain.Counter:
			w.Write([]byte(fmt.Sprint(*metric.Delta)))
		case domain.Gauge:
			w.Write([]byte(fmt.Sprint(*metric.Value)))
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h ValueMetricHandler) PostValueMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buffer []byte
		buffer, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Info("cant read request body", zap.Error(err))
			fmt.Println("Cant read request body", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var askedMetric domain.Metric
		if err := json.Unmarshal(buffer, &askedMetric); err != nil {
			logger.Log.Error("cant unmarshal metric to JSON", zap.Error(err))
			fmt.Println("Cant unmarshal metric to JSON", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := askedMetric.Validate(); err != nil {
			logger.Log.Info("asked metric invalid", zap.Error(err), zap.String("metric", askedMetric.String()))
			fmt.Println("Validation Error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metric, err := h.metricService.Get(r.Context(), askedMetric.ID)
		if err != nil {
			metric = domain.Metric{
				ID:    askedMetric.ID,
				MType: askedMetric.MType,
			}

			if askedMetric.MType == domain.Counter {
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
