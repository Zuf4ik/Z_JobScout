package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Zuf4ik/Z_JobScout/internal/repository"
	"github.com/Zuf4ik/Z_JobScout/internal/user"
)

// fakeUserStore имитирует базу данных в памяти.
type fakeUserStore struct {
	users map[string]user.User
}

func newFakeUserStore() *fakeUserStore {
	return &fakeUserStore{
		users: make(map[string]user.User),
	}
}

func (f *fakeUserStore) CreateUser(
	ctx context.Context,
	email string,
	passwordHash string,
) (*user.User, error) {
	if _, exists := f.users[email]; exists {
		return nil, repository.ErrEmailTaken
	}

	u := user.User{
		ID:           int64(len(f.users) + 1),
		Email:        email,
		PasswordHash: passwordHash,
	}

	f.users[email] = u
	return &u, nil
}

func (f *fakeUserStore) GetUserByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	u, exists := f.users[email]
	if !exists {
		return nil, repository.ErrUserNotFound
	}

	return &u, nil
}

func TestAuthRegisterAndLogin(t *testing.T) {
	store := newFakeUserStore()
	auth := NewAuthService(store)
	ctx := context.Background()

	registered, err := auth.Register(
		ctx,
		" Test@Example.com ",
		"StrongPassword123!",
	)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	if registered.Email != "test@example.com" {
		t.Fatalf("Неправильный email: %s", registered.Email)
	}

	if registered.PasswordHash != "" {
		t.Fatal("Регистрация вернула хеш пароля")
	}

	stored := store.users["test@example.com"]
	if stored.PasswordHash == "" {
		t.Fatal("Пароль не был захеширован")
	}

	loggedIn, err := auth.Login(
		ctx,
		"test@example.com",
		"StrongPassword123!",
	)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if loggedIn.ID != registered.ID {
		t.Fatal("ID пользователей не совпадают")
	}

	if loggedIn.PasswordHash != "" {
		t.Fatal("Login вернул хеш пароля")
	}

	_, err = auth.Login(
		ctx,
		"test@example.com",
		"WrongPassword",
	)
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Ожидалась ErrInvalidCredentials, получено: %v", err)
	}
}

func TestAuthValidation(t *testing.T) {
	auth := NewAuthService(newFakeUserStore())
	ctx := context.Background()

	_, err := auth.Register(
		ctx,
		"wrong-email",
		"StrongPassword123!",
	)
	if !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("Ожидалась ErrInvalidEmail, получено: %v", err)
	}

	_, err = auth.Register(
		ctx,
		"test@example.com",
		"short",
	)
	if !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("Ожидалась ErrWeakPassword, получено: %v", err)
	}
}
