package database

import (
	"fmt"
	"time"

	"github.com/VladimirB/gometrics/internal/domain"
)

type metricDB struct {
	ID        string    `db:"id"`
	Type      string    `db:"type"`
	Delta     *int64    `db:"delta"`
	Value     *float64  `db:"value"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *metricDB) ToDomain() *domain.Metric {
	return &domain.Metric{
		ID:    m.ID,
		MType: m.Type,
		Delta: m.Delta,
		Value: m.Value,
	}
}

func FromDomain(m *domain.Metric, updatedAt time.Time) *metricDB {
	return &metricDB{
		ID:        m.ID,
		Type:      m.MType,
		Delta:     m.Delta,
		Value:     m.Value,
		UpdatedAt: updatedAt,
	}
}

func (m *metricDB) String() string {
	return fmt.Sprintf("MetricDB(id: %s, type: %s, delta: %p, value: %p, updated_at: %s)", m.ID, m.Type, m.Delta, m.Value, m.UpdatedAt)
}
