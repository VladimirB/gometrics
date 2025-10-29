package main

import (
	"net/http"
	"github.com/VladimirB/gometrics/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/update/gauge/", handler.GaugeHandler())
	mux.Handle("/update/counter/", handler.CounterHandler())

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
