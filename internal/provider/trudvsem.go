package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

type TrudvsemResponse struct {
	Status string `json:"status"`

	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`

	Results struct {
		Vacancies []struct {
			Vacancy TrudvsemVacancy `json:"vacancy"`
		} `json:"vacancies"`
	} `json:"results"`
}

type TrudvsemVacancy struct {
	ID           string `json:"id"`
	Name         string `json:"job-name"`
	URL          string `json:"vac_url"`
	Requirements string `json:"requirements"`
	Duty         string `json:"duty"`
	Salary       string `json:"salary"`
	Requirement  struct {
		Experience *int `json:"experience"`
	} `json:"requirement"`

	Company struct {
		Name string `json:"name"`
	} `json:"company"`

	Region struct {
		Name string `json:"name"`
	} `json:"region"`
}

func LoadTrudvsemVacancies(path string) ([]vacancy.Vacancy, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("чтение файла Работы России: %w", err)
	}

	return ParseTrudvsemVacancies(body)
}

func ParseTrudvsemVacancies(body []byte) ([]vacancy.Vacancy, error) {

	var response TrudvsemResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("разбор ответа Работы России: %w", err)
	}

	if response.Status != "200" {
		return nil, fmt.Errorf(
			"Работа России вернула статус %q",
			response.Status,
		)
	}

	var items []vacancy.Vacancy

	for _, entry := range response.Results.Vacancies {
		v := entry.Vacancy

		items = append(items, vacancy.Vacancy{
			Name:            v.Name,
			Company:         v.Company.Name,
			URL:             v.URL,
			Location:        v.Region.Name,
			Description:     v.Requirements + "\n" + v.Duty,
			Salary:          v.Salary,
			ExperienceYears: v.Requirement.Experience,
		})
	}

	return items, nil
}

func FetchTrudvsemVacancies(query string) ([]vacancy.Vacancy, error) {
	params := url.Values{}
	params.Set("text", query)
	params.Set("limit", "20")
	params.Set("offset", "0")

	address := "https://opendata.trudvsem.ru/api/v1/vacancies?" +
		params.Encode()

	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return nil, fmt.Errorf("создание запроса: %w", err)
	}

	req.Header.Set("User-Agent", "JobScout/0.1")

	client := &http.Client{
		Timeout: 45 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("запрос к Работе России: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("чтение HTTP-ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Работа России: HTTP %d",
			resp.StatusCode,
		)
	}

	return ParseTrudvsemVacancies(body)
}
