package handler

import (
	"net/http"
	"strconv"

	model "github.com/VladimirB/gometrics/internal/model"
	"github.com/go-chi/chi/v5"
)

const UpdateMetricPath = "/update/"

func PostMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		if metricType != model.Gauge && metricType != model.Counter {
			http.Error(w, "incorrect metric type", http.StatusBadRequest)
			return
		}

		if _, err := strconv.ParseFloat(chi.URLParam(r, "metricValue"), 64); err != nil {
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func PostMetricNoNameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "incorrect path", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
	}
}
