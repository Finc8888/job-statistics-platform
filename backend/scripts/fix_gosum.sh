#!/bin/bash

# 1. Удалите текущий go.sum
rm go.sum

# 2. Скачайте зависимости и создайте go.sum
go mod download
go mod tidy

# 3. Проверьте, что go.sum создался и не пустой
cat go.sum | head -20
