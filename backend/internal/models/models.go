package models

import (
	"database/sql"
	"encoding/json"
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

// MarshalJSON реализует кастомную сериализацию для корректной обработки sql.Null* типов
func (j Job) MarshalJSON() ([]byte, error) {
	type jobJSON struct {
		ID               int      `json:"id"`
		CompanyID        int      `json:"company_id"`
		Title            string   `json:"title"`
		Level            string   `json:"level"`
		Specialization   *string  `json:"specialization"`
		SalaryMin        *float64 `json:"salary_min"`
		SalaryMax        *float64 `json:"salary_max"`
		SalaryCurrency   string   `json:"salary_currency"`
		ExperienceYears  *string  `json:"experience_years"`
		Location         *string  `json:"location"`
		RemoteAvailable  bool     `json:"remote_available"`
		Description      *string  `json:"description"`
		Responsibilities *string  `json:"responsibilities"`
		Benefits         *string  `json:"benefits"`
		PostedDate       *string  `json:"posted_date"`
		IsActive         bool     `json:"is_active"`
		SourceURL        *string  `json:"source_url"`
		CreatedAt        string   `json:"created_at"`
		UpdatedAt        string   `json:"updated_at"`
	}

	out := jobJSON{
		ID:              j.ID,
		CompanyID:       j.CompanyID,
		Title:           j.Title,
		Level:           j.Level,
		SalaryCurrency:  j.SalaryCurrency,
		RemoteAvailable: j.RemoteAvailable,
		IsActive:        j.IsActive,
		CreatedAt:       j.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       j.UpdatedAt.Format(time.RFC3339),
	}

	if j.Specialization.Valid {
		out.Specialization = &j.Specialization.String
	}
	if j.SalaryMin.Valid {
		out.SalaryMin = &j.SalaryMin.Float64
	}
	if j.SalaryMax.Valid {
		out.SalaryMax = &j.SalaryMax.Float64
	}
	if j.ExperienceYears.Valid {
		out.ExperienceYears = &j.ExperienceYears.String
	}
	if j.Location.Valid {
		out.Location = &j.Location.String
	}
	if j.Description.Valid {
		out.Description = &j.Description.String
	}
	if j.Responsibilities.Valid {
		out.Responsibilities = &j.Responsibilities.String
	}
	if j.Benefits.Valid {
		out.Benefits = &j.Benefits.String
	}
	if j.PostedDate.Valid {
		s := j.PostedDate.Time.Format(time.RFC3339)
		out.PostedDate = &s
	}
	if j.SourceURL.Valid {
		out.SourceURL = &j.SourceURL.String
	}

	return json.Marshal(out)
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

// MarshalJSON реализует кастомную сериализацию для Location
func (l Location) MarshalJSON() ([]byte, error) {
	type locationJSON struct {
		ID           int     `json:"id"`
		JobID        int     `json:"job_id"`
		City         string  `json:"city"`
		MetroStation *string `json:"metro_station"`
		IsPrimary    bool    `json:"is_primary"`
	}

	out := locationJSON{
		ID:        l.ID,
		JobID:     l.JobID,
		City:      l.City,
		IsPrimary: l.IsPrimary,
	}

	if l.MetroStation.Valid {
		out.MetroStation = &l.MetroStation.String
	}

	return json.Marshal(out)
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
