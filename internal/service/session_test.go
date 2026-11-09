package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Zuf4ik/Z_JobScout/internal/repository"
	"github.com/Zuf4ik/Z_JobScout/internal/user"
)

type fakeSessionStore struct {
	sessions map[[32]byte]user.Session
}

func newFakeSessionStore() *fakeSessionStore {
	return &fakeSessionStore{
		sessions: make(map[[32]byte]user.Session),
	}
}

func (f *fakeSessionStore) Create(
	ctx context.Context,
	userID int64,
	hash [32]byte,
	expiresAt time.Time,
) (*user.Session, error) {
	session := user.Session{
		ID:        int64(len(f.sessions) + 1),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	f.sessions[hash] = session
	return &session, nil
}

func (f *fakeSessionStore) GetByToken(
	ctx context.Context,
	hash [32]byte,
) (*user.Session, error) {
	session, exists := f.sessions[hash]

	if !exists || !session.ExpiresAt.After(time.Now()) {
		return nil, repository.ErrSessionNotFound
	}

	return &session, nil
}

func (f *fakeSessionStore) Delete(
	ctx context.Context,
	hash [32]byte,
) error {
	delete(f.sessions, hash)
	return nil
}

func TestSessionCreateAndAuthenticate(t *testing.T) {
	store := newFakeSessionStore()
	service := NewSessionService(store)
	ctx := context.Background()

	token, created, err := service.Create(ctx, 42)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if token == "" {
		t.Fatal("токен пустой")
	}

	if created.UserID != 42 {
		t.Fatalf("неверный UserID: %d", created.UserID)
	}

	found, err := service.Authenticate(ctx, token)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}

	if found.UserID != created.UserID {
		t.Fatal("пользователь сессии не совпадает")
	}
}

func TestSessionRevoke(t *testing.T) {
	store := newFakeSessionStore()
	service := NewSessionService(store)
	ctx := context.Background()

	token, _, err := service.Create(ctx, 42)
	if err != nil {
		t.Fatal(err)
	}

	if err := service.Revoke(ctx, token); err != nil {
		t.Fatal(err)
	}

	_, err = service.Authenticate(ctx, token)
	if !errors.Is(err, repository.ErrSessionNotFound) {
		t.Fatalf("ожидалась ошибка отсутствия сессии: %v", err)
	}
}

func TestInvalidSession(t *testing.T) {
	service := NewSessionService(newFakeSessionStore())
	ctx := context.Background()

	_, err := service.Authenticate(ctx, "invalid-token")
	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("ожидалась ErrInvalidSession: %v", err)
	}

	_, _, err = service.Create(ctx, 0)
	if err == nil {
		t.Fatal("создана сессия с некорректным UserID")
	}
}

func TestExpiredSession(t *testing.T) {
	store := newFakeSessionStore()
	service := NewSessionService(store)
	ctx := context.Background()

	token, _, err := service.Create(ctx, 42)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := HashSessionToken(token)
	if err != nil {
		t.Fatal(err)
	}

	session := store.sessions[hash]
	session.ExpiresAt = time.Now().Add(-time.Minute)
	store.sessions[hash] = session

	_, err = service.Authenticate(ctx, token)
	if !errors.Is(err, repository.ErrSessionNotFound) {
		t.Fatalf("просроченная сессия принята: %v", err)
	}
}
