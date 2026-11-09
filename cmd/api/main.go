package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/Zuf4ik/Z_JobScout/internal/database"
	"github.com/Zuf4ik/Z_JobScout/internal/httpapi"
	"github.com/Zuf4ik/Z_JobScout/internal/repository"
	"github.com/Zuf4ik/Z_JobScout/internal/service"
)

func main() {
	// Загружаем переменные из .env.
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения")
	}

	// Подключаем PostgreSQL.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	db, err := database.Connect(
		ctx,
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL: ", err)
	}
	defer db.Close()

	log.Println("PostgreSQL подключён")

	// Существующий поиск вакансий.
	searchService := service.NewSearchService()
	searchHandler := httpapi.NewHandler(searchService)

	// Репозитории пользователей и сессий.
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Бизнес-логика авторизации.
	authService := service.NewAuthService(userRepo)
	sessionService := service.NewSessionService(sessionRepo)

	// HTTP-обработчики авторизации.
	// Secure=false допускается только для локального HTTP.
	authHandler := httpapi.NewAuthHandler(
		authService,
		sessionService,
		false,
	)

	// Общий маршрутизатор.
	mux := http.NewServeMux()

	// Регистрация и вход.
	mux.HandleFunc(
		"POST /api/v1/auth/register",
		authHandler.Register,
	)
	mux.HandleFunc(
		"POST /api/v1/auth/login",
		authHandler.Login,
	)

	mux.HandleFunc("GET /api/v1/auth/me", authHandler.Me)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)

	// Все остальные запросы передаём существующему API.
	mux.Handle("/", searchHandler)

	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("JobScout API запущен: http://127.0.0.1:8080")

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Ошибка HTTP-сервера: ", err)
	}
}
