package main

import (
	"net/http"

	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/VladimirB/gometrics/internal/repository"
)

func main() {
	memStorage := repository.NewMemStorage()
	metricHandler := handler.NewUpdateMetricHandler(memStorage)
	
	err := http.ListenAndServe(":8080", handler.NewRouter(metricHandler))
	if err != nil {
		panic(err)
	}
}
