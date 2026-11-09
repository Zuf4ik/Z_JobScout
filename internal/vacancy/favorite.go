package vacancy

import "time"

type Favorite struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	VacancyID   string    `json:"vacancy_id"`
	Source      string    `json:"source"`
	Title       string    `json:"title"`
	CompanyName string    `json:"company_name"`
	URL         string    `json:"url"`
	SavedAt     time.Time `json:"saved_at"`
}
