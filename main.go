package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
	queryFlag := flag.String(
		"query", "",
		"Поисковый запрос",
	)

	sourceFlag := flag.String(
		"source", "remotive",
		"Источник: remotive, trudvsem или trudvsem-online",
	)

	fileFlag := flag.String(
		"file", "examples/vacancies.json",
		"Путь к JSON-файлу",
	)

	maxYearsFlag := flag.Int(
		"max-experience", -1,
		"Максимальный требуемый опыт в годах; -1 — без ограничения",
	)

	telegramFlag := flag.Bool(
		"telegram", false,
		"Запустить Telegram-бота",
	)

	flag.Parse()

	if *telegramFlag {
		return runTelegram()
	}

	if flag.NArg() != 0 {
		return fmt.Errorf("используйте именованные параметры; справка: go run . -h")
	}

	query := strings.TrimSpace(*queryFlag)
	if query == "" {
		return fmt.Errorf("укажите запрос: go run . -query Golang")
	}

	source := *sourceFlag
	path := *fileFlag
	maxYears := *maxYearsFlag

	if maxYears < -1 {
		return fmt.Errorf("максимальный опыт должен быть -1 или неотрицательным числом")
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
