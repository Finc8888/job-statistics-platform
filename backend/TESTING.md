# Тестирование Backend

## Обзор

Тесты покрывают два слоя приложения:

| Слой | Файлы | Подход |
|------|-------|--------|
| Repository | `internal/repository/*_test.go` | Табличные тесты + `go-sqlmock` (без реальной БД) |
| Handlers | `internal/handlers/*_handler_test.go` | `net/http/httptest` + mock-репозитории |

**Тестов нет зависимости от реальной базы данных** — используется `go-sqlmock` для имитации SQL-запросов и mock-структуры для handlers.

---

## Структура тестовых файлов

```
backend/
├── internal/
│   ├── repository/
│   │   ├── interfaces.go                  # Интерфейсы для DI и тестирования
│   │   ├── company_repository_test.go     # Тесты CompanyRepository (CRUD)
│   │   └── job_repository_test.go         # Тесты JobRepository (CRUD)
│   └── handlers/
│       ├── company_handler_test.go        # Тесты HTTP-хендлеров компаний
│       └── job_handler_test.go            # Тесты HTTP-хендлеров вакансий
├── Dockerfile.test                        # Docker-образ для тестов
├── docker-compose.test.yml                # Compose для запуска тестов
└── coverage/                              # Отчёты покрытия (генерируется)
```

---

## Способы запуска

### 1. Локально (требуется Go 1.21+)

```bash
# Запустить все тесты
make test

# Или напрямую через go test
go test ./... -v -count=1
```

### 2. Только определённый пакет

```bash
# Тесты репозитория
go test ./internal/repository/... -v

# Тесты хендлеров
go test ./internal/handlers/... -v
```

### 3. Один тест или группа тестов

```bash
# Запустить конкретный тест-кейс
go test ./internal/repository/... -run TestCompanyRepository_GetAll -v

# Запустить все тесты компании
go test ./internal/repository/... -run TestCompanyRepository -v

# Запустить все тесты хендлеров вакансий
go test ./internal/handlers/... -run TestJobHandler -v
```

### 4. С отчётом о покрытии (локально)

```bash
make test-coverage
```

Команда:
1. Запускает тесты
2. Выводит покрытие по функциям в консоль
3. Генерирует `coverage.html` — откройте в браузере для детального отчёта

### 5. В Docker-контейнере (рекомендуется для CI)

```bash
make test-docker
```

Или вручную:

```bash
mkdir -p coverage
docker-compose -f docker-compose.test.yml run --rm backend-test
```

После выполнения файл `coverage/coverage.out` будет доступен на хост-машине.

---

## Как работают тесты

### Repository-слой (sqlmock)

Тесты репозитория используют `github.com/DATA-DOG/go-sqlmock` — библиотеку, которая подменяет реальное `*sql.DB` объектом-заглушкой. Это позволяет:
- Проверять что репозиторий формирует **правильные SQL-запросы**
- Проверять обработку **ошибок базы данных**
- Тестировать **без запущенного MySQL**

Каждый тест-кейс:
1. Создаёт новый mock (`sqlmock.New()`)
2. Настраивает ожидаемые запросы (`ExpectQuery`, `ExpectExec`)
3. Вызывает метод репозитория
4. Проверяет результат и выполнение всех ожиданий (`ExpectationsWereMet`)

### Handler-слой (httptest + mock)

Тесты хендлеров используют:
- `net/http/httptest` — записывает HTTP-ответы без запуска сервера
- `mockJobRepo` / `mockCompanyRepo` — структуры, реализующие интерфейсы репозитория

Каждый тест-кейс:
1. Создаёт mock-репозиторий с нужным поведением (`err` или тестовые данные)
2. Формирует HTTP-запрос через `httptest.NewRequest`
3. Записывает ответ в `httptest.NewRecorder`
4. Проверяет HTTP-статус и тело ответа

---

## Табличные тесты (table-driven tests)

Все тесты написаны в стиле [table-driven tests](https://go.dev/wiki/TableDrivenTests) — стандартный идиоматичный подход в Go:

```go
tests := []struct {
    name       string
    // входные данные
    // ожидаемый результат
}{
    {"успешный случай", ...},
    {"ошибка БД", ...},
    {"не найдено", ...},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // логика теста
    })
}
```

---

## Добавление новых тестов

### Тест для нового метода репозитория

```go
func TestMyRepository_NewMethod(t *testing.T) {
    tests := []struct {
        name      string
        mockSetup func(mock sqlmock.Sqlmock)
        wantErr   bool
    }{
        {
            name: "успешный случай",
            mockSetup: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery("SELECT ...").WillReturnRows(...)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db, mock, _ := sqlmock.New()
            defer db.Close()

            tt.mockSetup(mock)
            repo := NewMyRepository(db)

            _, err := repo.NewMethod()
            if (err != nil) != tt.wantErr {
                t.Errorf("...")
            }
            mock.ExpectationsWereMet()
        })
    }
}
```

### Тест для нового хендлера

```go
// 1. Добавить метод в mockJobRepo / создать новый mock
// 2. Написать тест:
func TestMyHandler_NewEndpoint(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/api/v1/my-endpoint", nil)
    rec := httptest.NewRecorder()

    handler := NewMyHandler(&mockMyRepo{})
    handler.NewEndpoint(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("status = %d, want 200", rec.Code)
    }
}
```

---

## Частые проблемы

**`unfulfilled mock expectations`**
SQL-запрос в коде не совпадает с ожидаемым в тесте. Проверьте точность строки запроса (включая пробелы).

**`cannot use ... as driver.Value`**
В sqlmock v1 метод `AddRow` принимает `...driver.Value`. Убедитесь, что хелперы возвращают `[]driver.Value`, а числа приводятся к `int64`.

**Тест прошёл, но проверка не сработала**
Убедитесь что вы вызываете `t.Fatal` / `t.Error` а не просто логируете — тест без явного фейла считается прошедшим.
