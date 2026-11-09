package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

func TestFavoriteRepositoryIntegration(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("Используем DATABASE_URL из окружения")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL не задан")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	if err := pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	var userID int64

	err = pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash)
		VALUES (
			'test@jobscout.local',
			'test-only-not-a-real-password-hash'
		)
		ON CONFLICT (email)
		DO UPDATE SET email = EXCLUDED.email
		RETURNING id
	`).Scan(&userID)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewFavoriteRepository(pool)

	favorite := vacancy.Favorite{
		UserID:      userID,
		VacancyID:   "integration-test-001",
		Source:      "test",
		Title:       "Junior Go Developer",
		CompanyName: "Test Company",
		URL:         "https://example.com/job",
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cleanupCancel()

		if _, err := pool.Exec(cleanupCtx, `
			DELETE FROM favorite_vacancies
			WHERE user_id = $1 AND source = 'test'
		`, userID); err != nil {
			t.Logf("Ошибка очистки тестовых данных: %v", err)
		}
	})

	if err := repo.Add(ctx, favorite); err != nil {
		t.Fatalf("Add: %v", err)
	}

	items, err := repo.GetAll(ctx, userID)
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}

	found := false
	for _, item := range items {
		if item.Source == favorite.Source &&
			item.VacancyID == favorite.VacancyID {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("Добавленная вакансия не найдена")
	}

	if err := repo.Delete(
		ctx,
		userID,
		favorite.Source,
		favorite.VacancyID,
	); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	items, err = repo.GetAll(ctx, userID)
	if err != nil {
		t.Fatalf("GetAll после удаления: %v", err)
	}

	for _, item := range items {
		if item.Source == favorite.Source &&
			item.VacancyID == favorite.VacancyID {
			t.Fatal("Вакансия осталась после удаления")
		}
	}
}
