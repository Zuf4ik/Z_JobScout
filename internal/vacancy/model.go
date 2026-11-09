package vacancy

type Vacancy struct {
	ID              int    `json:"id"`
	Name            string `json:"title"`
	Company         string `json:"company_name"`
	URL             string `json:"url"`
	Location        string `json:"candidate_required_location"`
	Description     string `json:"description"`
	Salary          string `json:"salary"`
	ExperienceYears *int   `json:"experience_years"`
}
