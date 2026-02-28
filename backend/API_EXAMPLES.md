# Примеры API запросов

## Health Check
```bash
curl http://localhost:8081/health
```

## Компании

### Получить все компании
```bash
curl http://localhost:8081/api/v1/companies
```

### Создать компанию
```bash
curl -X POST http://localhost:8081/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{"name": "Тинькофф", "description": "Финтех компания"}'
```

## Вакансии

### Получить все вакансии
```bash
curl http://localhost:8081/api/v1/jobs
```

### Создать вакансию
```bash
curl -X POST http://localhost:8081/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": 1,
    "title": "Senior Go Developer",
    "level": "Senior",
    "specialization": "Golang",
    "salary_min": 350000,
    "salary_max": 550000,
    "salary_currency": "RUB",
    "remote_available": true,
    "is_active": true
  }'
```

## Статистика

### Топ навыков
```bash
curl "http://localhost:8081/api/v1/stats/top-skills?limit=10"
```

### Зарплаты по навыкам
```bash
curl "http://localhost:8081/api/v1/stats/skill-salaries?min_vacancies=2"
```

### Статистика по компаниям
```bash
curl http://localhost:8081/api/v1/stats/companies
```

### Статистика по базам данных
```bash
curl http://localhost:8081/api/v1/stats/databases
```

### Статистика по языкам программирования
```bash
curl http://localhost:8081/api/v1/stats/languages
```
