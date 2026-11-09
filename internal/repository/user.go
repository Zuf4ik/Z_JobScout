package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Zuf4ik/Z_JobScout/internal/user"
)

var (
	ErrEmailTaken   = errors.New("email уже зарегистрирован")
	ErrUserNotFound = errors.New("пользователь не найден")
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	email string,
	passwordHash string,
) (*user.User, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at
	`

	var u user.User

	err := r.db.QueryRow(
		ctx,
		query,
		email,
		passwordHash,
	).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailTaken
		}

		return nil, fmt.Errorf(
			"ошибка создания пользователя: %w",
			err,
		)
	}

	return &u, nil
}

func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`

	var u user.User

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"ошибка поиска пользователя: %w",
			err,
		)
	}

	return &u, nil
}
