package handler

import (
	"github.com/VladimirB/gometrics/internal/module/server/adapter/handler/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(mainPageHandler *MainPageHandler,
	updateHandler *UpdateMetricHandler,
	valueHandler *ValueMetricHandler,
	dbPingHandler *DatabasePingHandler) *chi.Mux {

	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/", middleware.RequestLogger(middleware.GZip(mainPageHandler.GetMainPage())))

		r.Route("/update", func(r chi.Router) {
			r.Post("/", middleware.RequestLogger(middleware.GZip(updateHandler.UpdateMetricJSONHandler())))
			r.Post("/{metricType}/{metricValue}", middleware.RequestLogger(updateHandler.UpdateNoNameMetricHandler()))
			r.Post("/{metricType}/{metricName}/{metricValue}", middleware.RequestLogger(updateHandler.UpdateMetricByNameHandler()))
		})

		r.Post("/updates", middleware.RequestLogger(middleware.GZip(updateHandler.UpdateMetricsBunchHandler())))

		r.Route("/value", func(r chi.Router) {
			r.Post("/", middleware.RequestLogger(middleware.GZip(valueHandler.PostValueMetricHandler())))
			r.Get("/{metricType}/{metricName}", middleware.RequestLogger(valueHandler.GetMetricHandler()))
		})

		r.Get("/ping", middleware.RequestLogger(dbPingHandler.Ping))
	})

	return router
}
