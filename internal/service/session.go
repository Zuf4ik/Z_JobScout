package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zuf4ik/Z_JobScout/internal/user"
)

const SessionTTL = 7 * 24 * time.Hour

var ErrInvalidSession = errors.New("недействительная сессия")

// SessionStore описывает операции хранилища сессий.
type SessionStore interface {
	Create(
		ctx context.Context,
		userID int64,
		tokenHash [32]byte,
		expiresAt time.Time,
	) (*user.Session, error)

	GetByToken(
		ctx context.Context,
		tokenHash [32]byte,
	) (*user.Session, error)

	Delete(
		ctx context.Context,
		tokenHash [32]byte,
	) error
}

type SessionService struct {
	store SessionStore
}

func NewSessionService(store SessionStore) *SessionService {
	return &SessionService{store: store}
}

// Create создаёт новую сессию и возвращает токен для cookie.
func (s *SessionService) Create(
	ctx context.Context,
	userID int64,
) (string, *user.Session, error) {
	if userID <= 0 {
		return "", nil, errors.New("некорректный ID пользователя")
	}

	token, tokenHash, err := GenerateSessionToken()
	if err != nil {
		return "", nil, fmt.Errorf("генерация токена: %w", err)
	}

	expiresAt := time.Now().UTC().Add(SessionTTL)

	session, err := s.store.Create(
		ctx,
		userID,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		return "", nil, fmt.Errorf("сохранение сессии: %w", err)
	}

	return token, session, nil
}

// Authenticate находит действующую сессию по токену.
func (s *SessionService) Authenticate(
	ctx context.Context,
	token string,
) (*user.Session, error) {
	hash, err := HashSessionToken(token)
	if err != nil {
		return nil, ErrInvalidSession
	}

	session, err := s.store.GetByToken(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("проверка сессии: %w", err)
	}

	if !session.ExpiresAt.After(time.Now()) {
		return nil, ErrInvalidSession
	}

	return session, nil
}

// Revoke завершает сессию.
func (s *SessionService) Revoke(
	ctx context.Context,
	token string,
) error {
	hash, err := HashSessionToken(token)
	if err != nil {
		return ErrInvalidSession
	}

	if err := s.store.Delete(ctx, hash); err != nil {
		return fmt.Errorf("удаление сессии: %w", err)
	}

	return nil
}
