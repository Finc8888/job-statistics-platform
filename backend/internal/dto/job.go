package dto

import (
	"database/sql"
	"job-statistics-api/internal/models"
	"time"
)

// JobResponse — DTO для отправки вакансии клиенту.
// Все nullable поля представлены указателями — JSON сериализует их как значение или null.
type JobResponse struct {
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

// JobRequest — DTO для приёма данных от клиента при создании/обновлении вакансии.
type JobRequest struct {
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
	IsActive         bool     `json:"is_active"`
}

// JobResponseFromModel конвертирует models.Job → JobResponse
func JobResponseFromModel(j models.Job) JobResponse {
	r := JobResponse{
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
		r.Specialization = &j.Specialization.String
	}
	if j.SalaryMin.Valid {
		r.SalaryMin = &j.SalaryMin.Float64
	}
	if j.SalaryMax.Valid {
		r.SalaryMax = &j.SalaryMax.Float64
	}
	if j.ExperienceYears.Valid {
		r.ExperienceYears = &j.ExperienceYears.String
	}
	if j.Location.Valid {
		r.Location = &j.Location.String
	}
	if j.Description.Valid {
		r.Description = &j.Description.String
	}
	if j.Responsibilities.Valid {
		r.Responsibilities = &j.Responsibilities.String
	}
	if j.Benefits.Valid {
		r.Benefits = &j.Benefits.String
	}
	if j.PostedDate.Valid {
		s := j.PostedDate.Time.Format(time.RFC3339)
		r.PostedDate = &s
	}
	if j.SourceURL.Valid {
		r.SourceURL = &j.SourceURL.String
	}

	return r
}

// JobResponseList конвертирует []models.Job → []JobResponse
func JobResponseList(jobs []models.Job) []JobResponse {
	result := make([]JobResponse, len(jobs))
	for i, j := range jobs {
		result[i] = JobResponseFromModel(j)
	}
	return result
}

// ToModel конвертирует JobRequest → models.Job
func (req JobRequest) ToModel() models.Job {
	j := models.Job{
		CompanyID:       req.CompanyID,
		Title:           req.Title,
		Level:           req.Level,
		SalaryCurrency:  req.SalaryCurrency,
		RemoteAvailable: req.RemoteAvailable,
		IsActive:        req.IsActive,
	}

	if req.Specialization != nil && *req.Specialization != "" {
		j.Specialization = sql.NullString{String: *req.Specialization, Valid: true}
	}
	if req.SalaryMin != nil {
		j.SalaryMin = sql.NullFloat64{Float64: *req.SalaryMin, Valid: true}
	}
	if req.SalaryMax != nil {
		j.SalaryMax = sql.NullFloat64{Float64: *req.SalaryMax, Valid: true}
	}
	if req.ExperienceYears != nil && *req.ExperienceYears != "" {
		j.ExperienceYears = sql.NullString{String: *req.ExperienceYears, Valid: true}
	}
	if req.Location != nil && *req.Location != "" {
		j.Location = sql.NullString{String: *req.Location, Valid: true}
	}
	if req.Description != nil && *req.Description != "" {
		j.Description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.Responsibilities != nil && *req.Responsibilities != "" {
		j.Responsibilities = sql.NullString{String: *req.Responsibilities, Valid: true}
	}
	if req.Benefits != nil && *req.Benefits != "" {
		j.Benefits = sql.NullString{String: *req.Benefits, Valid: true}
	}

	return j
}
