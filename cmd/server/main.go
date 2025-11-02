package main

import (
	"net/http"

	"github.com/VladimirB/gometrics/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(handler.UPDATE_METRIC_PATTERN, handler.UpdateMetricHandler())

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
