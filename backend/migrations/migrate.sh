#!/bin/bash

# migrate.sh — только схема БД (CREATE TABLE IF NOT EXISTS)
# Безопасно запускать повторно: не трогает данные пользователя.

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MYSQL_CONTAINER="job_stats_mysql"

echo "🔧 Применение схемы БД..."

if ! docker ps --filter "name=${MYSQL_CONTAINER}" --filter "status=running" -q | grep -q .; then
    echo -e "${RED}❌ MySQL контейнер не запущен!${NC}"
    echo "Запустите: make up"
    exit 1
fi

echo "Ожидание готовности MySQL..."
sleep 5

docker exec -i "${MYSQL_CONTAINER}" mysql --default-character-set=utf8mb4 -u jobuser -pjobpassword job_stats \
    < "$SCRIPT_DIR/001_create_tables.sql"

echo -e "${GREEN}✅ Схема применена (таблицы созданы / уже существуют)${NC}"
