package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Zuf4ik/Z_JobScout/internal/service"
	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

// Searcher описывает то, что необходимо HTTP-слою.
type Searcher interface {
	Search(params service.SearchParams) (service.SearchResult, error)
}

type Handler struct {
	searcher Searcher
}

type VacanciesResponse struct {
	Items       []vacancy.Vacancy `json:"items"`
	Loaded      int               `json:"loaded"`
	TextMatches int               `json:"text_matches"`
	Total       int               `json:"total"`
}

func NewHandler(searcher Searcher) http.Handler {
	h := &Handler{
		searcher: searcher,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/v1/vacancies", h.vacancies)

	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "JobScout API",
	})
}

func (h *Handler) vacancies(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))

	if query == "" {
		writeError(w, http.StatusBadRequest, "параметр query обязателен")
		return
	}

	if len(query) > 200 {
		writeError(w, http.StatusBadRequest, "слишком длинный поисковый запрос")
		return
	}

	source := r.URL.Query().Get("source")
	if source == "" {
		source = "trudvsem-online"
	}

	if source != "trudvsem-online" {
		writeError(w, http.StatusBadRequest, "неподдерживаемый источник")
		return
	}

	maxExperience := -1

	rawExperience := r.URL.Query().Get("max_experience")
	if rawExperience != "" {
		value, err := strconv.Atoi(rawExperience)
		if err != nil || value < -1 {
			writeError(w, http.StatusBadRequest,
				"max_experience должен быть целым числом от -1")
			return
		}
		maxExperience = value
	}

	result, err := h.searcher.Search(service.SearchParams{
		Query:         query,
		Source:        source,
		MaxExperience: maxExperience,
	})
	if err != nil {
		// Внутренние подробности не показываем клиенту.
		writeError(w, http.StatusBadGateway, "не удалось получить вакансии")
		return
	}

	items := result.Vacancies
	if items == nil {
		items = []vacancy.Vacancy{}
	}

	writeJSON(w, http.StatusOK, VacanciesResponse{
		Items:       items,
		Loaded:      result.Loaded,
		TextMatches: result.TextMatches,
		Total:       len(items),
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Если запись уже началась, HTTP-статус менять поздно.
		// Ошибку фиксируем на стороне сервера.
		fmt.Printf("ошибка записи JSON: %v\n", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
