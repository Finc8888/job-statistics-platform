package repository

import (
	"database/sql"
	"job-statistics-api/internal/models"
)

type LocationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

// GetAll возвращает все локации
func (r *LocationRepository) GetAll() ([]models.Location, error) {
	query := `SELECT id, job_id, city, metro_station, is_primary FROM locations`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var l models.Location
		if err := rows.Scan(&l.ID, &l.JobID, &l.City, &l.MetroStation, &l.IsPrimary); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}

	return locations, nil
}

// GetByJobID возвращает локации для вакансии
func (r *LocationRepository) GetByJobID(jobID int) ([]models.Location, error) {
	query := `SELECT id, job_id, city, metro_station, is_primary FROM locations WHERE job_id = ?`

	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var l models.Location
		if err := rows.Scan(&l.ID, &l.JobID, &l.City, &l.MetroStation, &l.IsPrimary); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}

	return locations, nil
}

// Create создает новую локацию
func (r *LocationRepository) Create(l *models.Location) error {
	query := `INSERT INTO locations (job_id, city, metro_station, is_primary) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, l.JobID, l.City, l.MetroStation, l.IsPrimary)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	l.ID = int(id)
	return nil
}

// Update обновляет локацию
func (r *LocationRepository) Update(l *models.Location) error {
	query := `UPDATE locations SET job_id = ?, city = ?, metro_station = ?, is_primary = ? WHERE id = ?`

	_, err := r.db.Exec(query, l.JobID, l.City, l.MetroStation, l.IsPrimary, l.ID)
	return err
}

// Delete удаляет локацию
func (r *LocationRepository) Delete(id int) error {
	query := `DELETE FROM locations WHERE id = ?`

	_, err := r.db.Exec(query, id)
	return err
}
