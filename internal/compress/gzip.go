package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
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
