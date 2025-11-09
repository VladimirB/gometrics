package handler

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

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
	// Не удается пройти тесты на CI. Предполагаю, что потому что в пайплайне тестов запуск производится
	// не из корня проекта: cd cmd/server 
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Не удалось получить текущий путь к файлу")
	}
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filename))) 
	absoulutePath := filepath.Join(rootDir, "web/template/index.html")

	mainPage = template.Must(template.ParseFiles(absoulutePath))
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
