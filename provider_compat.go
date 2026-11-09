package main

import (
	"github.com/Zuf4ik/Z_JobScout/internal/provider"
)

// Временные обёртки для совместимости со старым кодом.

func loadTrudvsemVacancies(path string) ([]Vacancy, error) {
	return provider.LoadTrudvsemVacancies(path)
}

func parseTrudvsemVacancies(body []byte) ([]Vacancy, error) {
	return provider.ParseTrudvsemVacancies(body)
}

func fetchTrudvsemVacancies(query string) ([]Vacancy, error) {
	return provider.FetchTrudvsemVacancies(query)
}
