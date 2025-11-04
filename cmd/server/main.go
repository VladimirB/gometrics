package main

import (
	"net/http"

	"github.com/VladimirB/gometrics/internal/handler"
)

func main() {
	err := http.ListenAndServe(":8080", handler.NewRouter())
	if err != nil {
		panic(err)
	}
}
