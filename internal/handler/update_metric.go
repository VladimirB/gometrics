package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VladimirB/gometrics/internal/logger"
	model "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UpdateMetricHandler struct {
	storage *repository.MemStorage
}

func NewUpdateMetricHandler(storage *repository.MemStorage) *UpdateMetricHandler {
	return &UpdateMetricHandler{
		storage: storage,
	}
}

func (h UpdateMetricHandler) UpdateMetricJsonHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buffer bytes.Buffer
		_, err := buffer.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var metric model.Metrics
		if err = json.Unmarshal(buffer.Bytes(), &metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if metric.ID == "" {
			http.Error(w, "incorrect metric ID", http.StatusBadRequest)
			return
		}

		if metric.MType != model.Gauge && metric.MType != model.Counter {
			http.Error(w, "incorrect metric type", http.StatusBadRequest)
			return
		}

		logger.Log.Info("Received metric", zap.String("id", metric.ID), zap.String("type", metric.MType), zap.Float64("value", *metric.Value))
		h.storage.Save(metric.ID, metric)

		w.WriteHeader(http.StatusOK)
	}
}

func (h UpdateMetricHandler) UpdateMetricByNameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		if metricType != model.Gauge && metricType != model.Counter {
			http.Error(w, "incorrect metric type", http.StatusBadRequest)
			return
		}

		metricName := chi.URLParam(r, "metricName")
		if metricName == "" {
			http.Error(w, "incorrect metric name", http.StatusBadRequest)
			return
		}

		var metricValue float64
		if parsed, err := strconv.ParseFloat(chi.URLParam(r, "metricValue"), 64); err != nil {
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
			return
		} else {
			metricValue = parsed
		}

		metric := model.Metrics{
			ID:    metricName,
			MType: metricType,
			Value: new(float64),
		}
		metric.Value = &metricValue
		h.storage.Save(metricName, metric)

		w.WriteHeader(http.StatusOK)
	}
}

func (UpdateMetricHandler) UpdateNoNameMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "incorrect path", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
	}
}
