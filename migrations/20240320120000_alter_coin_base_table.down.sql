-- migrations/20240320120000_alter_coin_base_table.down.sql
BEGIN;

-- Удаляем ограничения
ALTER TABLE coin_base
    DROP CONSTRAINT IF EXISTS unique_title_date;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_title;
DROP INDEX IF EXISTS idx_date;

-- Удаляем добавленные колонки
ALTER TABLE coin_base
    DROP COLUMN IF EXISTS id,
    DROP COLUMN IF EXISTS created_at;

COMMIT;