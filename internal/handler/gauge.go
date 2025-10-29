package handler

import (
	"net/http"
	"strconv"
	"strings"
)

func GaugeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		paths = paths[3:] // обрезаем часть пути /update/gauge

		if len(paths) != 2 { // metricname/value
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}

		metricValue := paths[1]
		if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
			http.Error(w, "incorrect metric value type, only float allowed", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}