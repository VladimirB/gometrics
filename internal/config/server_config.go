package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerAddress   = "localhost:8080"
	defaultStoreInterval   = 3
	defaultFileStoragePath = "metrics.json"
	defaultRestore         = true
)

type ServerConfig struct {
	Address         string        `env:"ADDRESS"`           // адрес сервера
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`    // интервал записи метрик на диск в секундах
	FileStoragePath string        `env:"FILE_STORAGE_PATH"` // путь до файла для записи метрик
	Restore         bool          `env:"RESTORE"`           // флаг, инициализировать ли значения метрик из файла на старте сервера
}

type serverFlags struct {
	address         string
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func GetServerConfig() ServerConfig {
	config := ServerConfig{}
	if err := env.Parse(&config); err != nil {
		log.Println(err)
	}

	flags := parseServerFlags()

	if config.Address == "" {
		config.Address = flags.address
	}

	if config.StoreInterval == 0 {
		config.StoreInterval = time.Duration(flags.storeInterval) * time.Second
	}

	if config.FileStoragePath == "" {
		config.FileStoragePath = flags.fileStoragePath
	}

	if _, exists := os.LookupEnv("RESTORE"); !exists {
		config.Restore = flags.restore
	}

	return config
}

func parseServerFlags() serverFlags {
	flags := serverFlags{}
	flag.StringVar(&flags.address, "a", defaultServerAddress, "server address")
	flag.IntVar(&flags.storeInterval, "i", defaultStoreInterval, "interval to save metrics into file")
	flag.StringVar(&flags.fileStoragePath, "f", defaultFileStoragePath, "file path to save metrics")
	flag.BoolVar(&flags.restore, "r", defaultRestore, "flag to restore metrics from file on server start")
	flag.Parse()
	return flags
}
