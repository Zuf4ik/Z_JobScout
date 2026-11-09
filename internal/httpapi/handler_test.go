package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Zuf4ik/Z_JobScout/internal/service"
	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

type fakeSearcher struct {
	result service.SearchResult
	err    error
	params service.SearchParams
	called bool
}

func (f *fakeSearcher) Search(
	params service.SearchParams,
) (service.SearchResult, error) {
	f.called = true
	f.params = params
	return f.result, f.err
}

func TestHealth(t *testing.T) {
	handler := NewHandler(&fakeSearcher{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался статус 200, получен %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("неожиданный ответ: %s", rec.Body.String())
	}
}

func TestVacanciesSuccess(t *testing.T) {
	fake := &fakeSearcher{
		result: service.SearchResult{
			Loaded:      3,
			TextMatches: 2,
			Vacancies: []vacancy.Vacancy{
				{Name: "Junior Go Developer"},
			},
		},
	}

	handler := NewHandler(fake)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/vacancies?query=Golang&max_experience=1",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался статус 200, получен %d", rec.Code)
	}

	var response VacanciesResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.Total != 1 {
		t.Errorf("ожидалась 1 вакансия, получено %d", response.Total)
	}

	if response.Items[0].Name != "Junior Go Developer" {
		t.Errorf("неверная вакансия: %s", response.Items[0].Name)
	}

	if fake.params.Source != "trudvsem-online" {
		t.Errorf("неверный источник: %s", fake.params.Source)
	}

	if fake.params.Query != "Golang" ||
		fake.params.MaxExperience != 1 {
		t.Errorf("неверные параметры поиска: %+v", fake.params)
	}
}

func TestVacanciesInvalidParams(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"пустой запрос", "/api/v1/vacancies"},
		{"неверный опыт", "/api/v1/vacancies?query=Go&max_experience=abc"},
		{"отрицательный опыт", "/api/v1/vacancies?query=Go&max_experience=-2"},
		{"неподдерживаемый источник", "/api/v1/vacancies?query=Go&source=unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &fakeSearcher{}
			handler := NewHandler(fake)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("ожидался 400, получен %d", rec.Code)
			}

			if fake.called {
				t.Error("поисковый сервис не должен вызываться")
			}
		})
	}
}

func TestVacanciesProviderError(t *testing.T) {
	fake := &fakeSearcher{
		err: errors.New("API недоступен"),
	}
	handler := NewHandler(fake)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/vacancies?query=Go",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("ожидался статус 502, получен %d", rec.Code)
	}

	if strings.Contains(rec.Body.String(), "API недоступен") {
		t.Error("внутренняя ошибка не должна передаваться клиенту")
	}
}

func TestVacanciesMethodNotAllowed(t *testing.T) {
	handler := NewHandler(&fakeSearcher{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/vacancies",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("ожидался статус 405, получен %d", rec.Code)
	}
}
