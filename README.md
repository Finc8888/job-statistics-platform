# Job Statistics Platform

Платформа для управления вакансиями и анализа статистики рынка труда IT-разработчиков.

## Быстрый старт

```bash
make setup          # Первый запуск: Docker stack + schema + seed data
```

## Команды

```bash
make up            # Запустить (данные сохраняются)
make down          # Остановить (данные сохраняются)
make clean         # Остановить + удалить данные

make migrate       # Применить схему БД
make seed          # Перезалить тестовые данные (DESTRUCTIVE)

make logs          # Логи всех контейнеров
make ps            # Статус контейнеров

make dev-api       # Локально: запустить Go API
make dev-frontend  # Локально: собрать и запустить Frontend

make test          # Unit-тесты backend (без БД)
make test-coverage # Тесты + HTML-отчёт покрытия
make test-docker   # Тесты в Docker-контейнере
```

---

## Порты

| Сервис   | Порт |
|----------|------|
| Frontend | 3000 |
| API      | 8081 |
| MySQL    | 3307 |

---

## Структура

```
job-statistics-platform/
├── backend/
│   ├── cmd/api/main.go          # точка входа
│   ├── internal/
│   │   ├── dto/                 # Data Transfer Objects (API ↔ модели)
│   │   ├── handlers/            # HTTP-хендлеры
│   │   ├── repository/          # слой данных + интерфейсы
│   │   ├── models/              # внутренние типы (sql.Null*)
│   │   └── database/            # подключение к БД
│   ├── migrations/              # SQL-миграции + seed
│   ├── TESTING.md               # документация по тестам
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── pages/               # Dashboard, Companies, Jobs, Skills, Statistics
│   │   ├── stores/              # MobX
│   │   ├── services/api.ts      # Axios-клиент
│   │   └── types/               # TypeScript типы
│   └── Dockerfile
├── CLAUDE.md                    # Руководство для Claude Code
├── UPGRADE_PLANS.md             # Планы развития + DDD roadmap
└── docker-compose.full.yml
```

---

## Архитектура backend

```
HTTP Request → Handler → DTO (JSON ↔ Go) → Model (sql.Null*) → Repository → MySQL
```

- **DTO** (`internal/dto/`) — разделяет API-контракты и внутренние модели. `*Request` для входящих данных, `*Response` для ответов
- **Models** (`internal/models/`) — внутренние структуры с `sql.Null*` для nullable колонок. Не сериализуются в JSON напрямую
- **Repository** — raw SQL, интерфейсы для DI и мок-тестирования
- **Handlers** — принимают `dto.*Request`, отвечают `dto.*Response`

---

## API

Base URL: `http://localhost:8081/api/v1`

| Метод | Путь | Описание |
|-------|------|----------|
| GET/POST | `/companies` | список / создать |
| GET/PUT/DELETE | `/companies/:id` | получить / обновить / удалить |
| GET/POST | `/jobs` | вакансии |
| GET/PUT/DELETE | `/jobs/:id` | — |
| GET/POST | `/skills` | навыки |
| GET | `/stats/top-skills?limit=10` | топ навыков |
| GET | `/stats/skill-salaries` | зарплаты по навыкам |
| GET | `/stats/companies` | статистика компаний |
| GET | `/stats/languages` | статистика языков |
| GET | `/health` | health check |

---

## Troubleshooting

**БД не инициализировалась:**
```bash
make clean && make up
```

**Порт занят:**
```bash
lsof -i :3000   # или :8081, :3307
```

**Пересобрать образы:**
```bash
docker-compose -f docker-compose.full.yml up -d --build
```
