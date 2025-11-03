package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	model "github.com/VladimirB/gometrics/internal/model"
)

const UpdateMetricPath = "/update/"

func UpdateMetricHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)

		if r.Method != http.MethodPost {
			http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, UpdateMetricPath)
		paths := strings.Split(path, "/")
		if len(paths) != 3 { // type/name/value
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}

		metricType := paths[0]
		metricValue := paths[2]
		switch metricType {
		case model.Gauge:
			if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
				http.Error(w, "incorrect metric value type, only float allowed", http.StatusBadRequest)
				return
			}
		case model.Counter:
			if _, err := strconv.Atoi(metricValue); err != nil {
				http.Error(w, "incorrect metric value type, only int allowed", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "incorrect metric type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
