package models

import (
	"database/sql"
	"time"
)

// Company представляет компанию
type Company struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Job представляет вакансию
type Job struct {
	ID               int            `json:"id"`
	CompanyID        int            `json:"company_id"`
	Title            string         `json:"title"`
	Level            string         `json:"level"`
	Specialization   sql.NullString `json:"specialization"`
	SalaryMin        sql.NullFloat64 `json:"salary_min"`
	SalaryMax        sql.NullFloat64 `json:"salary_max"`
	SalaryCurrency   string         `json:"salary_currency"`
	ExperienceYears  sql.NullString `json:"experience_years"`
	Location         sql.NullString `json:"location"`
	RemoteAvailable  bool           `json:"remote_available"`
	Description      sql.NullString `json:"description"`
	Responsibilities sql.NullString `json:"responsibilities"`
	Benefits         sql.NullString `json:"benefits"`
	PostedDate       sql.NullTime   `json:"posted_date"`
	IsActive         bool           `json:"is_active"`
	SourceURL        sql.NullString `json:"source_url"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// Skill представляет навык
type Skill struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// Location представляет локацию
type Location struct {
	ID           int            `json:"id"`
	JobID        int            `json:"job_id"`
	City         string         `json:"city"`
	MetroStation sql.NullString `json:"metro_station"`
	IsPrimary    bool           `json:"is_primary"`
}

// JobSkill представляет связь вакансии и навыка
type JobSkill struct {
	ID           int       `json:"id"`
	JobID        int       `json:"job_id"`
	SkillID      int       `json:"skill_id"`
	IsRequired   bool      `json:"is_required"`
	IsNiceToHave bool      `json:"is_nice_to_have"`
	CreatedAt    time.Time `json:"created_at"`
}

// Статистические модели

// TopSkill - топ навыков
type TopSkill struct {
	Name            string `json:"name"`
	Category        string `json:"category"`
	VacancyCount    int    `json:"vacancy_count"`
	RequiredCount   int    `json:"required_count"`
	NiceToHaveCount int    `json:"nice_to_have_count"`
}

// SkillSalary - средняя зарплата по навыкам
type SkillSalary struct {
	Name         string  `json:"name"`
	AvgSalary    float64 `json:"avg_salary"`
	VacancyCount int     `json:"vacancy_count"`
}

// SkillByLevel - навыки по уровню
type SkillByLevel struct {
	Level string `json:"level"`
	Name  string `json:"skill_name"`
	Count int    `json:"count"`
}

// CompanyStats - статистика по компаниям
type CompanyStats struct {
	Company         string  `json:"company"`
	VacanciesCount  int     `json:"vacancies_count"`
	Levels          string  `json:"levels"`
	MinSalary       float64 `json:"min_salary"`
	MaxSalary       float64 `json:"max_salary"`
	LocationsCount  int     `json:"locations_count"`
	RemoteVacancies int     `json:"remote_vacancies"`
}

// DatabaseStats - статистика по БД
type DatabaseStats struct {
	Database     string  `json:"database"`
	Vacancies    int     `json:"vacancies"`
	Required     int     `json:"required"`
	NiceToHave   int     `json:"nice_to_have"`
	Companies    string  `json:"companies"`
	AvgSalary    float64 `json:"avg_salary"`
}

// LanguageStats - статистика по языкам программирования
type LanguageStats struct {
	Language      string  `json:"language"`
	VacancyCount  int     `json:"vacancy_count"`
	Companies     string  `json:"companies"`
	Levels        string  `json:"levels"`
	AvgSalaryRub  float64 `json:"avg_salary_rub"`
}
