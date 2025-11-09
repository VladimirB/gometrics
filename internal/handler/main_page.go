package handler

import (
	"html/template"
	"net/http"

	models "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
)

type MainPageHandler struct {
	storage *repository.MemStorage
}

func NewMainPageHandler(storage *repository.MemStorage) *MainPageHandler {
	return &MainPageHandler{
		storage: storage,
	}
}

type MainPageData struct {
	Metrics map[string]models.Metrics
}

var mainPage *template.Template

func init() {
	mainPage = template.Must(template.ParseFiles("web/template/index.html"))
}

func (h MainPageHandler) GetMainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		metrics := h.storage.GetAll()
		pageData := MainPageData{
			Metrics: metrics,
		}

		err := mainPage.Execute(w, pageData)
		if err != nil {
			http.Error(w, "error on parsing of main page: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
