package main

import (
	"fmt"
	"job-statistics-api/internal/database"
	"job-statistics-api/internal/handlers"
	"job-statistics-api/internal/repository"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения системы")
	}

	// Подключаемся к БД
	if err := database.Connect(); err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer database.Close()

	// Инициализируем репозитории
	companyRepo := repository.NewCompanyRepository(database.DB)
	jobRepo := repository.NewJobRepository(database.DB)
	skillRepo := repository.NewSkillRepository(database.DB)
	locationRepo := repository.NewLocationRepository(database.DB)
	statsRepo := repository.NewStatsRepository(database.DB)

	// Инициализируем обработчики
	companyHandler := handlers.NewCompanyHandler(companyRepo)
	jobHandler := handlers.NewJobHandler(jobRepo)
	skillHandler := handlers.NewSkillHandler(skillRepo)
	locationHandler := handlers.NewLocationHandler(locationRepo)
	statsHandler := handlers.NewStatsHandler(statsRepo)

	// Создаем роутер
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Companies endpoints
	api.HandleFunc("/companies", companyHandler.GetAll).Methods("GET")
	api.HandleFunc("/companies/{id}", companyHandler.GetByID).Methods("GET")
	api.HandleFunc("/companies", companyHandler.Create).Methods("POST")
	api.HandleFunc("/companies/{id}", companyHandler.Update).Methods("PUT")
	api.HandleFunc("/companies/{id}", companyHandler.Delete).Methods("DELETE")

	// Jobs endpoints
	api.HandleFunc("/jobs", jobHandler.GetAll).Methods("GET")
	api.HandleFunc("/jobs/{id}", jobHandler.GetByID).Methods("GET")
	api.HandleFunc("/jobs", jobHandler.Create).Methods("POST")
	api.HandleFunc("/jobs/{id}", jobHandler.Update).Methods("PUT")
	api.HandleFunc("/jobs/{id}", jobHandler.Delete).Methods("DELETE")

	// Skills endpoints
	api.HandleFunc("/skills", skillHandler.GetAll).Methods("GET")
	api.HandleFunc("/skills/{id}", skillHandler.GetByID).Methods("GET")
	api.HandleFunc("/skills", skillHandler.Create).Methods("POST")
	api.HandleFunc("/skills/{id}", skillHandler.Update).Methods("PUT")
	api.HandleFunc("/skills/{id}", skillHandler.Delete).Methods("DELETE")

	// Locations endpoints
	api.HandleFunc("/locations", locationHandler.GetAll).Methods("GET")
	api.HandleFunc("/locations/job/{job_id}", locationHandler.GetByJobID).Methods("GET")
	api.HandleFunc("/locations", locationHandler.Create).Methods("POST")
	api.HandleFunc("/locations/{id}", locationHandler.Update).Methods("PUT")
	api.HandleFunc("/locations/{id}", locationHandler.Delete).Methods("DELETE")

	// Statistics endpoints
	api.HandleFunc("/stats/top-skills", statsHandler.GetTopSkills).Methods("GET")
	api.HandleFunc("/stats/skill-salaries", statsHandler.GetSkillSalaries).Methods("GET")
	api.HandleFunc("/stats/skills-by-level", statsHandler.GetSkillsByLevel).Methods("GET")
	api.HandleFunc("/stats/companies", statsHandler.GetCompanyStats).Methods("GET")
	api.HandleFunc("/stats/databases", statsHandler.GetDatabaseStats).Methods("GET")
	api.HandleFunc("/stats/languages", statsHandler.GetLanguageStats).Methods("GET")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Получаем порт
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8081"
	}

	// Настраиваем сервер.
	// corsMiddleware оборачивает весь роутер — это гарантирует что
	// preflight OPTIONS-запросы получают CORS-заголовки до роутинга.
	srv := &http.Server{
		Handler:      corsMiddleware(r),
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер
	log.Printf("🚀 Сервер запущен на порту %s", port)
	log.Printf("📊 API доступен по адресу: http://localhost:%s/api/v1", port)
	log.Fatal(srv.ListenAndServe())
}

// CORS middleware для разрешения запросов с других доменов
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого origin
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Обрабатываем preflight запросы
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
