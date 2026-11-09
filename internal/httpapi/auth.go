package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Zuf4ik/Z_JobScout/internal/repository"
	"github.com/Zuf4ik/Z_JobScout/internal/service"
)

type AuthHandler struct {
	auth     *service.AuthService
	sessions *service.SessionService
	secure   bool
}

func NewAuthHandler(
	auth *service.AuthService,
	sessions *service.SessionService,
	secureCookie bool,
) *AuthHandler {
	return &AuthHandler{
		auth:     auth,
		sessions: sessions,
		secure:   secureCookie,
	}
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

// readAuthRequest ограничивает размер и проверяет JSON.
func readAuthRequest(
	w http.ResponseWriter,
	r *http.Request,
) (authRequest, error) {
	var req authRequest

	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return req, err
	}

	// Запрещаем несколько JSON-объектов в одном запросе.
	if err := decoder.Decode(new(any)); err != io.EOF {
		if err == nil {
			return req, errors.New("лишние данные в запросе")
		}
		return req, err
	}

	return req, nil
}

func (h *AuthHandler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	req, err := readAuthRequest(w, r)
	if err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	u, err := h.auth.Register(
		r.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailTaken):
			http.Error(w, "Email уже зарегистрирован", http.StatusConflict)
		case errors.Is(err, service.ErrInvalidEmail),
			errors.Is(err, service.ErrWeakPassword):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	writeAuthJSON(w, http.StatusCreated, authResponse{
		ID:    u.ID,
		Email: u.Email,
	})
}

func (h *AuthHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	req, err := readAuthRequest(w, r)
	if err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	u, err := h.auth.Login(
		r.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(
				w,
				"Неверный email или пароль",
				http.StatusUnauthorized,
			)
			return
		}
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	token, _, err := h.sessions.Create(r.Context(), u.ID)
	if err != nil {
		http.Error(w, "Ошибка создания сессии", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jobscout_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(service.SessionTTL / time.Second),
	})

	writeAuthJSON(w, http.StatusOK, authResponse{
		ID:    u.ID,
		Email: u.Email,
	})
}

func writeAuthJSON(
	w http.ResponseWriter,
	status int,
	data authResponse,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {
	cookie, err := r.Cookie("jobscout_session")
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	session, err := h.sessions.Authenticate(
		r.Context(),
		cookie.Value,
	)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	writeAuthJSON(w, http.StatusOK, authResponse{
		ID: session.UserID,
	})
}

func (h *AuthHandler) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	cookie, err := r.Cookie("jobscout_session")

	if err == nil {
		if err := h.sessions.Revoke(
			r.Context(),
			cookie.Value,
		); err != nil {
			// Пока фиксируем ошибку, не раскрывая детали клиенту.
			http.Error(w, "Ошибка завершения сессии",
				http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jobscout_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
