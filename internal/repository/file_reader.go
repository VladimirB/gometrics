package repository

import (
	"encoding/json"
	"os"

	"github.com/VladimirB/gometrics/internal/domain"
)

type MetricsFileReader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewMetricsFileReader(fileName string) (*MetricsFileReader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &MetricsFileReader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (fr *MetricsFileReader) Read() ([]domain.Metric, error) {
	var result = make([]domain.Metric, 0)

	if err := fr.decoder.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (fr *MetricsFileReader) Close() error {
	return fr.file.Close()
}
