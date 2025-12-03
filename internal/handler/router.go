package handler

import (
	"github.com/VladimirB/gometrics/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(mainPageHandler *MainPageHandler,
	updateHandler *UpdateMetricHandler,
	valueHandler *ValueMetricHandler) *chi.Mux {

	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/", RequestLogger(mainPageHandler.GetMainPage()))

		r.Route("/update", func(r chi.Router) {
			r.Post("/", RequestLogger(middleware.GZip(updateHandler.UpdateMetricJSONHandler())))
			r.Post("/{metricType}/{metricValue}", RequestLogger(updateHandler.UpdateNoNameMetricHandler()))
			r.Post("/{metricType}/{metricName}/{metricValue}", RequestLogger(updateHandler.UpdateMetricByNameHandler()))
		})

		r.Route("/value", func(r chi.Router) {
			r.Post("/", RequestLogger(middleware.GZip(valueHandler.PostValueMetricHandler())))
			r.Get("/{metricType}/{metricName}", RequestLogger(valueHandler.GetMetricHandler()))
		})
	})

	return router
}
