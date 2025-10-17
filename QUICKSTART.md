# Quick Start Guide

## 🚀 Запуск проекта

### 1. Настройка окружения

```bash
# Скопируйте файл с переменными окружения
cp .env.example .env
```

### 2. Запуск PostgreSQL контейнеров

```bash
# Запуск всех контейнеров
docker-compose up -d

# Проверка статуса
docker-compose ps
```

Будут запущены 3 контейнера:
- **psql-mcp-registry** (порт 5434) - БД для хранения зарегистрированных экземпляров
- **psql-mcp-test** (порт 5432) - тестовый PostgreSQL экземпляр "prod"
- **psql-mcp-test-dev** (порт 5433) - тестовый PostgreSQL экземпляр "dev"

### 3. Сборка и запуск приложения

```bash
# Сборка
go build -o psql-mcp-registry .

# Запуск
export $(cat .env | xargs) && ./psql-mcp-registry
```

При успешном запуске вы увидите:
```
Successfully connected to PostgreSQL
Migrations applied successfully!
Initialized MCP server
Starting HTTP API server on :8080
```

## 📊 Тестирование

### Регистрация экземпляров

```bash
# Регистрация PROD экземпляра
curl -X POST http://localhost:8080/api/v1/instances \
  -H "Content-Type: application/json" \
  -d '{
    "name": "prod",
    "database_name": "testdb",
    "description": "Production instance",
    "creator_username": "admin"
  }'

# Регистрация DEV экземпляра
curl -X POST http://localhost:8080/api/v1/instances \
  -H "Content-Type: application/json" \
  -d '{
    "name": "dev",
    "database_name": "devdb",
    "description": "Development instance",
    "creator_username": "admin"
  }'
```

### Проверка здоровья

```bash
curl http://localhost:8080/health
```

## 🔧 Доступные MCP инструменты

Теперь доступно **13 инструментов** для анализа производительности:

### Основные (8):
1. `database_overview` - статистика транзакций и буферов
2. `cache_hit_rate` - эффективность кэша
3. `checkpoints_stats` - статистика чекпоинтов
4. `wal_activity` - активность WAL
5. `tables_info` - статистика таблиц (с bloat-метриками)
6. `locking_info` - текущие блокировки
7. `changed_settings` - измененные настройки
8. `version` - версия PostgreSQL

### Новые (5):
9. `index_stats` - использование индексов
10. `active_queries` - долгие активные запросы
11. `connection_stats` - статистика соединений
12. `slow_queries` - топ медленных запросов (pg_stat_statements)
13. `database_sizes` - размеры баз данных

## 🧪 Прямое тестирование

```bash
# Подключение к test instance
docker exec -it psql-mcp-test psql -U testuser -d testdb

# Проверка pg_stat_statements
testdb=# SELECT * FROM pg_stat_statements LIMIT 3;
testdb=# \dx  -- показать расширения
```

## 🔄 Остановка

```bash
# Остановить контейнеры
docker-compose down

# Остановить и удалить данные
docker-compose down -v
```

## 📝 Примечания

- Registry БД автоматически создает таблицу `instance_registry` при первом запуске
- pg_stat_statements автоматически устанавливается в тестовых экземплярах
- HTTP API работает на порту 8080 (можно изменить в .env)
- MCP сервер работает через stdio для интеграции с Claude Desktop

