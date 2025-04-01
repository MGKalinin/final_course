-- migrations/20240320120000_alter_coin_base_table.up.sql
BEGIN;

-- Добавляем недостающие колонки
ALTER TABLE coin_base
    ADD COLUMN IF NOT EXISTS id SERIAL PRIMARY KEY,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Удаляем дубликаты (оставляем последнюю запись для каждой пары title-date)
DELETE FROM coin_base
WHERE ctid NOT IN (
    SELECT max(ctid)
    FROM coin_base
    GROUP BY title, date
);

-- Добавляем ограничения
ALTER TABLE coin_base
    ADD CONSTRAINT unique_title_date UNIQUE (title, date);

-- Создаем индексы
CREATE INDEX IF NOT EXISTS idx_title ON coin_base (title);
CREATE INDEX IF NOT EXISTS idx_date ON coin_base (date);

COMMIT;