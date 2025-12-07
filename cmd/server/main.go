package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/handler"
	"github.com/VladimirB/gometrics/internal/logger"
	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/VladimirB/gometrics/internal/service"
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

	config := config.GetServerConfig()
	logger.Log.Info("Running Server", zap.String("Start time", time.Now().Local().String()), zap.Any("config", config))

	memStorage := repository.NewMemStorage()
	metricsService := service.NewMetricsService(memStorage)

	go runStoreInFileTicker(config, metricsService)

	mainPageTemplate := template.Must(template.ParseFS(templateFiles, "web/template/index.html"))
	mainPageHandler := handler.NewMainPageHandler(metricsService, mainPageTemplate)
	updateHandler := handler.NewUpdateMetricHandler(metricsService)
	valueHandler := handler.NewValueMetricHandler(metricsService)
	err := http.ListenAndServe(config.Address, handler.NewRouter(mainPageHandler, updateHandler, valueHandler))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func runStoreInFileTicker(config config.ServerConfig, service *service.MetricsService) {
	storeInFileTicker := time.NewTicker(config.StoreInterval)
	defer storeInFileTicker.Stop()

	for {
		select {
		case <-storeInFileTicker.C:
			logger.Log.Info("Write metrics to file")

			fileWriter, err := repository.NewJsonFileWriter(config.FileStoragePath)
			if err != nil {
				logger.Log.Fatal("Cant initialize metrics file writer", zap.Error(err))
			}
			defer fileWriter.Close()

			var metrics = make([]models.Metrics, 0)
			for _, value := range service.GetAll() {
				metrics = append(metrics, value)
			}

			err = fileWriter.Write(metrics)
			if err != nil {
				logger.Log.Error("Error on metrics write to file", zap.Error(err))
			}
		}
	}
}
