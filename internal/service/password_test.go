package service

import (
	"strings"
	"testing"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "MyStrongPassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(hash, password) {
		t.Fatal("хеш содержит исходный пароль")
	}

	ok, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("правильный пароль не прошёл проверку")
	}

	ok, err = VerifyPassword("WrongPassword", hash)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("неправильный пароль прошёл проверку")
	}
}

func TestPasswordSalt(t *testing.T) {
	first, err := HashPassword("same-password")
	if err != nil {
		t.Fatal(err)
	}

	second, err := HashPassword("same-password")
	if err != nil {
		t.Fatal(err)
	}

	if first == second {
		t.Fatal("хеши совпали: ожидались разные соли")
	}
}
