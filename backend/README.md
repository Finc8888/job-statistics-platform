# Job Statistics API

REST API для управления вакансиями и анализа статистики рынка труда IT-разработчиков.

## 🚀 Возможности

- CRUD операции для компаний, вакансий, навыков и локаций
- Статистика по самым востребованным навыкам
- Анализ зарплат по навыкам
- Статистика по языкам программирования и базам данных
- DTO-слой для разделения API-контрактов и внутренних моделей
- Unit-тесты (sqlmock + httptest)
- Docker-контейнеризация
- MySQL база данных

## 📋 Требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)

## 🛠️ Установка и запуск

### Запуск через Docker Compose (рекомендуется)

```bash
# Клонируйте репозиторий
git clone <repository-url>
cd backend

# Запустите приложение
docker-compose up -d

# Проверьте статус
docker-compose ps

# Просмотр логов
docker-compose logs -f api
```

API будет доступен по адресу: `http://localhost:8081`

### Локальная разработка

```bash
# Установите зависимости
go mod download

# Запустите MySQL (через Docker)
docker-compose up -d mysql

# Запустите API
go run cmd/api/main.go
```

## 🏗️ Архитектура

```
HTTP Request → Handler → DTO → Model → Repository → MySQL
HTTP Response ← Handler ← DTO ← Model ← Repository ← MySQL
```

```
backend/
├── cmd/api/main.go              # Точка входа, wiring
├── internal/
│   ├── database/db.go           # Подключение к MySQL
│   ├── dto/                     # Data Transfer Objects
│   │   ├── job.go               # JobRequest, JobResponse, маппинг
│   │   └── location.go          # LocationRequest, LocationResponse, маппинг
│   ├── models/models.go         # Внутренние модели (sql.Null* для nullable колонок)
│   ├── repository/              # SQL-запросы + интерфейсы для DI
│   │   ├── interfaces.go
│   │   ├── job_repository.go
│   │   ├── company_repository.go
│   │   ├── skill_repository.go
│   │   ├── location_repository.go
│   │   └── stats_repository.go
│   └── handlers/                # HTTP-хендлеры
│       ├── job_handler.go
│       ├── company_handler.go
│       ├── skill_location_handler.go
│       └── stats_handler.go
├── migrations/
│   ├── 001_create_tables.sql    # Схема (idempotent)
│   └── 002_seed_data.sql        # Тестовые данные (destructive)
└── TESTING.md                   # Документация по тестам
```

### DTO-слой

Пакет `internal/dto/` отвечает за преобразование данных между JSON (API) и внутренними моделями:

- `dto.JobRequest` — входящий JSON → `models.Job` (через `ToModel()`)
- `dto.JobResponse` — `models.Job` → исходящий JSON (через `JobResponseFromModel()`)

Nullable поля (`sql.NullString`, `sql.NullFloat64`) в моделях представлены как указатели (`*string`, `*float64`) в DTO и сериализуются как значение или `null`.

## 📚 API Endpoints

### Companies (Компании)

- `GET /api/v1/companies` - Получить все компании
- `GET /api/v1/companies/{id}` - Получить компанию по ID
- `POST /api/v1/companies` - Создать компанию
- `PUT /api/v1/companies/{id}` - Обновить компанию
- `DELETE /api/v1/companies/{id}` - Удалить компанию

**Пример создания компании:**
```bash
curl -X POST http://localhost:8081/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Google",
    "description": "Технологическая компания"
  }'
```

### Jobs (Вакансии)

- `GET /api/v1/jobs` - Получить все вакансии
- `GET /api/v1/jobs/{id}` - Получить вакансию по ID
- `POST /api/v1/jobs` - Создать вакансию
- `PUT /api/v1/jobs/{id}` - Обновить вакансию
- `DELETE /api/v1/jobs/{id}` - Удалить вакансию

**Пример создания вакансии:**
```bash
curl -X POST http://localhost:8081/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": 1,
    "title": "Senior Go Developer",
    "level": "Senior",
    "specialization": "Golang",
    "salary_min": 300000,
    "salary_max": 500000,
    "salary_currency": "RUB",
    "experience_years": "5+",
    "remote_available": true,
    "is_active": true
  }'
```

### Skills (Навыки)

- `GET /api/v1/skills` - Получить все навыки
- `GET /api/v1/skills/{id}` - Получить навык по ID
- `POST /api/v1/skills` - Создать навык
- `PUT /api/v1/skills/{id}` - Обновить навык
- `DELETE /api/v1/skills/{id}` - Удалить навык

**Пример создания навыка:**
```bash
curl -X POST http://localhost:8081/api/v1/skills \
  -H "Content-Type: application/json" \
  -d '{
    "name": "TypeScript",
    "category": "Язык программирования"
  }'
```

### Locations (Локации)

- `GET /api/v1/locations` - Получить все локации
- `GET /api/v1/locations/job/{job_id}` - Получить локации для вакансии
- `POST /api/v1/locations` - Создать локацию
- `PUT /api/v1/locations/{id}` - Обновить локацию
- `DELETE /api/v1/locations/{id}` - Удалить локацию

### Statistics (Статистика)

#### Топ навыков
```bash
GET /api/v1/stats/top-skills?limit=10
```
Возвращает топ самых востребованных навыков.

#### Зарплаты по навыкам
```bash
GET /api/v1/stats/skill-salaries?min_vacancies=1
```

#### Навыки по уровням
```bash
GET /api/v1/stats/skills-by-level
```

#### Статистика компаний
```bash
GET /api/v1/stats/companies
```

#### Статистика баз данных
```bash
GET /api/v1/stats/databases
```

#### Статистика языков программирования
```bash
GET /api/v1/stats/languages
```

## 🧪 Тестирование

```bash
# Все тесты
go test ./... -v -count=1

# С покрытием
make test-coverage

# В Docker
make test-docker
```

Подробнее — см. [TESTING.md](TESTING.md).
