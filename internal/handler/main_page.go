package handler

import (
	"fmt"
	"html/template"
	"net/http"

	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/service"
)

type MainPageHandler struct {
	service  *service.MetricsService
	template *template.Template
}

func NewMainPageHandler(service *service.MetricsService, template *template.Template) *MainPageHandler {
	return &MainPageHandler{
		service:  service,
		template: template,
	}
}

type Raw struct {
	ID string
	Type string
	Value string
}

type MainPageData struct {
	Raws []Raw
}

func (h MainPageHandler) GetMainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		metrics := h.service.GetAll()
		pageData := preparePageData(metrics)
		err := h.template.Execute(w, pageData)
		if err != nil {
			http.Error(w, "error on parsing of main page: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func preparePageData(metrics map[string]models.Metrics) MainPageData {
	var raws []Raw
	for _, m := range metrics {
		raw := Raw{
			ID: m.ID,
			Type: m.MType,
		}

		if m.MType == models.Counter {
			raw.Value = fmt.Sprint(*m.Delta)
		} else {
			raw.Value = fmt.Sprint(*m.Value)
		}
		raws = append(raws, raw)
	}

	return MainPageData{Raws: raws}
}
