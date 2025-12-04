package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func Zip(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data: %v", err)
	}

	err := w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	return buf.Bytes(), nil
}

func Unzip(data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	return io.ReadAll(gzReader)
}
