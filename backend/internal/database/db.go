package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB - глобальное соединение с БД
var DB *sql.DB

// Connect устанавливает соединение с базой данных
func Connect() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Формируем DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("ошибка открытия соединения: %v", err)
	}

	// Настройки пула соединений
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем соединение
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ошибка ping базы данных: %v", err)
	}

	log.Println("✅ Успешное подключение к базе данных MySQL")
	return nil
}

// Close закрывает соединение с БД
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
