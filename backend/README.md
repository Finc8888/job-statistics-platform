# Job Statistics API

REST API для управления вакансиями и анализа статистики рынка труда IT-разработчиков.

## 🚀 Возможности

- CRUD операции для компаний, вакансий, навыков и локаций
- Статистика по самым востребованным навыкам
- Анализ зарплат по навыкам
- Статистика по языкам программирования и базам данных
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
GET /api/v1/stats/skill-salaries?min_vacancies=2
```
Возвращает среднюю зарплату по навыкам (с минимальным количеством вакансий).

#### Навыки по уровню
```bash
GET /api/v1/stats/skills-by-level
```
Возвращает востребованность навыков по уровню (Junior/Middle/Senior).

#### Статистика по компаниям
```bash
GET /api/v1/stats/companies
```
Возвращает статистику по компаниям (количество вакансий, зарплаты, локации).

#### Статистика по базам данных
```bash
GET /api/v1/stats/databases
```
Возвращает статистику по базам данных.

#### Статистика по языкам программирования
```bash
GET /api/v1/stats/languages
```
Возвращает статистику по языкам программирования.

## 🗄️ Структура базы данных

```
companies
├── id
├── name
├── description
├── created_at
└── updated_at

jobs
├── id
├── company_id (FK)
├── title
├── level
├── specialization
├── salary_min
├── salary_max
├── salary_currency
├── experience_years
├── location
├── remote_available
├── description
├── responsibilities
├── benefits
├── posted_date
├── is_active
├── source_url
├── created_at
└── updated_at

skills
├── id
├── name
├── category
└── created_at

job_skills
├── id
├── job_id (FK)
├── skill_id (FK)
├── is_required
├── is_nice_to_have
└── created_at

locations
├── id
├── job_id (FK)
├── city
├── metro_station
└── is_primary
```

## 🔧 Переменные окружения

```env
DB_HOST=mysql
DB_PORT=3306
DB_USER=jobuser
DB_PASSWORD=jobpassword
DB_NAME=job_stats
API_PORT=8081
```

## 📊 Примеры использования

### Получить топ-10 навыков
```bash
curl http://localhost:8081/api/v1/stats/top-skills?limit=10
```

### Получить статистику по компаниям
```bash
curl http://localhost:8081/api/v1/stats/companies
```

### Получить все активные вакансии
```bash
curl http://localhost:8081/api/v1/jobs
```

## 🐳 Docker команды

```bash
# Запустить контейнеры
docker-compose up -d

# Остановить контейнеры
docker-compose down

# Пересобрать образы
docker-compose build

# Просмотр логов
docker-compose logs -f

# Выполнить команду в контейнере
docker-compose exec api sh
docker-compose exec mysql mysql -u jobuser -p job_stats
```

## 🧪 Тестирование API

Вы можете использовать Postman, Insomnia или curl для тестирования API.

Примеры curl-запросов находятся выше в разделе API Endpoints.

## 📝 Лицензия

MIT
