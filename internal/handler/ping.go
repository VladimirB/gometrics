package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type DatabasePingHandler struct {
	db *sql.DB
}

func NewDatabasePingHandler(db *sql.DB) *DatabasePingHandler {
	return &DatabasePingHandler{
		db: db,
	}
}

func (h *DatabasePingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		http.Error(w, "Database not connected", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
