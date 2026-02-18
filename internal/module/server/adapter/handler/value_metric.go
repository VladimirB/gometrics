package handler

import (
	"encoding/json"
	"fmt"
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
		var req metricRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Log.Error("cant decode request from JSON", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metric, err := h.metricService.Get(r.Context(), req.ID)
		if err != nil { // в случае ошибки по требования нужно возвращать пустое значение метрики
			logger.Log.Error("metric /value error", zap.Error(err))

			metric = domain.Metric{
				ID:    req.ID,
				MType: req.MType,
			}

			if req.MType == domain.Counter {
				metric.Delta = new(int64)
			} else {
				metric.Value = new(float64)
			}
		}

		response := mapToMetricResponse(metric)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
