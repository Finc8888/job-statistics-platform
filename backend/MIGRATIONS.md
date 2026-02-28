# 🔧 Решение проблемы с пустой базой данных

Если вы получили ошибку `Table 'job_stats.skills' doesn't exist`, выполните следующие шаги:

## Способ 1: Полный пересоздание (рекомендуется)

```bash
# 1. Остановите контейнеры и удалите volumes
docker-compose down -v

# 2. Запустите заново
docker-compose up -d

# 3. Подождите 10 секунд для инициализации MySQL
sleep 10

# 4. Проверьте таблицы
docker-compose exec mysql mysql -u jobuser -pjobpassword job_stats -e "SHOW TABLES;"
```

## Способ 2: Выполнить миграции вручную

```bash
# 1. Убедитесь, что MySQL запущен
docker-compose up -d mysql

# 2. Выполните скрипт миграции
cd migrations
chmod +x migrate.sh
./migrate.sh

# 3. Запустите API
docker-compose up -d api

# 4. Проверьте работу
curl http://localhost:8081/api/v1/stats/top-skills?limit=10
```

## Способ 3: Пошаговое выполнение миграций

```bash
# 1. Создание таблиц
docker-compose exec -T mysql mysql -u jobuser -pjobpassword job_stats < migrations/001_create_tables.sql

# 2. Загрузка тестовых данных
docker-compose exec -T mysql mysql -u jobuser -pjobpassword job_stats < migrations/002_seed_data.sql

# 3. Проверка
docker-compose exec mysql mysql -u jobuser -pjobpassword job_stats -e "SELECT COUNT(*) FROM companies;"
docker-compose exec mysql mysql -u jobuser -pjobpassword job_stats -e "SELECT COUNT(*) FROM skills;"
docker-compose exec mysql mysql -u jobuser -pjobpassword job_stats -e "SELECT COUNT(*) FROM jobs;"
```

## Проверка успешной миграции

После выполнения миграций проверьте:

```bash
# Количество записей
docker-compose exec mysql mysql -u jobuser -pjobpassword job_stats << EOF
SELECT
    (SELECT COUNT(*) FROM companies) as companies,
    (SELECT COUNT(*) FROM skills) as skills,
    (SELECT COUNT(*) FROM jobs) as jobs,
    (SELECT COUNT(*) FROM locations) as locations;
EOF
```

Должно быть:
- **companies**: 4
- **skills**: 40+
- **jobs**: 4
- **locations**: 5+

## Если ничего не помогло

Войдите в MySQL и проверьте вручную:

```bash
# Войти в MySQL
docker-compose exec mysql mysql -u jobuser -pjobpassword job_stats

# В MySQL консоли:
SHOW TABLES;
SELECT * FROM companies;
SELECT COUNT(*) FROM skills;
EXIT;
```

## Обновление Makefile

Добавьте в Makefile команду для миграций:

```makefile
migrate: ## Выполнить миграции
	cd migrations && ./migrate.sh

db-reset: ## Пересоздать БД с нуля
	docker-compose down -v
	docker-compose up -d mysql
	sleep 10
	cd migrations && ./migrate.sh
	docker-compose up -d api
```

Использование:

```bash
make migrate      # Выполнить миграции
make db-reset     # Полный сброс БД
```
