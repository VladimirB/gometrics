package handler

import (
	"fmt"
	"net/http"

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

func (h MainPageHandler) GetMainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		fmt.Fprint(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Таблица с метриками</title>
			<head>
			<style>
				table {
  					border-spacing: 5px;
					width: 50%;
				}
				th {
  					text-align: left;
				}
			</style>
			<body>
				<h1>Список метрик</h1>
				<table>
					<tr>
						<th>ID</th>
						<th>Type</th>
						<th>Value</th>
					</tr>
		`)

		for _, metric := range h.storage.GetAll() {
			row := fmt.Sprintf(`
					<tr>
						<td>%s</td>
						<td>%s</td>
						<td>%v</td>
					</tr>
			`, metric.ID, metric.MType, *metric.Value)
			fmt.Fprint(w, row)
		}

		fmt.Fprintf(w, `
				</table>
			</body>
			</html>
		`)
		w.WriteHeader(http.StatusOK)
	}
}
