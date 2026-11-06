package main

import (
	"encoding/json"
	"fmt"

	"log"

	"os"
	"strings"
)

type Vacancy struct {
	ID          int    `json:"id"`
	Name        string `json:"title"`
	Company     string `json:"company_name"`
	URL         string `json:"url"`
	Location    string `json:"candidate_required_location"`
	Description string `json:"description"`
}

type SearchResponse struct {
	Items []Vacancy `json:"jobs"`
	Found int       `json:"job-count"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	if len(os.Args) < 2 {
		fmt.Println("Укажите поисковый запрос: go run . Golang")
		return nil
	}

	query := strings.TrimSpace(os.Args[1])

	if query == "" {
		fmt.Println("Поисковый запрос не должен быть пустым.")
		return nil
	}

	body, err := os.ReadFile("remotive-response.json")
	if err != nil {
		return fmt.Errorf("чтение файла: %w", err)
	}

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("разбор JSON: %w", err)
	}

	vacancies := filterVacancies(result.Items, query)

	fmt.Printf(
		"Загружено: %d. Подходит: %d\n\n",
		len(result.Items),
		len(vacancies),
	)

	for _, vacancy := range vacancies {
		fmt.Println(vacancy.Name)
		fmt.Println("Источник: Remotive")
		fmt.Println(vacancy.URL)
		fmt.Println()
	}

	return nil
}

func filterVacancies(items []Vacancy, query string) []Vacancy {
	var filtered []Vacancy
	query = strings.ToLower(strings.TrimSpace(query))

	for _, vacancy := range items {
		text := strings.ToLower(
			vacancy.Name + " " + vacancy.Description,
		)

		if strings.Contains(text, query) {
			filtered = append(filtered, vacancy)
		}
	}

	return filtered
}
