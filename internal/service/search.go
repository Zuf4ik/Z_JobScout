package service

import (
	"fmt"
	"strings"

	"github.com/Zuf4ik/Z_JobScout/internal/provider"
	"github.com/Zuf4ik/Z_JobScout/internal/search"
	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

// Loader описывает любой источник вакансий.
type Loader func(value string) ([]vacancy.Vacancy, error)

type SearchService struct {
	loaders map[string]Loader
}

type SearchParams struct {
	Query         string
	Source        string
	File          string
	MaxExperience int // -1 означает отсутствие ограничения
}

type SearchResult struct {
	Vacancies   []vacancy.Vacancy
	Loaded      int
	TextMatches int
}

// NewSearchService создаёт сервис с источниками данных.
func NewSearchService() *SearchService {
	return &SearchService{
		loaders: map[string]Loader{
			"remotive":        provider.LoadRemotiveVacancies,
			"trudvsem":        provider.LoadTrudvsemVacancies,
			"trudvsem-online": provider.FetchTrudvsemVacancies,
		},
	}
}

func (s *SearchService) Search(
	params SearchParams,
) (SearchResult, error) {
	var result SearchResult

	query := strings.TrimSpace(params.Query)
	if query == "" {
		return result, fmt.Errorf("поисковый запрос не может быть пустым")
	}

	if params.MaxExperience < -1 {
		return result, fmt.Errorf("максимальный опыт должен быть -1 или больше")
	}

	if s == nil {
		return result, fmt.Errorf("сервис поиска не инициализирован")
	}

	loader, ok := s.loaders[params.Source]
	if !ok || loader == nil {
		return result, fmt.Errorf("неизвестный источник: %s", params.Source)
	}

	value := params.File
	if params.Source == "trudvsem-online" {
		value = query
	}

	items, err := loader(value)
	if err != nil {
		return result, fmt.Errorf("загрузка вакансий: %w", err)
	}

	result.Loaded = len(items)

	filtered := search.FilterByText(items, query)
	result.TextMatches = len(filtered)

	if params.MaxExperience >= 0 {
		filtered = search.FilterByExperience(
			filtered,
			params.MaxExperience,
		)
	}

	result.Vacancies = filtered
	return result, nil
}
