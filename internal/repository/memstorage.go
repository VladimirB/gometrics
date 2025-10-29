package repository

type MemStorage struct {
	storage map[string]any
}

func (ms MemStorage) Save(metricName string, value any) {
	ms.storage[metricName] = value
}

func (ms MemStorage) GetGauge(metricName string) *float64 {
	var value float64
	var ok bool
	if value, ok = ms.storage[metricName].(float64); !ok {
		return nil
	}

	return &value
}

func (ms MemStorage) GetCounter(metricName string) *int {
	var value int
	var ok bool
	if value, ok = ms.storage[metricName].(int); !ok {
		return nil
	}

	return &value
}