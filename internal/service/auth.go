package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/Zuf4ik/Z_JobScout/internal/repository"
	"github.com/Zuf4ik/Z_JobScout/internal/user"
)

var (
	ErrInvalidCredentials = errors.New("неверный email или пароль")
	ErrInvalidEmail       = errors.New("некорректный email")
	ErrWeakPassword       = errors.New("пароль должен содержать от 12 до 128 символов")
)

type UserStore interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
}

type AuthService struct {
	users UserStore
}

func NewAuthService(users UserStore) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) Register(
	ctx context.Context,
	email string,
	password string,
) (*user.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return nil, ErrInvalidEmail
	}

	if len([]rune(password)) < 12 || len([]rune(password)) > 128 {
		return nil, ErrWeakPassword
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования: %w", err)
	}

	u, err := s.users.CreateUser(ctx, email, passwordHash)
	if err != nil {
		return nil, err
	}

	// Хеш не должен попадать в HTTP-ответ.
	u.PasswordHash = ""

	return u, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*user.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	u, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	valid, err := VerifyPassword(password, u.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки хеша: %w", err)
	}
	if !valid {
		return nil, ErrInvalidCredentials
	}

	u.PasswordHash = ""
	return u, nil
}
