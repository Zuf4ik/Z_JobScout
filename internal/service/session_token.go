package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const sessionTokenBytes = 32

func GenerateSessionToken() (string, [32]byte, error) {
	raw := make([]byte, sessionTokenBytes)

	if _, err := rand.Read(raw); err != nil {
		return "", [32]byte{}, err
	}

	token := hex.EncodeToString(raw)
	hash := sha256.Sum256(raw)

	return token, hash, nil
}

func HashSessionToken(token string) ([32]byte, error) {
	raw, err := hex.DecodeString(token)
	if err != nil || len(raw) != sessionTokenBytes {
		return [32]byte{}, errors.New("некорректный токен сессии")
	}

	return sha256.Sum256(raw), nil
}
