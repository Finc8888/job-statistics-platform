#!/bin/bash

# Переместить файлы в правильные директории

# Database
mv db.go internal/database/

# Models
mv models.go internal/models/

# Handlers
mv company_handler.go internal/handlers/
mv job_handler.go internal/handlers/
mv skill_location_handler.go internal/handlers/
mv stats_handler.go internal/handlers/

# Repository
mv company_repository.go internal/repository/
mv job_repository.go internal/repository/
mv skill_repository.go internal/repository/
mv location_repository.go internal/repository/
mv stats_repository.go internal/repository/
