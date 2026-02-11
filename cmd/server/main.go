package main

import (
	"context"
	"database/sql"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/module/server/adapter/handler"
	"github.com/VladimirB/gometrics/internal/module/server/port"
	"github.com/VladimirB/gometrics/internal/module/server/service"
	"github.com/VladimirB/gometrics/internal/repository"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		initMetricsFromFile(ctx, config.FileStoragePath, metricsService)
	}

	go runStoreInFileTicker(ctx, config, metricsService)

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

func initMetricsFromFile(ctx context.Context, fileName string, metricService port.MetricService) {
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
		err := metricService.Save(ctx, m)
		if err != nil {
			logger.Log.Error("Error on save metric", zap.Error(err), zap.String("Metric", m.String()))
		}
	}
}

func runStoreInFileTicker(ctx context.Context, config config.ServerConfig, metricService port.MetricService) {
	storeInFileTicker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
	defer storeInFileTicker.Stop()

	for tick := range storeInFileTicker.C {
		logger.Log.Info("Write metrics to file", zap.Any("Seconds from start", tick.Second()))

		currentMetrics := metricService.GetAll(ctx)
		if len(currentMetrics) == 0 {
			continue
		}

		fileWriter, err := repository.NewMetricsFileWriter(config.FileStoragePath)
		if err != nil {
			logger.Log.Fatal("Cant initialize metrics file writer", zap.Error(err))
			continue
		}

		var metrics []domain.Metric
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
