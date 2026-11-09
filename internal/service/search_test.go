package service

import (
	"errors"
	"testing"

	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

func TestSearchService(t *testing.T) {
	zero := 0
	three := 3

	s := &SearchService{
		loaders: map[string]Loader{
			"test": func(value string) ([]vacancy.Vacancy, error) {
				return []vacancy.Vacancy{
					{
						Name:            "Junior Go Developer",
						ExperienceYears: &zero,
					},
					{
						Name:            "Senior Go Developer",
						ExperienceYears: &three,
					},
					{
						Name: "Python Developer",
					},
				}, nil
			},
		},
	}

	result, err := s.Search(SearchParams{
		Query:         "Go",
		Source:        "test",
		MaxExperience: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.Loaded != 3 {
		t.Errorf("загружено %d, ожидалось 3", result.Loaded)
	}

	if result.TextMatches != 2 {
		t.Errorf("текстовых совпадений %d, ожидалось 2", result.TextMatches)
	}

	if len(result.Vacancies) != 1 {
		t.Fatalf("найдено %d вакансий, ожидалась 1", len(result.Vacancies))
	}

	if result.Vacancies[0].Name != "Junior Go Developer" {
		t.Errorf("неожиданная вакансия: %s", result.Vacancies[0].Name)
	}
}

func TestSearchServiceInvalidParams(t *testing.T) {
	s := NewSearchService()

	tests := []SearchParams{
		{Query: "", Source: "trudvsem-online", MaxExperience: -1},
		{Query: "Go", Source: "unknown", MaxExperience: -1},
		{Query: "Go", Source: "trudvsem-online", MaxExperience: -2},
	}

	for _, params := range tests {
		_, err := s.Search(params)
		if err == nil {
			t.Errorf("ожидалась ошибка для параметров: %+v", params)
		}
	}
}

func TestSearchServiceProviderError(t *testing.T) {
	s := &SearchService{
		loaders: map[string]Loader{
			"test": func(string) ([]vacancy.Vacancy, error) {
				return nil, errors.New("источник недоступен")
			},
		},
	}

	_, err := s.Search(SearchParams{
		Query:         "Go",
		Source:        "test",
		MaxExperience: -1,
	})

	if err == nil {
		t.Fatal("ожидалась ошибка загрузки")
	}
}
