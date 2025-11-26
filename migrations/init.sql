CREATE TABLE IF NOT EXISTS vacancies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id TEXT NOT NULL,
    profession TEXT NOT NULL,
    position TEXT NOT NULL,
    requirements TEXT[],
    tasks TEXT[],
    task_ideas TEXT[],
    metrics TEXT[],
    is_active BOOLEAN DEFAULT true,
    duration INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS interview_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vacancy_id UUID NOT NULL REFERENCES vacancies(id) ON DELETE CASCADE,
    profession TEXT,
    name TEXT,
    surname TEXT,
    resume_link TEXT,
    tasks TEXT[],
    solutions TEXT[],
    chat_history TEXT[],
    metrics TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_interview_vacancy_id ON interview_results(vacancy_id);