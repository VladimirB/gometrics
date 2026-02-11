package repository

import (
	"encoding/json"
	"os"

	"github.com/VladimirB/gometrics/internal/domain"
)

type MetricsFileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewMetricsFileWriter(fileName string) (*MetricsFileWriter, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &MetricsFileWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (fw *MetricsFileWriter) Write(metrics []domain.Metric) error {
	return fw.encoder.Encode(metrics)
}

func (fw *MetricsFileWriter) Close() error {
	return fw.file.Close()
}
