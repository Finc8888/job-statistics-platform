# 🚀 Быстрый старт

## Запуск приложения

### Вариант 1: Docker Compose (рекомендуется)

```bash
# 1. Перейдите в директорию проекта
cd backend

# 2. Запустите приложение
docker-compose up -d

# 3. Проверьте, что контейнеры запустились
docker-compose ps

# 4. Проверьте работу API
curl http://localhost:8081/health
```

API будет доступен по адресу: **http://localhost:8081**

### Вариант 2: Makefile

```bash
# Запустить
make up

# Просмотр логов
make logs

# Остановить
make stop

# Показать все команды
make help
```

## Первые шаги

### 1. Проверьте, что API работает
```bash
curl http://localhost:8081/health
# Должен вернуть: OK
```

### 2. Получите список компаний (тестовые данные)
```bash
curl http://localhost:8081/api/v1/companies
```

### 3. Получите топ навыков
```bash
curl "http://localhost:8081/api/v1/stats/top-skills?limit=10"
```

### 4. Создайте свою компанию
```bash
curl -X POST http://localhost:8081/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{"name": "Моя компания", "description": "Описание"}'
```

## Структура проекта

```
backend/
├── cmd/api/              # Главный файл приложения
├── internal/
│   ├── database/         # Подключение к БД
│   ├── dto/              # Data Transfer Objects (API ↔ модели)
│   ├── models/           # Внутренние модели данных (sql.Null*)
│   ├── repository/       # Репозитории (работа с БД) + интерфейсы
│   └── handlers/         # HTTP обработчики
├── migrations/           # SQL-миграции + seed data
├── TESTING.md            # Документация по тестам
├── docker-compose.yml    # Docker Compose конфигурация
├── Dockerfile            # Docker образ
└── README.md             # Полная документация
```

## Полезные команды

```bash
# Просмотр логов API
docker-compose logs -f api

# Просмотр логов MySQL
docker-compose logs -f mysql

# Войти в MySQL
docker-compose exec mysql mysql -u jobuser -p job_stats
# Пароль: jobpassword

# Перезапустить приложение
docker-compose restart

# Остановить и удалить все
docker-compose down -v
```

## Endpoints

### CRUD операции
- **Companies**: `/api/v1/companies`
- **Jobs**: `/api/v1/jobs`
- **Skills**: `/api/v1/skills`
- **Locations**: `/api/v1/locations`

### Статистика
- **Топ навыков**: `/api/v1/stats/top-skills`
- **Зарплаты**: `/api/v1/stats/skill-salaries`
- **По компаниям**: `/api/v1/stats/companies`
- **Базы данных**: `/api/v1/stats/databases`
- **Языки**: `/api/v1/stats/languages`

## Подробнее

- **README.md** - полная документация
- **TESTING.md** - документация по тестированию
- **API_EXAMPLES.md** - примеры всех запросов

## Troubleshooting

### Порт 3306 занят
Если MySQL порт занят, измените в `docker-compose.yml`:
```yaml
ports:
  - "3307:3306"  # Используйте другой внешний порт
```

### Порт 8081 занят
Измените в `.env`:
```
API_PORT=8081
```

И в `docker-compose.yml`:
```yaml
ports:
  - "8081:8081"
```
