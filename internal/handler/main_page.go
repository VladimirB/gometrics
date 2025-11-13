package handler

import (
	"html/template"
	"net/http"

	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
)

type MainPageHandler struct {
	storage *repository.MemStorage
	template *template.Template
}

func NewMainPageHandler(storage *repository.MemStorage, template *template.Template) *MainPageHandler {
	return &MainPageHandler{
		storage: storage,
		template: template,
	}
}

type MainPageData struct {
	Metrics map[string]models.Metrics
}

func (h MainPageHandler) GetMainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		metrics := h.storage.GetAll()
		pageData := MainPageData{
			Metrics: metrics,
		}

		err := h.template.Execute(w, pageData)
		if err != nil {
			http.Error(w, "error on parsing of main page: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
