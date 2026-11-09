package provider

import "testing"

func TestParseTrudvsemExperience(t *testing.T) {
	body := []byte(`{
		"status": "200",
		"results": {
			"vacancies": [
				{
					"vacancy": {
						"job-name": "Опыт не указан"
					}
				},
				{
					"vacancy": {
						"job-name": "Без опыта",
						"requirement": {
							"experience": 0
						}
					}
				},
				{
					"vacancy": {
						"job-name": "Опытный разработчик",
						"requirement": {
							"experience": 3
						}
					}
				}
			]
		}
	}`)

	items, err := ParseTrudvsemVacancies(body)
	if err != nil {
		t.Fatalf("ошибка разбора JSON: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("ожидалось 3 вакансии, получено: %d", len(items))
	}

	if items[0].ExperienceYears != nil {
		t.Error("отсутствующий опыт должен быть nil")
	}

	if items[1].ExperienceYears == nil {
		t.Fatal("нулевой опыт не должен быть nil")
	}

	if *items[1].ExperienceYears != 0 {
		t.Errorf(
			"ожидалось 0 лет, получено: %d",
			*items[1].ExperienceYears,
		)
	}

	if items[2].ExperienceYears == nil {
		t.Fatal("опыт 3 года не должен быть nil")
	}

	if *items[2].ExperienceYears != 3 {
		t.Errorf(
			"ожидалось 3 года, получено: %d",
			*items[2].ExperienceYears,
		)
	}
}

func TestParseTrudvsemID(t *testing.T) {
	body := []byte(`{
		"status": "200",
		"results": {
			"vacancies": [
				{
					"vacancy": {
						"id": "abc-123",
						"job-name": "Go Developer"
					}
				}
			]
		}
	}`)

	items, err := ParseTrudvsemVacancies(body)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 {
		t.Fatalf("ожидалась 1 вакансия, получено %d", len(items))
	}

	if items[0].ID != "abc-123" {
		t.Errorf("неверный ID: %s", items[0].ID)
	}

	if items[0].Source != "trudvsem" {
		t.Errorf("неверный источник: %s", items[0].Source)
	}
}
