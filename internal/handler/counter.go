package handler

import (
	"net/http"
	"strconv"
	"strings"
)

func CounterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		paths = paths[3:] // обрезаем часть пути /update/counter

		if len(paths) != 2 { // metricname/value
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}

		metricValue := paths[1]
		if _, err := strconv.Atoi(metricValue); err != nil {
			http.Error(w, "incorrect metric value type, only int allowed", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}