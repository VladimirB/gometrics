package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/VladimirB/gometrics/internal/domain"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Save(ctx context.Context, metric domain.Metric) error {
	query := `
		INSERT INTO metrics (id, type, delta, value, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id)
		DO UPDATE SET
			delta = EXCLUDED.delta,
			value = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Delta, metric.Value, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("error on save metric in db: %w", err)
	}

	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, metricID string) (domain.Metric, error) {
	query := "SELECT id, type, delta, value FROM metrics WHERE id = $1"

	row := r.db.QueryRowContext(ctx, query, metricID)

	metric := domain.Metric{}
	err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if err != nil {
		return metric, fmt.Errorf("metric %s not found in db: %w", metricID, err)
	}

	return metric, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context) (map[string]domain.Metric, error) {
	query := "SELECT id, type, delta, value FROM metrics";

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error on obtain all metrics from db: %w", err)
	}
	defer rows.Close()

	metrics := make(map[string]domain.Metric)

	for rows.Next() {
		metric := domain.Metric{}
		err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		metrics[metric.ID] = metric
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return metrics, nil
}
