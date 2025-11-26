package handlers

import (
	"net/http"

	"github.com/TOPfiIT/vacancy-service/internal/domain/models"
	"github.com/TOPfiIT/vacancy-service/internal/services"
	"github.com/TOPfiIT/vacancy-service/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type VacancyHandler struct {
	service *services.VacancyService
}

func NewVacancyHandler(service *services.VacancyService) *VacancyHandler {
	return &VacancyHandler{
		service: service,
	}
}

func (h *VacancyHandler) CreateVacancy(c *gin.Context) {
	var req models.CreateVacancyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID := middlewares.GetCompanyID(c)
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "company id not found"})
		return
	}
	req.CompanyID = companyID

	vacancy, err := h.service.CreateVacancy(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vacancy)
}

func (h *VacancyHandler) GetVacancy(c *gin.Context) {
	vacancyID := c.Param("vacancy_id")
	if vacancyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty id"})
		return
	}

	vacancy, err := h.service.GetVacancy(c.Request.Context(), vacancyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vacancy)
}

func (h *VacancyHandler) GetVacancies(c *gin.Context) {
	companyID := middlewares.GetCompanyID(c)
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "company id not found"})
		return
	}

	vacancies, err := h.service.GetVacanciesByCompany(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.GetVacancyResponse, len(vacancies))
	for i, v := range vacancies {
		responses[i] = models.GetVacancyResponse{
			ID:           v.ID,
			CompanyID:    v.CompanyID,
			Profession:   v.Profession,
			Position:     v.Position,
			Requirements: v.Requirements,
			Tasks:        v.Tasks,
			TaskIdeas:    v.TaskIdeas,
			Metrics:      v.Metrics,
			IsActive:     v.IsActive,
			Duration:     v.Duration,
			CreatedAt:    v.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, responses)
}

func (h *VacancyHandler) GetInterviewResults(c *gin.Context) {
	vacancyID := c.Param("vacancy_id")
	if vacancyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty id"})
		return
	}

	companyID := middlewares.GetCompanyID(c)
	vacancy, err := h.service.GetVacancy(c.Request.Context(), vacancyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if vacancy.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	results, err := h.service.GetInterviewResults(c.Request.Context(), vacancyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.InterviewResultResponse, len(results))
	for i, result := range results {
		responses[i] = models.InterviewResultResponse{
			ID:          result.ID,
			VacancyID:   result.VacancyID,
			Profession:  result.Profession,
			Name:        result.Name,
			Surname:     result.Surname,
			ResumeLink:  result.ResumeLink,
			Tasks:       result.Tasks,
			Solutions:   result.Solutions,
			ChatHistory: result.ChatHistory,
			Metrics:     result.Metrics,
			CreatedAt:   result.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, responses)
}

func (h *VacancyHandler) AddInterviewResult(c *gin.Context) {
	vacancyID := c.Param("vacancy_id")
	if vacancyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty id"})
		return
	}

	var req models.AddInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.AddInterviewResult(c.Request.Context(), vacancyID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
