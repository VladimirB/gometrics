package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VladimirB/gometrics/internal/shared/logger"
	model "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UpdateMetricHandler struct {
	metricsService *service.MetricsService
}

func NewUpdateMetricHandler(service *service.MetricsService) *UpdateMetricHandler {
	return &UpdateMetricHandler{
		metricsService: service,
	}
}

func (h UpdateMetricHandler) UpdateMetricJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buffer bytes.Buffer
		_, err := buffer.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var metric model.Metrics
		if err = json.Unmarshal(buffer.Bytes(), &metric); err != nil {
			logger.Log.Info("error unmarshaling metric", zap.String("Metric", metric.String()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.metricsService.Save(metric); err != nil {
			http.Error(w, "incorrect metric ID", http.StatusBadRequest)
			return
		}

		resp, err := json.Marshal(metric)
		if err != nil {
			logger.Log.Info("error marshaling metric", zap.String("Metric", metric.String()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		w.WriteHeader(http.StatusOK)
	}
}

func (h UpdateMetricHandler) UpdateMetricByNameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		var metricValue float64
		if parsed, err := strconv.ParseFloat(chi.URLParam(r, "metricValue"), 64); err != nil {
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
			return
		} else {
			metricValue = parsed
		}

		if err := h.metricsService.SaveByFields(metricType, metricName, metricValue); err != nil {
			http.Error(w, "incorrect metric ID", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (UpdateMetricHandler) UpdateNoNameMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "incorrect path", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
	}
}
