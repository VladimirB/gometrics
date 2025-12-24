-- Создаем тип для метрики
CREATE TYPE metric_type AS ENUM ('Gauge', 'Counter');

-- Таблица для хранения метрик
CREATE TABLE metrics (
    id VARCHAR(15) PRIMARY KEY,
    type metric_type NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION
);