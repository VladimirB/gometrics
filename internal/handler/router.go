package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(updateHandler *UpdateMetricHandler, valueHandler *ValueMetricHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/", func(r chi.Router) {
		r.Route("/update", func(r chi.Router) {
			r.Post("/{metricType}/{metricValue}", updateHandler.PostMetricNoNameHandler())
			r.Post("/{metricType}/{metricName}/{metricValue}", updateHandler.PostMetricHandler())
		})

		r.Route("/value", func(r chi.Router) {
			r.Get("/{metricType}/{metricName}", valueHandler.GetMetricHandler())
		})
	})

	return router
}