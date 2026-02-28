#!/bin/bash

# seed.sh — загрузка тестовых данных.
# ⚠️  ДЕСТРУКТИВНО: очищает все таблицы и заново заполняет тестовыми данными.
# Все пользовательские изменения будут потеряны.

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MYSQL_CONTAINER="job_stats_mysql"

echo -e "${YELLOW}⚠️  Загрузка тестовых данных очистит все таблицы!${NC}"

if ! docker ps --filter "name=${MYSQL_CONTAINER}" --filter "status=running" -q | grep -q .; then
    echo -e "${RED}❌ MySQL контейнер не запущен!${NC}"
    echo "Запустите: make up"
    exit 1
fi

docker exec -i "${MYSQL_CONTAINER}" mysql --default-character-set=utf8mb4 -u jobuser -pjobpassword job_stats \
    < "$SCRIPT_DIR/002_seed_data.sql"

echo -e "${GREEN}✅ Тестовые данные загружены${NC}"
echo ""
echo "Компании:"
docker exec "${MYSQL_CONTAINER}" mysql --default-character-set=utf8mb4 -u jobuser -pjobpassword job_stats \
    -e "SELECT id, name FROM companies;"
echo ""
echo "Вакансии:"
docker exec "${MYSQL_CONTAINER}" mysql --default-character-set=utf8mb4 -u jobuser -pjobpassword job_stats \
    -e "SELECT id, title, level FROM jobs;"
