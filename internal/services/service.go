package services

import (
	"context"
	"fmt"

	"github.com/TOPfiIT/vacancy-service/internal/db"
	"github.com/TOPfiIT/vacancy-service/internal/domain/models"
	"github.com/google/uuid"
)

type VacancyService struct {
	repo *db.PostgresDB
}

func NewVacancyService(db *db.PostgresDB) *VacancyService {
	return &VacancyService{repo: db}
}

func (s *VacancyService) CreateVacancy(ctx context.Context, req *models.CreateVacancyRequest) (*models.Vacancy, error) {
	vacancy := &models.Vacancy{
		ID:           uuid.New().String(),
		CompanyID:    req.CompanyID,
		Profession:   req.Profession,
		Position:     req.Position,
		Requirements: req.Requirements,
		Tasks:        req.Tasks,
		TaskIdeas:    req.TaskIdeas,
		Metrics:      req.Metrics,
		IsActive:     req.IsActive,
		Duration:     req.Duration,
	}

	if err := s.repo.CreateVacancy(ctx, vacancy); err != nil {
		return nil, fmt.Errorf("failed to create vacancy: %w", err)
	}

	return vacancy, nil
}

func (s *VacancyService) GetVacancy(ctx context.Context, vacancyID string) (*models.Vacancy, error) {
	vacancy, err := s.repo.GetVacancyByID(ctx, vacancyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vacancy: %w", err)
	}

	return vacancy, nil
}

func (s *VacancyService) GetVacanciesByCompany(ctx context.Context, companyID string) ([]*models.Vacancy, error) {
	vacancies, err := s.repo.GetVacanciesByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vacancies: %w", err)
	}

	return vacancies, nil
}

func (s *VacancyService) GetInterviewResults(ctx context.Context, vacancyID string) ([]*models.InterviewResult, error) {
	results, err := s.repo.GetInterviewResults(ctx, vacancyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get interview results: %w", err)
	}

	return results, nil
}

func (s *VacancyService) AddInterviewResult(ctx context.Context, vacancyID string, req *models.AddInterviewRequest) (*models.InterviewResult, error) {
	if _, err := s.repo.GetVacancyByID(ctx, vacancyID); err != nil {
		return nil, fmt.Errorf("invalid vacancy ID: %w", err)
	}

	result := &models.InterviewResult{
		ID:          uuid.New().String(),
		VacancyID:   vacancyID,
		Profession:  req.Profession,
		Name:        req.Name,
		Surname:     req.Surname,
		ResumeLink:  req.ResumeLink,
		Tasks:       req.Tasks,
		Solutions:   req.Solutions,
		ChatHistory: req.ChatHistory,
		Metrics:     req.Metrics,
	}

	if err := s.repo.CreateInterviewResult(ctx, result); err != nil {
		return nil, fmt.Errorf("failed to add interview result: %w", err)
	}

	return result, nil
}
