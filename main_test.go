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

func TestFilterByExperience(t *testing.T) {
	zero := 0
	one := 1
	three := 3

	items := []Vacancy{
		{ID: 1, ExperienceYears: &zero},
		{ID: 2, ExperienceYears: &one},
		{ID: 3, ExperienceYears: &three},
		{ID: 4},
	}

	got := filterByExperience(items, 1)

	if len(got) != 3 {
		t.Fatalf(
			"ожидалось 3 вакансии, получено: %d",
			len(got),
		)
	}

	expectedIDs := []int{1, 2, 4}

	for i, vacancy := range got {
		if vacancy.ID != expectedIDs[i] {
			t.Errorf(
				"позиция %d: ожидался ID %d, получен %d",
				i,
				expectedIDs[i],
				vacancy.ID,
			)
		}
	}
}
