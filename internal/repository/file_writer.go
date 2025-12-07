package repository

import (
	"encoding/json"
	"os"
)

type JsonFileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewJsonFileWriter(fileName string) (*JsonFileWriter, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &JsonFileWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (fw *JsonFileWriter) Write(v any) error {
	return fw.encoder.Encode(v)
}

func (fw *JsonFileWriter) Close() error {
	return fw.file.Close()
}
