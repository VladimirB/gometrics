package main

import (
	"net/http"
	"github.com/VladimirB/gometrics/internal/handler"
)

func main() {
	http.HandleFunc("/update/", handler.MetricsHandler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
