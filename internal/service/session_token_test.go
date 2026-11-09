package service

import "testing"

func TestSessionToken(t *testing.T) {
	token, hash, err := GenerateSessionToken()
	if err != nil {
		t.Fatal(err)
	}

	if len(token) != 64 {
		t.Fatalf("неожиданная длина токена: %d", len(token))
	}

	computed, err := HashSessionToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if hash != computed {
		t.Fatal("хеши токена не совпадают")
	}

	secondToken, _, err := GenerateSessionToken()
	if err != nil {
		t.Fatal(err)
	}

	if token == secondToken {
		t.Fatal("разные сессии получили одинаковый токен")
	}
}

func TestInvalidSessionToken(t *testing.T) {
	_, err := HashSessionToken("invalid-token")

	if err == nil {
		t.Fatal("ожидалась ошибка некорректного токена")
	}
}
