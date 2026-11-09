package provider

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

type SearchResponse struct {
	Items []vacancy.Vacancy `json:"jobs"`
	Found int               `json:"job-count"`
}

func LoadRemotiveVacancies(path string) ([]vacancy.Vacancy, error) {
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
