package main

import (
	"net/http"

	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/VladimirB/gometrics/internal/repository"
)

func main() {
	memStorage := repository.NewMemStorage()
	mainPageHandler := handler.NewMainPageHandler(memStorage)
	updateHandler := handler.NewUpdateMetricHandler(memStorage)
	valueHandler := handler.NewValueMetricHandler(memStorage)

	err := http.ListenAndServe(":8080", handler.NewRouter(mainPageHandler, updateHandler, valueHandler))
	if err != nil {
		panic(err)
	}
}
