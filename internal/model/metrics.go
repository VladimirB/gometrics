package models

import "fmt"

const (
	Counter = "counter"

	PollCount = "PollCount"
)

const (
	Gauge = "gauge"

	Alloc         = "Alloc"
	BuckHashSys   = "BuckHashSys"
	Frees         = "Frees"
	GCCPUFraction = "GCCPUFraction"
	GCSys         = "GCSys"
	HeapAlloc     = "HeapAlloc"
	HeapIdle      = "HeapIdle"
	HeapInuse     = "HeapInuse"
	HeapObjects   = "HeapObjects"
	HeapReleased  = "HeapReleased"
	HeapSys       = "HeapSys"
	LastGC        = "LastGC"
	Lookups       = "Lookups"
	MCacheInuse   = "MCacheInuse"
	MSpanSys      = "MSpanSys"
	Mallocs       = "Mallocs"
	NextGC        = "NextGC"
	NumForcedGC   = "NumForcedGC"
	NumGC         = "NumGC"
	OtherSys      = "OtherSys"
	PauseTotalNs  = "PauseTotalNs"
	StackInuse    = "StackInuse"
	StackSys      = "StackSys"
	Sys           = "Sys"
	TotalAlloc    = "TotalAlloc"
	RandomValue   = "RandomValue"
)

var allowedIds = map[string]bool {
	PollCount: true,
	Alloc: true,
	BuckHashSys: true,
	Frees: true,
	GCCPUFraction: true,
	GCSys: true,
	HeapAlloc: true,
	HeapIdle: true,
	HeapInuse: true,
	HeapObjects: true,
	HeapReleased: true,
	HeapSys: true,
	LastGC: true,
	Lookups: true,
	MCacheInuse: true,
	MSpanSys: true,
	Mallocs: true,
	NextGC: true,
	NumForcedGC: true,
	NumGC: true,
	OtherSys: true,
	PauseTotalNs: true,
	StackInuse: true,
	StackSys: true,
	Sys: true,
	TotalAlloc: true,
	RandomValue: true,
}

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`
}

func (m Metrics) Validate() error {
	if !allowedIds[m.ID] {
		return fmt.Errorf("not allowed metric ID")
	}

	if m.MType != Gauge && m.MType != Counter {
		return fmt.Errorf("not allowed metric type")
	}

	return nil
}
