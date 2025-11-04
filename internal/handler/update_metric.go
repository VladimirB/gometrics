package handler

import (
	"net/http"
	"strconv"

	model "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

type UpdateMetricHandler struct {
	storage *repository.MemStorage
}

func NewUpdateMetricHandler(storage *repository.MemStorage) *UpdateMetricHandler {
	return &UpdateMetricHandler{
		storage: storage,
	}
}

func (h UpdateMetricHandler) PostMetricHandler() http.HandlerFunc {
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

		metric := model.Metrics {
			ID: metricName,
			MType: metricType,
			Value: new(float64),
		}
		metric.Value = &metricValue
		h.storage.Save(metricName, metric)

		w.WriteHeader(http.StatusOK)
	}
}

func (UpdateMetricHandler) PostMetricNoNameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "incorrect path", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
	}
}
