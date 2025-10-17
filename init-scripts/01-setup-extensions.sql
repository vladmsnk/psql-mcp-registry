-- Включаем pg_stat_statements для анализа медленных запросов
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Создаем простую тестовую таблицу
CREATE TABLE IF NOT EXISTS test_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON test_users(username);

-- Добавляем тестовые данные
INSERT INTO test_users (username, email)
SELECT 
    'user_' || i,
    'user_' || i || '@example.com'
FROM generate_series(1, 100) AS i
ON CONFLICT DO NOTHING;

