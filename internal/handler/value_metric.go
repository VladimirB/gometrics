package handler

import (
	"fmt"
	"net/http"

	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

type ValueMetricHandler struct {
	storage *repository.MemStorage
}

func NewValueMetricHandler(storage *repository.MemStorage) *ValueMetricHandler {
	return &ValueMetricHandler{
		storage: storage,
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
		if metric, err := h.storage.Get(metricName); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			w.Write([]byte(fmt.Sprint(*metric.Value)))
			w.WriteHeader(http.StatusOK)
		}
	}
}
