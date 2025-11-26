package models

type Vacancy struct {
	ID           string   `json:"id"`
	CompanyID    string   `json:"company_id"`
	Profession   string   `json:"profession"`
	Position     string   `json:"position"`
	Requirements []string `json:"requirements"`
	Tasks        []string `json:"tasks"`
	TaskIdeas    []string `json:"task_ideas"`
	Metrics      []string `json:"metrics"`
	IsActive     bool     `json:"is_active"`
	Duration     int      `json:"duration"`
	CreatedAt    string   `json:"created_at"`
}

type InterviewResult struct {
	ID          string   `json:"id"`
	VacancyID   string   `json:"vacancy_id"`
	Profession  string   `json:"profession"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	ResumeLink  string   `json:"resume_link"`
	Tasks       []string `json:"tasks"`
	Solutions   []string `json:"solutions"`
	ChatHistory []string `json:"chat_history"`
	Metrics     []string `json:"metrics"`
	CreatedAt   string   `json:"created_at"`
}
