package main

import (
	"log"
	"net/http"

	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/VladimirB/gometrics/internal/repository"
)

func main() {
	memStorage := repository.NewMemStorage()
	mainPageHandler := handler.NewMainPageHandler(memStorage)
	updateHandler := handler.NewUpdateMetricHandler(memStorage)
	valueHandler := handler.NewValueMetricHandler(memStorage)

	config := config.NewServerConfig()
	parseFlags(config)
	log.Println("run server:", config.Address.String())

	err := http.ListenAndServe(config.Address.String(), handler.NewRouter(mainPageHandler, updateHandler, valueHandler))
	if err != nil {
		log.Fatal(err)
	}
}
