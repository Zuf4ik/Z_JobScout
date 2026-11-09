package provider

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

// RemotiveJob отражает формат отдельной вакансии в API Remotive.
type RemotiveJob struct {
	ID              json.Number `json:"id"`
	Name            string      `json:"title"`
	Company         string      `json:"company_name"`
	URL             string      `json:"url"`
	Location        string      `json:"candidate_required_location"`
	Description     string      `json:"description"`
	Salary          string      `json:"salary"`
	ExperienceYears *int        `json:"experience_years"`
}

type SearchResponse struct {
	Items []RemotiveJob `json:"jobs"`
	Found int           `json:"job-count"`
}

func LoadRemotiveVacancies(path string) ([]vacancy.Vacancy, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("чтение файла: %w", err)
	}

	var response SearchResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("разбор JSON: %w", err)
	}

	items := make([]vacancy.Vacancy, 0, len(response.Items))

	for _, job := range response.Items {
		items = append(items, vacancy.Vacancy{
			ID:              job.ID.String(),
			Source:          "remotive",
			Name:            job.Name,
			Company:         job.Company,
			URL:             job.URL,
			Location:        job.Location,
			Description:     job.Description,
			Salary:          job.Salary,
			ExperienceYears: job.ExperienceYears,
		})
	}

	return items, nil
}
