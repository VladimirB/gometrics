BEGIN;
-- Откат создания поля updated_at для отслеживания даты последнего обновления метрики
ALTER TABLE metrics DROP COLUMN IF EXISTS updated_at;
COMMIT;