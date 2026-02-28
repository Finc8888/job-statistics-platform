# Job Statistics Platform

Платформа для анализа рынка труда IT: CRUD вакансий/компаний/навыков + интерактивные графики статистики.

**Стек:** Go 1.21 · MySQL 8 · React 18 · TypeScript · MobX · Recharts · Docker

---

## Быстрый старт

### Первый запуск

```bash
make setup   # поднять стек + создать схему + загрузить тестовые данные
```

### Обычный перезапуск (данные сохраняются)

```bash
make down    # остановить (данные в БД остаются)
make up      # запустить снова — все ваши изменения на месте
```

> `make down` не удаляет volume с данными MySQL. Данные теряются только при `make clean`.

### Локальная разработка

Требуется: Go 1.21+, Node.js 18+, Yarn 4, Docker (только для MySQL).

```bash
make dev-local     # запустить MySQL в Docker + инструкция
make migrate       # создать схему (после запуска MySQL)
make dev-api       # Терминал 1 — Go API на :8081
make dev-frontend  # Терминал 2 — Frontend на :3000
```

---

## Все Make-команды

```bash
make setup         # Первый запуск: стек + схема + тестовые данные
make up            # Запустить стек (данные сохраняются)
make down          # Остановить (данные сохраняются)
make clean         # ⚠️  Остановить + удалить volumes (все данные удалятся)

make migrate       # Создать/обновить схему БД (безопасно, данные не трогает)
make seed          # ⚠️  Загрузить тестовые данные (очищает таблицы!)

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
│   │   ├── handlers/            # HTTP-хендлеры
│   │   ├── repository/          # слой данных + интерфейсы
│   │   └── models/              # типы
│   ├── migrations/              # SQL-миграции + seed
│   ├── TESTING.md               # документация по тестам
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── pages/               # Dashboard, Companies, Jobs, Skills, Statistics
│   │   ├── stores/              # MobX
│   │   └── services/api.ts      # Axios-клиент
│   └── Dockerfile
└── docker-compose.full.yml
```

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
