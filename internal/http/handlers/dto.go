package handlers

type CreateVacancyRequest struct {
	CompanyID    string   `json:"company_id" binding:"required"`
	Profession   string   `json:"profession" binding:"required"`
	Position     string   `json:"position" binding:"required"`
	Requirements []string `json:"requirements"`
	Tasks        []string `json:"tasks"`
	TaskIdeas    []string `json:"task_ideas"`
	Metrics      []string `json:"metrics"`
	IsActive     bool     `json:"is_active"`
	Duration     int      `json:"duration"`
}

type GetVacanciesRequest struct {
	CompanyID string `form:"company_id" binding:"required"`
}

type AddInterviewRequest struct {
	Profession  string   `json:"profession" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Surname     string   `json:"surname" binding:"required"`
	ResumeLink  string   `json:"resume_link"`
	Tasks       []string `json:"tasks"`
	Solutions   []string `json:"solutions"`
	ChatHistory []string `json:"chat_history"`
	Metrics     []string `json:"metrics"`
}
