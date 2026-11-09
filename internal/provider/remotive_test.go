package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRemotiveVacancies(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vacancies.json")

	data := `{
        "job-count": 1,
        "jobs": [
            {
                "id": 123,
                "title": "Junior Go Developer",
                "company_name": "Example",
                "description": "Go и PostgreSQL"
            }
        ]
    }`

	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}

	items, err := LoadRemotiveVacancies(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 {
		t.Fatalf("ожидалась 1 вакансия, получено %d", len(items))
	}

	if items[0].Name != "Junior Go Developer" {
		t.Errorf("неверное название: %s", items[0].Name)
	}

	if items[0].ID != "123" {
		t.Errorf("ожидался ID 123, получен %s", items[0].ID)
	}

	if items[0].Source != "remotive" {
		t.Errorf("неверный источник: %s", items[0].Source)
	}
}
