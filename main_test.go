package main

import "testing"

func TestFilterVacancies(t *testing.T) {
	items := []Vacancy{
		{
			ID:          1,
			Name:        "Backend Developer",
			Description: "Разработка на GOLANG",
		},
		{
			ID:          2,
			Name:        "Python Developer",
			Description: "Разработка на Python",
		},
	}

	got := filterVacancies(items, "golang")

	if len(got) != 1 {
		t.Fatalf(
			"ожидалась 1 вакансия, получено: %d",
			len(got),
		)
	}

	if got[0].ID != 1 {
		t.Errorf(
			"ожидалась вакансия с ID 1, получен ID %d",
			got[0].ID,
		)
	}
}

func TestFilterVacanciesNoMatches(t *testing.T) {
	items := []Vacancy{
		{
			ID:          1,
			Name:        "Python Developer",
			Description: "Разработка на Python",
		},
	}

	got := filterVacancies(items, "Java")

	// Здесь проверь, что длина got равна нулю.

	if len(got) != 0 {
		t.Errorf(
			"ожидалось 0 вакансий, получено: %d",
			len(got),
		)
	}
}
