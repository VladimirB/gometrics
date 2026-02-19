package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/VladimirB/gometrics/internal/domain"
)

type FileStorage struct{}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (fs *FileStorage) Read(ctx context.Context, fileName string) ([]domain.Metric, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("read file aborted: %w", err)
	}

	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]domain.Metric, 0)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (fs *FileStorage) Write(ctx context.Context, metrics []domain.Metric, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(metrics)
}
