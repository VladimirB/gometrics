package database

import (
	"context"
	"fmt"
	"time"

	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Save(ctx context.Context, metric domain.Metric) error {
	metricDB := FromDomain(&metric, time.Now().UTC())

	query := `
		INSERT INTO metrics (id, type, delta, value, updated_at)
		VALUES (:id, :type, :delta, :value, :updated_at)
		ON CONFLICT (id)
		DO UPDATE SET
			delta = EXCLUDED.delta,
			value = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.NamedExecContext(ctx, query, metricDB)
	if err != nil {
		return fmt.Errorf("error on save metric in db: %w", err)
	}

	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, metricID string) (domain.Metric, error) {
	var row metricDB

	query := "SELECT id, type, delta, value FROM metrics WHERE id = $1"

	if err := r.db.GetContext(ctx, &row, query, metricID); err != nil {
		return domain.Metric{}, fmt.Errorf("metric %s not found in db: %w", metricID, err)
	}

	return *row.ToDomain(), nil
}

func (r *PostgresRepo) GetAll(ctx context.Context) (map[string]domain.Metric, error) {
	var rows []metricDB

	query := "SELECT id, type, delta, value FROM metrics"

	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("error on obtain all metrics from db: %w", err)
	}

	result := make(map[string]domain.Metric, len(rows))
	for _, row := range rows {
		metric := row.ToDomain()
		result[metric.ID] = *metric
	}

	return result, nil
}
