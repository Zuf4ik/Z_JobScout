package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Vacancy struct {
	ID              int    `json:"id"`
	Name            string `json:"title"`
	Company         string `json:"company_name"`
	URL             string `json:"url"`
	Location        string `json:"candidate_required_location"`
	Description     string `json:"description"`
	Salary          string `json:"salary"`
	ExperienceYears *int   `json:"experience_years"`
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

	path := "examples/vacancies.json"
	if len(os.Args) >= 3 {
		path = os.Args[2]
	}

	source := "remotive"
	if len(os.Args) >= 4 {
		source = os.Args[3]
	}

	maxYears := -1

	if len(os.Args) >= 5 {
		value, err := strconv.Atoi(os.Args[4])
		if err != nil {
			return fmt.Errorf(
				"максимальный опыт должен быть целым числом: %w",
				err,
			)
		}

		if value < 0 {
			return fmt.Errorf("максимальный опыт не может быть отрицательным")
		}

		maxYears = value
	}

	var items []Vacancy
	var err error

	switch source {
	case "remotive":
		items, err = loadVacancies(path)
	case "trudvsem":
		items, err = loadTrudvsemVacancies(path)
	case "trudvsem-online":
		items, err = fetchTrudvsemVacancies(query)
	default:
		return fmt.Errorf("неизвестный источник: %s", source)
	}

	if err != nil {
		return err
	}

	if source == "trudvsem-online" {
		fmt.Println("Источник: Работа России — https://trudvsem.ru")
		fmt.Println("Получена первая страница, до 20 вакансий.")
	} else {
		fmt.Println("Файл данных:", path)

		if path == "examples/vacancies.json" {
			fmt.Println("Демонстрационные данные: вакансии вымышленные.")
		}
	}

	vacancies := filterVacancies(items, query)

	textMatches := len(vacancies)

	if maxYears >= 0 {
		vacancies = filterByExperience(vacancies, maxYears)
	}

	fmt.Printf(
		"Загружено: %d. Совпадений по тексту: %d. После фильтров: %d\n",
		len(items),
		textMatches,
		len(vacancies),
	)

	if maxYears >= 0 {
		fmt.Printf(
			"Максимальный требуемый опыт, лет: %d. Неизвестный опыт сохраняется.\n",
			maxYears,
		)
	}

	fmt.Println()
	for _, vacancy := range vacancies {
		salary := strings.TrimSpace(vacancy.Salary)
		if salary == "" {
			salary = "не указана"
		}

		fmt.Println(vacancy.Name)
		fmt.Println("Компания:", vacancy.Company)
		fmt.Println("Регион:", vacancy.Location)
		fmt.Println("Зарплата:", salary)

		if vacancy.ExperienceYears == nil {
			fmt.Println("Требуемый опыт: не указан")
		} else {
			fmt.Printf(
				"Требуемый опыт, лет: %d\n",
				*vacancy.ExperienceYears,
			)
		}

		fmt.Println(vacancy.URL)
		fmt.Println()
	}

	return nil
}

func loadVacancies(path string) ([]Vacancy, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("чтение файла: %w", err)
	}

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("разбор JSON: %w", err)
	}

	return result.Items, nil
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

func filterByExperience(items []Vacancy, maxYears int) []Vacancy {
	var filtered []Vacancy

	for _, vacancy := range items {
		if vacancy.ExperienceYears != nil {
			if *vacancy.ExperienceYears > maxYears {
				continue
			}
		}

		filtered = append(filtered, vacancy)
	}

	return filtered
}
