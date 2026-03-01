# Job Statistics Platform

Платформа для управления вакансиями и анализа статистики рынка труда IT-разработчиков.

## Быстрый старт

```bash
make setup          # Первый запуск: Docker stack + schema + seed data
```

## Команды

### Основные

```bash
make up            # Запустить (данные сохраняются)
make down          # Остановить (данные сохраняются)
make clean         # Остановить + удалить данные
```

### Пересборка после изменения кода

> **`make up` может не подхватить изменения из-за Docker кэша.**
> После любых правок в коде используйте `make rebuild`:

```bash
make rebuild            # Пересобрать frontend + API без кэша
make rebuild-frontend   # Только frontend (после правок в frontend/src/)
make rebuild-api        # Только API (после правок в backend/)
```

### База данных

```bash
make migrate       # Применить схему БД
make seed          # Перезалить тестовые данные (DESTRUCTIVE)
```

### Тесты и мониторинг

```bash
make test          # Unit-тесты backend (без БД)
make test-coverage # Тесты + HTML-отчёт покрытия
make test-docker   # Тесты в Docker-контейнере

make ps            # Статус контейнеров
make logs          # Логи всех контейнеров
make logs-api      # Логи только API
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
| GET/POST | `/jobs/:id/skills` | навыки вакансии / привязать навыки |
| GET/POST | `/skills` | навыки |
| GET | `/stats/top-skills?limit=10` | топ навыков |
| GET | `/stats/skill-salaries` | зарплаты по навыкам |
| GET | `/stats/companies` | статистика компаний |
| GET | `/stats/languages` | статистика языков |
| GET | `/health` | health check |

---

## Troubleshooting

**Изменения в коде не применяются:**
```bash
make rebuild          # пересборка без Docker-кэша
```

**БД не инициализировалась:**
```bash
make clean && make setup
```

**Порт занят:**
```bash
lsof -i :3000   # или :8081, :3307
```

**Пересобрать образы:**
```bash
make rebuild
```
