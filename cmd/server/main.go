package main

import (
	"context"
	"database/sql"
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/module/server/adapter/database"
	"github.com/VladimirB/gometrics/internal/module/server/adapter/file"
	"github.com/VladimirB/gometrics/internal/module/server/adapter/handler"
	"github.com/VladimirB/gometrics/internal/module/server/adapter/memory"
	"github.com/VladimirB/gometrics/internal/module/server/port"
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

	serverConfig := config.GetServerConfig()
	logger.Log.Info("Running Server", zap.String("Start time", time.Now().Local().String()))

	ctx := context.Background()

	db, err := sql.Open("pgx", serverConfig.DatabaseDSN)
	if err != nil {
		logger.Log.Error("DB connection is not well", zap.Error(err))
	}
	defer db.Close()

	var metricRepo port.MetricRepository
	switch serverConfig.MetricStorageType() {
	case config.StorageTypeDB:
		metricRepo = database.NewPostgresRepo(db)
		logger.Log.Info("Database storage will be used")
	default:
		metricRepo = memory.NewMemStorage()
		logger.Log.Info("Memory storage will be used")
	}

	fileStorage := file.NewFileStorage()

	metricsService := service.NewMetricsService(metricRepo, fileStorage)

	// Если в конфигурации установлен флаг Restore, необходимо прочитать значения метрик из файла (путь в конфиге)
	// и инициализировать хранилище метрик прочитанными значениями
	if serverConfig.Restore {
		metricsService.RestoreFromFile(ctx, serverConfig.FileStoragePath)
	}

	if serverConfig.MetricStorageType() == config.StorageTypeMemory {
		go runStoreInFileTicker(ctx, serverConfig, metricsService)
	}

	mainPageTemplate := template.Must(template.ParseFS(templateFiles, "web/template/index.html"))
	mainPageHandler := handler.NewMainPageHandler(metricsService, mainPageTemplate)
	updateHandler := handler.NewUpdateMetricHandler(metricsService)
	valueHandler := handler.NewValueMetricHandler(metricsService)
	dbPingHandler := handler.NewDatabasePingHandler(db)
	err = http.ListenAndServe(serverConfig.Address, handler.NewRouter(mainPageHandler, updateHandler, valueHandler, dbPingHandler))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func runStoreInFileTicker(ctx context.Context, config config.ServerConfig, metricService port.MetricService) {
	storeInFileTicker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
	defer storeInFileTicker.Stop()

	for tick := range storeInFileTicker.C {
		logger.Log.Info("Write metrics to file", zap.Any("Seconds from start", tick.Second()))

		err := metricService.DumpToFile(ctx, config.FileStoragePath)
		if err != nil {
			logger.Log.Error("error on dump metrics to file", zap.Error(err))
		}
	}
}
