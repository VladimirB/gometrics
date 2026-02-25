package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type DatabasePingHandler struct {
	db *sqlx.DB
}

func NewDatabasePingHandler(db *sqlx.DB) *DatabasePingHandler {
	return &DatabasePingHandler{
		db: db,
	}
}

func (h *DatabasePingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		http.Error(w, "Database not connected", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
