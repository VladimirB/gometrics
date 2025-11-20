package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/VladimirB/gometrics/internal/logger"
	"github.com/VladimirB/gometrics/internal/repository"
	"go.uber.org/zap"
)

const logLevel = "info"

//go:embed web/template
var templateFiles embed.FS

func main() {
	if err := logger.Initialize(logLevel); err != nil {
		log.Fatal(err)
	}
	defer logger.Log.Sync()

	mainPageTemplate := template.Must(template.ParseFS(templateFiles, "web/template/index.html"))

	memStorage := repository.NewMemStorage()
	mainPageHandler := handler.NewMainPageHandler(memStorage, mainPageTemplate)
	updateHandler := handler.NewUpdateMetricHandler(memStorage)
	valueHandler := handler.NewValueMetricHandler(memStorage)

	config := config.GetServerConfig()
	logger.Log.Info("run server", zap.String("address", config.Address))

	err := http.ListenAndServe(config.Address, handler.NewRouter(mainPageHandler, updateHandler, valueHandler))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
