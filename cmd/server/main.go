package main

import (
	"fmt"
	"net/http"

	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/VladimirB/gometrics/internal/repository"
)

func main() {
	memStorage := repository.NewMemStorage()
	mainPageHandler := handler.NewMainPageHandler(memStorage)
	updateHandler := handler.NewUpdateMetricHandler(memStorage)
	valueHandler := handler.NewValueMetricHandler(memStorage)

	parseFlags()
	fmt.Printf("Run server with %s\n", serverConfig.Address.String())

	err := http.ListenAndServe(serverConfig.Address.String(), handler.NewRouter(mainPageHandler, updateHandler, valueHandler))
	if err != nil {
		panic(err)
	}
}
