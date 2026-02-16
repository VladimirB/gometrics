BEGIN;
-- Добавляю поле updated_at для отслеживания даты последнего обновления метрики
ALTER TABLE metrics ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
COMMIT;
