package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricValue}", PostMetricNoNameHandler())
		r.Post("/{metricType}/{metricName}/{metricValue}", PostMetricHandler())
	})
	return router
}