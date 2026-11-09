package search

import (
	"strings"

	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

// FilterByText ищет запрос в названии и описании вакансии.
func FilterByText(
	items []vacancy.Vacancy,
	query string,
) []vacancy.Vacancy {
	var filtered []vacancy.Vacancy

	query = strings.ToLower(strings.TrimSpace(query))

	for _, v := range items {
		text := strings.ToLower(
			v.Name + " " + v.Description,
		)

		if strings.Contains(text, query) {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

// FilterByExperience исключает вакансии,
// требующие опыта больше указанного.
func FilterByExperience(
	items []vacancy.Vacancy,
	maxYears int,
) []vacancy.Vacancy {
	var filtered []vacancy.Vacancy

	for _, v := range items {
		if v.ExperienceYears != nil &&
			*v.ExperienceYears > maxYears {
			continue
		}

		filtered = append(filtered, v)
	}

	return filtered
}
