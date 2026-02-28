package repository

import (
	"database/sql"
	"job-statistics-api/internal/models"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

// GetAll возвращает все вакансии
func (r *JobRepository) GetAll() ([]models.Job, error) {
	query := `SELECT id, company_id, title, level, specialization, salary_min, salary_max,
		salary_currency, experience_years, location, remote_available, description,
		responsibilities, benefits, posted_date, is_active, source_url, created_at, updated_at
		FROM jobs ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var j models.Job
		if err := rows.Scan(&j.ID, &j.CompanyID, &j.Title, &j.Level, &j.Specialization,
			&j.SalaryMin, &j.SalaryMax, &j.SalaryCurrency, &j.ExperienceYears, &j.Location,
			&j.RemoteAvailable, &j.Description, &j.Responsibilities, &j.Benefits,
			&j.PostedDate, &j.IsActive, &j.SourceURL, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}

	return jobs, nil
}

// GetByID возвращает вакансию по ID
func (r *JobRepository) GetByID(id int) (*models.Job, error) {
	query := `SELECT id, company_id, title, level, specialization, salary_min, salary_max,
		salary_currency, experience_years, location, remote_available, description,
		responsibilities, benefits, posted_date, is_active, source_url, created_at, updated_at
		FROM jobs WHERE id = ?`

	var j models.Job
	err := r.db.QueryRow(query, id).Scan(&j.ID, &j.CompanyID, &j.Title, &j.Level, &j.Specialization,
		&j.SalaryMin, &j.SalaryMax, &j.SalaryCurrency, &j.ExperienceYears, &j.Location,
		&j.RemoteAvailable, &j.Description, &j.Responsibilities, &j.Benefits,
		&j.PostedDate, &j.IsActive, &j.SourceURL, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

// Create создает новую вакансию
func (r *JobRepository) Create(j *models.Job) error {
	query := `INSERT INTO jobs (company_id, title, level, specialization, salary_min, salary_max,
		salary_currency, experience_years, location, remote_available, description,
		responsibilities, benefits, posted_date, is_active, source_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, j.CompanyID, j.Title, j.Level, j.Specialization,
		j.SalaryMin, j.SalaryMax, j.SalaryCurrency, j.ExperienceYears, j.Location,
		j.RemoteAvailable, j.Description, j.Responsibilities, j.Benefits,
		j.PostedDate, j.IsActive, j.SourceURL)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	j.ID = int(id)
	return nil
}

// Update обновляет вакансию
func (r *JobRepository) Update(j *models.Job) error {
	query := `UPDATE jobs SET company_id = ?, title = ?, level = ?, specialization = ?,
		salary_min = ?, salary_max = ?, salary_currency = ?, experience_years = ?,
		location = ?, remote_available = ?, description = ?, responsibilities = ?,
		benefits = ?, posted_date = ?, is_active = ?, source_url = ? WHERE id = ?`

	_, err := r.db.Exec(query, j.CompanyID, j.Title, j.Level, j.Specialization,
		j.SalaryMin, j.SalaryMax, j.SalaryCurrency, j.ExperienceYears, j.Location,
		j.RemoteAvailable, j.Description, j.Responsibilities, j.Benefits,
		j.PostedDate, j.IsActive, j.SourceURL, j.ID)
	return err
}

// Delete удаляет вакансию
func (r *JobRepository) Delete(id int) error {
	query := `DELETE FROM jobs WHERE id = ?`

	_, err := r.db.Exec(query, id)
	return err
}
