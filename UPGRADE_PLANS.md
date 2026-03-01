# Планы по улучшению проекта

## Текущее состояние

Проект имеет чистую слоистую архитектуру: **Handler → DTO → Repository → MySQL**, Docker-контейнеризацию, TypeScript strict mode, паттерн Repository с интерфейсами для DI. Реализован DTO-слой для разделения API-контрактов и внутренних моделей. Backend покрыт unit-тестами (repository + handlers) через `sqlmock` и `httptest`. Основные пробелы — отсутствие аутентификации, CI/CD, фронтенд-тестов и пагинации.

### Что уже сделано

- ✅ Repository-паттерн с интерфейсами (`JobRepositoryInterface`, `CompanyRepositoryInterface`, `JobSkillRepositoryInterface`)
- ✅ DTO-слой (`internal/dto/`) — разделение API-контрактов и внутренних моделей, маппинг `sql.Null*` ↔ `*string`/`*float64`
- ✅ Backend unit-тесты (49 тестов: repository через `sqlmock`, handlers через `httptest` + mock-репозитории)
- ✅ Docker: full-stack compose, backend-only compose, test compose; `make rebuild[-frontend|-api]` для форс-пересборки без кэша
- ✅ Локальная разработка (`make dev-local`, `make dev-api`, `make dev-frontend`)
- ✅ Миграции разделены на schema (`migrate.sh`) и seed (`seed.sh`)
- ✅ Навыки вакансии: `GET/POST /api/v1/jobs/{id}/skills`, мультиселект в форме, специализация из списка навыков, фильтр по вакансии в разделе Навыки

---

## 1. Переход на Domain-Driven Design (DDD)

Текущая архитектура с DTO-слоем уже подготовлена к плавному переходу на DDD. Основная идея — выделить чистые доменные сущности без инфраструктурных зависимостей.

### 1.1 Целевая архитектура

```
Текущая:     Handler → DTO ↔ Model (sql.Null*) → Repository → MySQL
Целевая:     Handler → DTO ↔ Domain Entity (чистый Go) → Repository (sql.Null*) → MySQL
```

Целевая структура пакетов:

```
backend/internal/
├── domain/                  # Чистые доменные сущности и бизнес-логика
│   ├── job.go               # Job entity + value objects (Salary, Level)
│   ├── company.go           # Company entity
│   ├── skill.go             # Skill entity
│   ├── location.go          # Location entity
│   ├── errors.go            # Доменные ошибки (ErrNotFound, ErrValidation)
│   └── repository.go        # Интерфейсы репозиториев (перенос из repository/interfaces.go)
├── dto/                     # Без изменений — уже работает с чистыми Go-типами
│   ├── job.go
│   └── location.go
├── repository/              # Инфраструктурный слой — sql.Null* живёт только здесь
│   ├── job_repository.go    # Маппинг domain.Job ↔ SQL rows
│   ├── company_repository.go
│   └── ...
├── service/                 # Доменные сервисы (бизнес-логика, оркестрация)
│   ├── job_service.go
│   └── stats_service.go
└── handlers/                # HTTP-слой — работает только с DTO и сервисами
```

### 1.2 Этапы миграции

**Этап 1: Выделение доменных сущностей**

Создать `internal/domain/` с чистыми структурами без `sql.Null*`, без `json`-тегов:

```go
// internal/domain/job.go
package domain

type Salary struct {
    Min      *float64
    Max      *float64
    Currency string
}

type Job struct {
    ID               int
    CompanyID        int
    Title            string
    Level            Level
    Specialization   string
    Salary           Salary
    ExperienceYears  string
    Location         string
    RemoteAvailable  bool
    Description      string
    IsActive         bool
    // ...
}

// Бизнес-методы
func (j Job) HasSalary() bool {
    return j.Salary.Min != nil || j.Salary.Max != nil
}

func (j Job) Validate() error {
    if j.Title == "" {
        return ErrValidation{Field: "title", Message: "title is required"}
    }
    // ...
}
```

**Этап 2: Перенос интерфейсов в domain**

Перенести `repository/interfaces.go` → `domain/repository.go`, изменив сигнатуры на доменные типы:

```go
// internal/domain/repository.go
package domain

type JobRepository interface {
    GetAll() ([]Job, error)
    GetByID(id int) (*Job, error)
    Create(j *Job) error
    Update(j *Job) error
    Delete(id int) error
}
```

**Этап 3: Маппинг в репозитории**

Добавить в каждый репозиторий приватные функции маппинга `domain ↔ sql.Null*`:

```go
// internal/repository/job_repository.go
func toDomain(row jobRow) domain.Job { ... }
func fromDomain(j domain.Job) jobRow { ... }
```

**Этап 4: Доменные сервисы**

Вынести бизнес-логику из хендлеров в сервисы:

```go
// internal/service/job_service.go
type JobService struct {
    repo domain.JobRepository
}

func (s *JobService) Create(j *domain.Job) error {
    if err := j.Validate(); err != nil {
        return err
    }
    return s.repo.Create(j)
}
```

**Этап 5: Обновление DTO-маперов**

Изменить маппинг `dto ↔ domain` вместо `dto ↔ model`. Поскольку DTO уже работает с `*string`/`*float64`, а домен тоже будет использовать чистые Go-типы, изменения будут минимальными.

### 1.3 Принципы миграции

- **Поэтапно, без big bang** — каждый этап можно мержить отдельно, не ломая работающий код
- **DTO-слой не меняется** — он уже изолирован от инфраструктуры
- **Тесты остаются зелёными** — на каждом этапе все тесты должны проходить
- **domain не импортирует** `database/sql`, `encoding/json`, `net/http` — это главное правило
- **Зависимости направлены внутрь** — domain ни от кого не зависит; repository, handlers, dto зависят от domain

### 1.4 Value Objects для рассмотрения

| Value Object | Поля | Где используется |
|---|---|---|
| `Salary` | Min, Max, Currency | Job |
| `Level` | string enum с валидацией | Job |
| `DateRange` | From, To | фильтрация, статистика |
| `Pagination` | Page, PerPage, Total | API-ответы (после реализации пагинации) |

---

## 2. Тестирование

- ✅ **Backend unit-тесты** — repository (`sqlmock`), handlers (`httptest` + mock-репозитории)
- ⬜ **DTO mapper тесты** — unit-тесты для `ToModel()`, `*ResponseFromModel()` (рекомендуется добавить перед DDD-миграцией)
- ⬜ **Frontend тесты** — Jest + React Testing Library для компонентов, тесты для MobX store
- ⬜ **E2E тесты** — Playwright или Cypress для критических сценариев (CRUD вакансий, статистика)
- ⬜ **Domain тесты** (после DDD) — тесты бизнес-логики и валидации в `domain/`

---

## 3. Безопасность

- ⬜ **Аутентификация/авторизация** — API полностью открыт. Добавить JWT или session-based auth
- ⬜ **CORS** — сейчас `*` (все origins разрешены). Ограничить до конкретных доменов
- ⬜ **Секреты** — пароли БД захардкожены в docker-compose и .env. Использовать Docker secrets или Vault
- ⬜ **Input validation** — минимальная валидация на бэкенде (при DDD переедет в `domain.Validate()`)
- ⬜ **Rate limiting** — нет защиты от перебора запросов

---

## 4. CI/CD

- ⬜ Проект не является git-репозиторием — начать с `git init`
- ⬜ Добавить GitHub Actions / GitLab CI:
  - Линтинг (`go vet`, `golangci-lint`, `eslint`)
  - Запуск тестов
  - Сборка Docker-образов
  - Деплой на staging/production

---

## 5. API

- ⬜ **OpenAPI/Swagger** — добавить спецификацию (сейчас документация только в Markdown)
- ⬜ **Пагинация** — эндпоинты `/jobs`, `/companies` возвращают все записи без лимита
- ⬜ **Фильтрация и сортировка** — для списковых эндпоинтов (по уровню, зарплате, навыкам)
- ⬜ **Structured logging** — заменить `log.Println` на `slog` (встроен в Go 1.21) с JSON-форматом
- ⬜ **Версионирование** — уже есть `/api/v1`, добавить middleware для обработки версий

---

## 6. База данных

- ⬜ **Миграции** — заменить ручной `migrate.sh` на `golang-migrate` или `goose` для версионирования
- ⬜ **Индексы** — добавить индексы на `jobs.level`, `jobs.specialization`, `jobs.salary_min/max`
- ⬜ **Бэкапы** — реализовать стратегию резервного копирования

---

## 7. Frontend

- ⬜ **Обработка ошибок** — добавить retry-логику и user-friendly сообщения при сбоях API
- ⬜ **Валидация форм** — использовать `react-hook-form` + `zod` вместо минимальной проверки
- ⬜ **Кеширование** — рассмотреть `tanstack-query` вместо ручных fetch в MobX (автокеш, refetch)
- ⬜ **Lazy loading** — `React.lazy()` + `Suspense` для code splitting страниц
- ⬜ **Accessibility** — добавить aria-атрибуты, проверить контрастность и навигацию с клавиатуры

---

## 8. Инфраструктура и мониторинг

- ⬜ **Health checks** — добавить реальную проверку соединения с БД в эндпоинт `/health`
- ⬜ **Мониторинг** — Prometheus metrics + Grafana для визуализации
- ⬜ **Логирование** — централизованный сбор логов (ELK / Loki)
- ⬜ **Graceful shutdown** — обработка SIGTERM в Go-сервере для корректного завершения

---

## 9. Функциональность
- ✅ **Специализация** — dropdown из списка навыков (не free text)
- ✅ **Мультиселект навыков** — при создании/редактировании вакансии
- ⬜ **Добавить поле 'Ссылка на вакансию'**
- ⬜ **Парсинг вакансий** — добавить парсер (hh.ru API, Habr Career) для автоматического сбора данных
- ⬜ **Экспорт данных** — CSV/Excel выгрузка статистики
- ⬜ **Фильтры на дашборде** — фильтрация по периоду, региону, уровню
- ⬜ **Сравнение** — сравнение навыков/зарплат между периодами

---

## Рекомендуемый порядок внедрения

| Приоритет | Задача | Зависимости |
|---|---|---|
| 1 | `git init` + CI/CD pipeline | — |
| 2 | DTO mapper тесты | — |
| 3 | JWT-аутентификация | — |
| 4 | Пагинация и фильтрация API | — |
| 5 | Input validation (DTO-уровень) | — |
| 6 | Frontend: валидация форм + тесты | — |
| 7 | OpenAPI спецификация | — |
| 8 | **DDD Этап 1:** `internal/domain/` — доменные сущности | Стабильный API |
| 9 | **DDD Этап 2:** перенос интерфейсов в domain | Этап 1 |
| 10 | **DDD Этап 3:** маппинг в репозиториях | Этап 2 |
| 11 | **DDD Этап 4:** доменные сервисы | Этап 3 |
| 12 | **DDD Этап 5:** обновление DTO-маперов | Этап 4 |
| 13 | Мониторинг и structured logging | — |
| 14 | Парсинг вакансий | DDD-сервисы |
