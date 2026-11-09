package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Zuf4ik/Z_JobScout/internal/user"
)

var ErrSessionNotFound = errors.New("сессия не найдена или истекла")

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(
	ctx context.Context,
	userID int64,
	tokenHash [32]byte,
	expiresAt time.Time,
) (*user.Session, error) {
	const query = `
		INSERT INTO user_sessions (
			user_id,
			token_hash,
			expires_at
		)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, created_at, expires_at
	`

	var session user.Session

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		tokenHash[:],
		expiresAt,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("создание сессии: %w", err)
	}

	return &session, nil
}

func (r *SessionRepository) GetByToken(
	ctx context.Context,
	tokenHash [32]byte,
) (*user.Session, error) {
	const query = `
		SELECT id, user_id, created_at, expires_at
		FROM user_sessions
		WHERE token_hash = $1
		  AND expires_at > NOW()
	`

	var session user.Session

	err := r.db.QueryRow(
		ctx,
		query,
		tokenHash[:],
	).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("поиск сессии: %w", err)
	}

	return &session, nil
}

func (r *SessionRepository) Delete(
	ctx context.Context,
	tokenHash [32]byte,
) error {
	const query = `
		DELETE FROM user_sessions
		WHERE token_hash = $1
	`

	_, err := r.db.Exec(ctx, query, tokenHash[:])
	if err != nil {
		return fmt.Errorf("удаление сессии: %w", err)
	}

	return nil
}

func (r *SessionRepository) DeleteExpired(
	ctx context.Context,
) (int64, error) {
	const query = `
		DELETE FROM user_sessions
		WHERE expires_at <= NOW()
	`

	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("очистка просроченных сессий: %w", err)
	}

	return result.RowsAffected(), nil
}
