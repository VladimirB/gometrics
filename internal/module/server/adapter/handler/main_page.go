package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/module/server/port"
	"github.com/VladimirB/gometrics/internal/shared/logger"
	"go.uber.org/zap"
)

type MainPageHandler struct {
	metricService port.MetricService
	template      *template.Template
}

func NewMainPageHandler(service port.MetricService, template *template.Template) *MainPageHandler {
	return &MainPageHandler{
		metricService: service,
		template:      template,
	}
}

type Raw struct {
	ID    string
	Type  string
	Value string
}

type MainPageData struct {
	Raws []Raw
}

func (h MainPageHandler) GetMainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		metrics, err := h.metricService.GetAll(r.Context())
		if err != nil {
			logger.Log.Error("cant get metrics for main page", zap.Error(err))
			http.Error(w, "cant get metrics for main page", http.StatusInternalServerError)
			return
		}

		pageData := preparePageData(metrics)
		err = h.template.Execute(w, pageData)
		if err != nil {
			logger.Log.Error("failed parse main page", zap.Error(err))
			http.Error(w, "failed parse main page: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func preparePageData(metrics map[string]domain.Metric) MainPageData {
	var raws []Raw
	for _, m := range metrics {
		raw := Raw{
			ID:   m.ID,
			Type: m.MType,
		}

		if m.MType == domain.Counter {
			raw.Value = fmt.Sprint(*m.Delta)
		} else {
			raw.Value = fmt.Sprint(*m.Value)
		}
		raws = append(raws, raw)
	}

	return MainPageData{Raws: raws}
}
