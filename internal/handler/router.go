package handler

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(mainPageHandler *MainPageHandler,
	updateHandler *UpdateMetricHandler,
	valueHandler *ValueMetricHandler) *chi.Mux {

	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/", RequestLogger(mainPageHandler.GetMainPage()))

		r.Route("/update", func(r chi.Router) {
			r.Post("/{metricType}/{metricValue}", RequestLogger(updateHandler.PostMetricNoNameHandler()))
			r.Post("/{metricType}/{metricName}/{metricValue}", RequestLogger(updateHandler.PostMetricHandler()))
		})

		r.Route("/value", func(r chi.Router) {
			r.Get("/{metricType}/{metricName}", RequestLogger(valueHandler.GetMetricHandler()))
		})
	})

	return router
}
