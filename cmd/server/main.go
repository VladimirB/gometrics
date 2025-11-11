package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/VladimirB/gometrics/internal/repository"
)

//go:embed web/template
var templateFiles embed.FS

func main() {
	mainPageTemplate := template.Must(template.ParseFS(templateFiles, "web/template/index.html"))

	memStorage := repository.NewMemStorage()
	mainPageHandler := handler.NewMainPageHandler(memStorage, mainPageTemplate)
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
