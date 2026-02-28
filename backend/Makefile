.PHONY: help build run stop clean logs test test-docker test-coverage

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Собрать Docker образы
	docker-compose build

up: ## Запустить приложение
	docker-compose up -d
	@echo "✅ Приложение запущено на http://localhost:8081"

run: up ## Алиас для up

stop: ## Остановить контейнеры
	docker-compose stop

down: ## Остановить и удалить контейнеры
	docker-compose down

clean: ## Остановить контейнеры и удалить volumes
	docker-compose down -v

logs: ## Показать логи
	docker-compose logs -f

logs-api: ## Показать логи API
	docker-compose logs -f api

logs-db: ## Показать логи MySQL
	docker-compose logs -f mysql

restart: ## Перезапустить приложение
	docker-compose restart

ps: ## Показать статус контейнеров
	docker-compose ps

shell-api: ## Войти в контейнер API
	docker-compose exec api sh

shell-db: ## Войти в MySQL
	docker-compose exec mysql mysql -u jobuser -p job_stats

test-api: ## Тестовый запрос к API
	@echo "Testing API health..."
	@curl -s http://localhost:8081/health
	@echo "\n\nTesting companies endpoint..."
	@curl -s http://localhost:8081/api/v1/companies | head -c 200

test: ## Запустить unit-тесты локально
	go test ./... -v -count=1

test-coverage: ## Запустить тесты с отчётом о покрытии (открывает HTML)
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Отчёт сохранён: coverage.html"

test-docker: ## Запустить тесты в изолированном Docker-контейнере
	@mkdir -p coverage
	docker-compose -f docker-compose.test.yml run --rm backend-test
	@echo "✅ Тесты завершены. Coverage: coverage/coverage.out"

dev: ## Запустить в режиме разработки
	go run cmd/api/main.go

install: ## Установить зависимости Go
	go mod download
	go mod tidy

migrate: ## Выполнить миграции
	cd migrations && chmod +x migrate.sh && ./migrate.sh

db-reset: ## Пересоздать БД с нуля
	docker-compose down -v
	docker-compose up -d mysql
	@echo "Ожидание запуска MySQL..."
	@sleep 10
	cd migrations && chmod +x migrate.sh && ./migrate.sh
	docker-compose up -d api
	@echo "✅ База данных пересоздана!"
