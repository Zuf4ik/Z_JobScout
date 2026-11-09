package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 3
	argonMemory  = 64 * 1024
	argonThreads = 2
	argonKeyLen  = 32
	saltLen      = 16
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("пароль не может быть пустым")
	}

	salt := make([]byte, saltLen)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("ошибка генерации соли: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	encoder := base64.RawStdEncoding

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory,
		argonTime,
		argonThreads,
		encoder.EncodeToString(salt),
		encoder.EncodeToString(hash),
	), nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	var memory, iterations, threads uint32
	var version int

	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("неверный формат хеша")
	}

	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != 19 {
		return false, errors.New("неподдерживаемая версия Argon2id")
	}

	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory, &iterations, &threads,
	); err != nil {
		return false, errors.New("неверные параметры Argon2id")
	}

	// Не допускаем произвольные затраты памяти и CPU
	// из повреждённого или недоверенного хеша.
	if memory != argonMemory ||
		iterations != argonTime ||
		threads != argonThreads {
		return false, errors.New("неподдерживаемые параметры Argon2id")
	}

	encoder := base64.RawStdEncoding

	salt, err := encoder.DecodeString(parts[4])
	if err != nil || len(salt) != saltLen {
		return false, errors.New("некорректная соль")
	}

	expectedHash, err := encoder.DecodeString(parts[5])
	if err != nil || len(expectedHash) != argonKeyLen {
		return false, errors.New("некорректный хеш")
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		uint8(threads),
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(
		actualHash,
		expectedHash,
	) == 1, nil
}
