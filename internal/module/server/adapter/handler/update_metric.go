package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/module/server/port"
	"github.com/VladimirB/gometrics/internal/shared/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UpdateMetricHandler struct {
	metricService port.MetricService
}

func NewUpdateMetricHandler(service port.MetricService) *UpdateMetricHandler {
	return &UpdateMetricHandler{
		metricService: service,
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

		var metric domain.Metric
		if err = json.Unmarshal(buffer.Bytes(), &metric); err != nil {
			logger.Log.Info("error unmarshaling metric", zap.String("Metric", metric.String()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.metricService.Save(r.Context(), metric); err != nil {
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

		if err := h.metricService.SaveByFields(r.Context(), metricType, metricName, metricValue); err != nil {
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
