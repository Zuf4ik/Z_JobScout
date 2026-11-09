package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Zuf4ik/Z_JobScout/internal/search"
	"github.com/Zuf4ik/Z_JobScout/internal/service"
	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

type Vacancy = vacancy.Vacancy

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
		return fmt.Errorf(
			"используйте именованные параметры; справка: go run . -h",
		)
	}

	query := strings.TrimSpace(*queryFlag)
	if query == "" {
		return fmt.Errorf("укажите запрос: go run . -query Golang")
	}

	source := *sourceFlag
	path := *fileFlag
	maxYears := *maxYearsFlag

	if maxYears < -1 {
		return fmt.Errorf(
			"максимальный опыт должен быть -1 или неотрицательным числом",
		)
	}

	searchService := service.NewSearchService()

	result, err := searchService.Search(service.SearchParams{
		Query:         query,
		Source:        source,
		File:          path,
		MaxExperience: maxYears,
	})
	if err != nil {
		return err
	}

	vacancies := result.Vacancies

	if source == "trudvsem-online" {
		fmt.Println("Источник: Работа России — https://trudvsem.ru")
		fmt.Println("Получена первая страница, до 20 вакансий.")
	} else {
		fmt.Println("Файл данных:", path)

		if path == "examples/vacancies.json" {
			fmt.Println("Демонстрационные данные: вакансии вымышленные.")
		}
	}

	fmt.Printf(
		"Загружено: %d. Совпадений по тексту: %d. После фильтров: %d\n",
		result.Loaded,
		result.TextMatches,
		len(vacancies),
	)

	if maxYears >= 0 {
		fmt.Printf(
			"Максимальный требуемый опыт, лет: %d. Неизвестный опыт сохраняется.\n",
			maxYears,
		)
	}

	fmt.Println()

	for _, v := range vacancies {
		salary := strings.TrimSpace(v.Salary)
		if salary == "" {
			salary = "не указана"
		}

		fmt.Println(v.Name)
		fmt.Println("Компания:", v.Company)
		fmt.Println("Регион:", v.Location)
		fmt.Println("Зарплата:", salary)

		if v.ExperienceYears == nil {
			fmt.Println("Требуемый опыт: не указан")
		} else {
			fmt.Printf(
				"Требуемый опыт, лет: %d\n",
				*v.ExperienceYears,
			)
		}

		fmt.Println(v.URL)
		fmt.Println()
	}

	return nil
}

// Временные обёртки для совместимости с существующими тестами.

func filterVacancies(items []Vacancy, query string) []Vacancy {
	return search.FilterByText(items, query)
}

func filterByExperience(items []Vacancy, maxYears int) []Vacancy {
	return search.FilterByExperience(items, maxYears)
}
