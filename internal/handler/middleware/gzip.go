package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/VladimirB/gometrics/internal/logger"
	"go.uber.org/zap"
)

type GZipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

type GZipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewGZipWriter(w http.ResponseWriter) *GZipWriter {
	return &GZipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func NewGZipReader(r io.ReadCloser) (*GZipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GZipReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c *GZipWriter) Header() http.Header {
	return c.w.Header()
}

func (c *GZipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *GZipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

func (c *GZipWriter) Close() error {
	return c.zw.Close()
}

func (c GZipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *GZipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GZip(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Сжатие работает только для указанных типов контента
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" && contentType != "text/html" {
			h.ServeHTTP(w, r)
			return
		}

		usedWriter := w

		// Если клиент поддерживает прием gzip, то устанавлияваем gzip writer как основной
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			// gzw := NewGZipWriter(w)
			// usedWriter = gzw
			// defer gzw.Close()
		}

		// Если клиент отправляет нам gzip, то используем reader с поддержкой декомпрессии
		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			gzr, err := NewGZipReader(r.Body)
			if err != nil {
				logger.Log.Error("cant read body as gzip", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = gzr

			defer gzr.Close()
		}

		logger.Log.Info("GZip middleware",
			zap.String("Accept-Encoding", acceptEncoding),
			zap.String("Content-Encoding", contentEncoding))

		h.ServeHTTP(usedWriter, r)
	}
}
