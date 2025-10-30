# MySQL TUI Editor - Execution Server

HTTP API сервер для безопасного выполнения SQL запросов студентов в изолированных песочницах.

## Возможности

✅ **Песочница**: Каждый запрос выполняется в изолированной временной БД  
✅ **Безопасность**: Блокировка DROP DATABASE, LOAD_FILE, CREATE USER и других опасных команд  
✅ **Форматирование**: Возврат красиво отформатированного текстового вывода MySQL  
✅ **Таймауты**: Автоматическое прерывание запросов после 30 секунд  
✅ **Rate Limiting**: Защита от перегрузки (10 запросов/сек с burst 20)  
✅ **Конкурентность**: Обработка множественных запросов одновременно  

## Быстрый старт

### 1. Запуск MySQL (если ещё не запущен)
```bash
docker run -d \
  --name mysql-tui \
  -e MYSQL_ROOT_PASSWORD=12345 \
  -p 3306:3306 \
  mysql:8.0
```

### 2. Настройка конфигурации
Отредактируйте `config/config.yml`:
```yaml
mysql:
  host: localhost
  port: 3306
  user: root
  password: 12345  # Измените на ваш пароль
```

### 3. Запуск сервера
```bash
# Установка зависимостей
go mod tidy

# Запуск
go run ./cmd/server -config config/config.yml

# Или сборка и запуск бинарника
go build -o bin/server ./cmd/server
./bin/server -config config/config.yml
```

Сервер запустится на `http://localhost:8080`

## API Endpoints

### POST /api/v1/execute
Выполняет SQL запрос в изолированной песочнице.

**Request:**
```json
{
  "query": "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100)); INSERT INTO users VALUES (1, 'John'); SELECT * FROM users;"
}
```

**Response (Success):**
```json
{
  "success": true,
  "output": "Query OK, 0 rows affected\n\nQuery OK, 1 row affected\n\n+----+------+\n| id | name |\n+----+------+\n|  1 | John |\n+----+------+\n1 row in set",
  "execution_time_ms": 45,
  "error": ""
}
```

**Response (Error):**
```json
{
  "success": false,
  "output": "",
  "execution_time_ms": 0,
  "error": "Security validation failed: DROP DATABASE command is not allowed"
}
```

### GET /health
Проверка здоровья сервера.

**Response:**
```json
{
  "status": "healthy",
  "message": "Server is running",
  "time": "2025-10-29T11:30:00+03:00"
}
```

## Безопасность

### Заблокированные команды:
- `DROP DATABASE` / `DROP SCHEMA`
- `SHUTDOWN`
- `LOAD_FILE()`, `INTO OUTFILE`, `INTO DUMPFILE`
- `CREATE USER`, `DROP USER`, `ALTER USER`, `RENAME USER`
- `GRANT`, `REVOKE`
- `SET GLOBAL`, `SET PASSWORD`
- `KILL`
- `INSTALL PLUGIN`, `UNINSTALL PLUGIN`

### Ограничения:
- Максимальное время выполнения запроса: **30 секунд**
- Максимальный размер запроса: **1 МБ**
- Rate limit: **10 запросов/сек** с burst 20

## Архитектура

```
Request → Validator → Executor → Sandbox
                                   ↓
                              Create temp DB
                                   ↓
                              Execute SQL
                                   ↓
                              Drop temp DB
                                   ↓
Response ← Formatted Output ←
```

### Ключевые компоненты:

- **`internal/api/`** - HTTP handlers и middleware (Gin framework)
- **`internal/executor/`** - Выполнение SQL и управление песочницами
- **`internal/security/`** - Валидация и блокировка опасных команд
- **`internal/domain/`** - Модели данных (Request/Response)
- **`internal/config/`** - Загрузка конфигурации (Viper)

## Тестирование

Запустите сервер и используйте примеры из `TEST_EXAMPLES.md`:

```bash
# Health check
curl http://localhost:8080/health

# Простой запрос
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"query":"SELECT 1 + 1 as result;"}'
```

## Production Deployment

### Docker
```bash
# Собрать образ
docker build -t mysql-tui-server .

# Запустить
docker run -d \
  -p 8080:8080 \
  -e MYSQL_HOST=mysql \
  -e MYSQL_PASSWORD=your_password \
  mysql-tui-server
```

### С Docker Compose
```bash
# Запустить MySQL и сервер
docker-compose up -d
```

## Мониторинг

Логи выводятся в stdout и содержат:
- HTTP запросы (метод, путь, статус, время)
- SQL выполнение (успех/неудача, время выполнения)
- Ошибки подключения к MySQL
- Предупреждения о cleanup временных БД

## Производительность

- Connection pool: 25 одновременных соединений
- Каждый запрос получает изолированную БД
- Автоматическая очистка через defer
- Graceful shutdown с завершением активных запросов

## Troubleshooting

### "Failed to ping MySQL"
- Проверьте, что MySQL запущен: `docker ps`
- Проверьте хост и порт в `config/config.yml`
- Проверьте пароль

### "Address already in use"
- Порт 8080 занят: `lsof -ti:8080 | xargs kill -9`

### "Query execution timeout"
- Запрос выполняется дольше 30 секунд
- Оптимизируйте запрос или увеличьте таймаут в конфиге
