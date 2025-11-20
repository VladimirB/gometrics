package handler

import (
	"net/http"
	"time"

	"github.com/VladimirB/gometrics/internal/logger"
)

func RequestLogger(h http.HandlerFunc) http.HandlerFunc {
	logFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Log.Sugar().Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	}
	return http.HandlerFunc(logFunc)
}