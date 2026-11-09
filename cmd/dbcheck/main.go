package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5/pgxpool"

	// Замени путь импорта, если в go.mod другое имя модуля.
	"github.com/Zuf4ik/Z_JobScout/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не загружен, используются переменные окружения")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	pool, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Ошибка подключения: ", err)
	}
	defer pool.Close()

	checkDatabase(ctx, pool)
}

func checkDatabase(ctx context.Context, pool *pgxpool.Pool) {
	var name string

	err := pool.QueryRow(
		ctx,
		"SELECT current_database()",
	).Scan(&name)

	if err != nil {
		log.Fatal("Ошибка SQL-запроса: ", err)
	}

	log.Printf("PostgreSQL подключён. База данных: %s", name)
}
