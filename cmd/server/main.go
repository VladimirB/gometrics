package main

import (
	"net/http"

	"github.com/VladimirB/gometrics/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(handler.UpdateMetricPath, handler.UpdateMetricHandler())

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
