.PHONY: help up down clean setup migrate seed rebuild rebuild-frontend rebuild-api logs logs-api ps test test-coverage test-docker dev-local dev-api dev-frontend

help: ## Показать справку
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ── Docker-режим ─────────────────────────────────────────────────────────────

up: ## Запустить стек (данные сохраняются между перезапусками)
	docker-compose -f docker-compose.full.yml up -d --build
	@echo ""
	@echo "✅ Стек запущен"
	@echo "   Frontend : http://localhost:3000"
	@echo "   API      : http://localhost:8081"

down: ## Остановить контейнеры (данные сохраняются)
	docker-compose -f docker-compose.full.yml down
	@echo "Контейнеры остановлены. Данные в БД сохранены."

clean: ## ⚠️  Остановить контейнеры и удалить volumes (все данные будут удалены)
	docker-compose -f docker-compose.full.yml down -v
	@echo "Контейнеры и volumes удалены."

setup: ## Первый запуск: поднять стек + создать схему + загрузить тестовые данные
	docker-compose -f docker-compose.full.yml up -d --build
	@echo "Ожидание запуска MySQL (~15 сек)..."
	@sleep 15
	@$(MAKE) migrate
	@$(MAKE) seed
	@echo ""
	@echo "✅ Готово! Первичная настройка завершена."
	@echo "   Frontend : http://localhost:3000"
	@echo "   API      : http://localhost:8081"

migrate: ## Применить схему БД (CREATE TABLE IF NOT EXISTS, данные не трогает)
	cd backend/migrations && chmod +x migrate.sh && ./migrate.sh

seed: ## ⚠️  Загрузить тестовые данные (очищает все таблицы!)
	@echo "Загрузка тестовых данных (все текущие данные будут удалены)..."
	cd backend/migrations && chmod +x seed.sh && ./seed.sh

rebuild-frontend: ## Пересобрать frontend без кэша (если make up не подхватил изменения)
	docker-compose -f docker-compose.full.yml build --no-cache frontend
	docker-compose -f docker-compose.full.yml up -d frontend
	@echo "✅ Frontend пересобран и перезапущен"

rebuild-api: ## Пересобрать API без кэша (если make up не подхватил изменения)
	docker-compose -f docker-compose.full.yml build --no-cache api
	docker-compose -f docker-compose.full.yml up -d api
	@echo "✅ API пересобран и перезапущен"

rebuild: ## Пересобрать frontend + API без кэша
	docker-compose -f docker-compose.full.yml build --no-cache frontend api
	docker-compose -f docker-compose.full.yml up -d frontend api
	@echo "✅ Стек пересобран и перезапущен"

ps: ## Статус контейнеров
	docker-compose -f docker-compose.full.yml ps

logs: ## Логи всех контейнеров
	docker-compose -f docker-compose.full.yml logs -f

logs-api: ## Логи только API
	docker-compose -f docker-compose.full.yml logs -f api

# ── Локальная разработка ─────────────────────────────────────────────────────

dev-local: ## Запустить MySQL в Docker, API и Frontend — локально (2 терминала)
	@echo "Запуск MySQL в Docker..."
	docker-compose -f docker-compose.full.yml up -d mysql
	@echo ""
	@echo "Через ~10 сек выполни:"
	@echo "   make migrate    (создать схему)"
	@echo ""
	@echo "Затем в двух терминалах:"
	@echo "   make dev-api       # Терминал 1"
	@echo "   make dev-frontend  # Терминал 2"

dev-api: ## Запустить Go API локально
	cd backend && go run cmd/api/main.go

dev-frontend: ## Запустить Frontend локально (сборка + сервер)
	cd frontend && yarn install && yarn build && yarn preview

# ── Тесты ────────────────────────────────────────────────────────────────────

test: ## Запустить unit-тесты backend
	cd backend && go test ./... -v -count=1

test-coverage: ## Тесты с HTML-отчётом покрытия
	cd backend && go test ./... -coverprofile=coverage.out -count=1
	cd backend && go tool cover -func=coverage.out
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Отчёт: backend/coverage.html"

test-docker: ## Запустить тесты в изолированном Docker-контейнере
	mkdir -p backend/coverage
	docker-compose -f backend/docker-compose.test.yml run --rm backend-test
