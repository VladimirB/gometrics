BEGIN;
-- Создаем тип для метрики
CREATE TYPE metric_type AS ENUM ('gauge', 'counter');

-- Таблица для хранения метрик
CREATE TABLE metrics (
    id VARCHAR(30) PRIMARY KEY,
    type metric_type NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION
);
COMMIT;