package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/TOPfiIT/vacancy-service/internal/config"
	"github.com/TOPfiIT/vacancy-service/internal/domain/models"
	"github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func InitPostgres(cfg *config.Config) *PostgresDB {
	db_opts := cfg.Database.GetDatabaseURL()
	db, err := sql.Open("postgres", db_opts)
	if err != nil {
		log.Fatal("Failed to connect to postgres: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping postgres: ", err)
	}

	return &PostgresDB{
		db: db,
	}
}

func (p *PostgresDB) CreateVacancy(ctx context.Context, vacancy *models.Vacancy) error {
	query := `
    INSERT INTO vacancies (id, company_id, profession, position, requirements, tasks, task_ideas, metrics, is_active, duration)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING created_at
    `

	return p.db.QueryRowContext(
		ctx,
		query,
		vacancy.ID,
		vacancy.CompanyID,
		vacancy.Profession,
		vacancy.Position,
		pq.Array(vacancy.Requirements),
		pq.Array(vacancy.Tasks),
		pq.Array(vacancy.TaskIdeas),
		pq.Array(vacancy.Metrics),
		vacancy.IsActive,
		vacancy.Duration,
	).Scan(&vacancy.CreatedAt)
}

func (p *PostgresDB) CreateInterviewResult(ctx context.Context, interview *models.InterviewResult) error {
	query := `
	INSERT INTO interview_results (id, vacancy_id, profession, name, surname, resume_link, tasks, solutions, chat_history, metrics)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING created_at
	`

	return p.db.QueryRowContext(
		ctx,
		query,
		interview.ID,
		interview.VacancyID,
		interview.Profession,
		interview.Name,
		interview.Surname,
		interview.ResumeLink,
		pq.Array(interview.Tasks),
		pq.Array(interview.Solutions),
		pq.Array(interview.ChatHistory),
		pq.Array(interview.Metrics),
	).Scan(&interview.CreatedAt)
}

func (p *PostgresDB) GetVacancyByID(ctx context.Context, vacancyID string) (*models.Vacancy, error) {
	query := `
	SELECT id, company_id, profession, position, requirements, tasks, task_ideas, metrics, is_active, duration, created_at
	FROM vacancies
	WHERE id = $1
	`

	vacancy := &models.Vacancy{}
	err := p.db.QueryRowContext(ctx, query, vacancyID).Scan(
		&vacancy.ID,
		&vacancy.CompanyID,
		&vacancy.Profession,
		&vacancy.Position,
		pq.Array(&vacancy.Requirements),
		pq.Array(&vacancy.Tasks),
		pq.Array(&vacancy.TaskIdeas),
		pq.Array(&vacancy.Metrics),
		&vacancy.IsActive,
		&vacancy.Duration,
		&vacancy.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vacancy not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get vacancy: %w", err)
	}

	return vacancy, nil
}

func (p *PostgresDB) GetVacanciesByCompany(ctx context.Context, companyID string) ([]*models.Vacancy, error) {
	query := `
	SELECT id, company_id, profession, position, requirements, tasks, task_ideas, metrics, is_active, duration, created_at
	FROM vacancies
	WHERE company_id = $1
	ORDER BY created_at DESC
	`

	rows, err := p.db.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query vacancies by company: %w", err)
	}
	defer rows.Close()

	vacancies := []*models.Vacancy{}
	for rows.Next() {
		v := &models.Vacancy{}
		if err := rows.Scan(
			&v.ID,
			&v.CompanyID,
			&v.Profession,
			&v.Position,
			pq.Array(&v.Requirements),
			pq.Array(&v.Tasks),
			pq.Array(&v.TaskIdeas),
			pq.Array(&v.Metrics),
			&v.IsActive,
			&v.Duration,
			&v.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan vacancy: %w", err)
		}

		vacancies = append(vacancies, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return vacancies, nil
}

func (p *PostgresDB) GetInterviewResults(ctx context.Context, vacancyID string) ([]*models.InterviewResult, error) {
	query := `
	SELECT id, vacancy_id, profession, name, surname, resume_link, tasks, solutions, chat_history, metrics, created_at
	FROM interview_results
	WHERE vacancy_id = $1
	ORDER BY created_at DESC
	`

	rows, err := p.db.QueryContext(ctx, query, vacancyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query interviews: %w", err)
	}
	defer rows.Close()

	results := []*models.InterviewResult{}
	for rows.Next() {
		result := &models.InterviewResult{}
		if err := rows.Scan(
			&result.ID,
			&result.VacancyID,
			&result.Profession,
			&result.Name,
			&result.Surname,
			&result.ResumeLink,
			pq.Array(&result.Tasks),
			pq.Array(&result.Solutions),
			pq.Array(&result.ChatHistory),
			pq.Array(&result.Metrics),
			&result.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan interview result: %w", err)
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}
