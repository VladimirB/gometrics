package config

import (
	"flag"
	"os"

	"github.com/VladimirB/gometrics/internal/shared/logger"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

const (
	defaultServerAddress   = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "metrics.json"
	defaultRestore         = true
)

type StorageType string

const (
	StorageTypeMemory StorageType = "memory"
	StorageTypeDB     StorageType = "db"
)

type ServerConfig struct {
	Address         string `env:"ADDRESS"`           // адрес сервера
	StoreInterval   int    `env:"STORE_INTERVAL"`    // интервал записи метрик на диск в секундах
	FileStoragePath string `env:"FILE_STORAGE_PATH"` // путь до файла для записи метрик
	Restore         bool   `env:"RESTORE"`           // флаг, инициализировать ли значения метрик из файла на старте сервера
	DatabaseDSN     string `env:"DATABASE_DSN"`      // строка для подключения к базе данных
}

type serverFlags struct {
	address         string
	storeInterval   int
	fileStoragePath string
	restore         bool
	databaseDSN     string
}

func GetServerConfig() ServerConfig {
	config := ServerConfig{}
	if err := env.Parse(&config); err != nil {
		logger.Log.Error("Cant parse envs", zap.Error(err))
	}

	flags := parseServerFlags()

	if config.Address == "" {
		config.Address = flags.address
	}

	if config.StoreInterval == 0 {
		config.StoreInterval = flags.storeInterval
	}

	if config.FileStoragePath == "" {
		config.FileStoragePath = flags.fileStoragePath
	}

	if _, exists := os.LookupEnv("RESTORE"); !exists {
		config.Restore = flags.restore
	}

	if config.DatabaseDSN == "" {
		config.DatabaseDSN = flags.databaseDSN
	}

	return config
}

func parseServerFlags() serverFlags {
	flags := serverFlags{}
	flag.StringVar(&flags.address, "a", defaultServerAddress, "server address")
	flag.IntVar(&flags.storeInterval, "i", defaultStoreInterval, "interval to save metrics into file")
	flag.StringVar(&flags.fileStoragePath, "f", defaultFileStoragePath, "file path to save metrics")
	flag.BoolVar(&flags.restore, "r", defaultRestore, "flag to restore metrics from file on server start")
	flag.StringVar(&flags.databaseDSN, "d", "", "database connection string")
	flag.Parse()
	return flags
}

func (c ServerConfig) MetricStorageType() StorageType {
	if c.DatabaseDSN != "" {
		return StorageTypeDB
	}

	return StorageTypeMemory
}
