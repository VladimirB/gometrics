package main

import (
	"database/sql"
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/module/server/adapter/handler"
	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
	"github.com/VladimirB/gometrics/internal/module/server/service"
	"github.com/VladimirB/gometrics/internal/shared/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	db, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		logger.Log.Error("DB connection is not well", zap.Error(err))
	}
	defer db.Close()

	memStorage := repository.NewMemStorage()
	metricsService := service.NewMetricsService(memStorage)

	// Если в конфигурации установлен флаг Restore, необходимо прочитать значения метрик из файла (путь в конфиге)
	// и инициализировать хранилище метрик прочитанными значениями
	if config.Restore {
		initMetricsFromFile(config.FileStoragePath, metricsService)
	}

	go runStoreInFileTicker(config, metricsService)

	mainPageTemplate := template.Must(template.ParseFS(templateFiles, "web/template/index.html"))
	mainPageHandler := handler.NewMainPageHandler(metricsService, mainPageTemplate)
	updateHandler := handler.NewUpdateMetricHandler(metricsService)
	valueHandler := handler.NewValueMetricHandler(metricsService)
	dbPingHandler := handler.NewDatabasePingHandler(db)
	err = http.ListenAndServe(config.Address, handler.NewRouter(mainPageHandler, updateHandler, valueHandler, dbPingHandler))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func initMetricsFromFile(fileName string, metricsService *service.MetricsService) {
	jsonFileReader, err := repository.NewMetricsFileReader(fileName)
	if err != nil {
		logger.Log.Error("Cant open file to read metrics", zap.Error(err))
		return
	}

	result, err := jsonFileReader.Read()
	if err != nil {
		logger.Log.Error("Cant read metrics from file", zap.Error(err))
		return
	}

	for _, m := range result {
		err := metricsService.Save(m)
		if err != nil {
			logger.Log.Error("Error on save metric", zap.Error(err), zap.String("Metric", m.String()))
		}
	}
}

func runStoreInFileTicker(config config.ServerConfig, service *service.MetricsService) {
	storeInFileTicker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
	defer storeInFileTicker.Stop()

	for tick := range storeInFileTicker.C {
		logger.Log.Info("Write metrics to file", zap.Any("Seconds from start", tick.Second()))

		currentMetrics := service.GetAll()
		if len(currentMetrics) == 0 {
			continue
		}

		fileWriter, err := repository.NewMetricsFileWriter(config.FileStoragePath)
		if err != nil {
			logger.Log.Fatal("Cant initialize metrics file writer", zap.Error(err))
			continue
		}

		var metrics []models.Metrics
		for _, value := range currentMetrics {
			metrics = append(metrics, value)
		}

		err = fileWriter.Write(metrics)
		if err != nil {
			logger.Log.Error("Error on metrics write to file", zap.Error(err))
		}
		fileWriter.Close()
	}
}
