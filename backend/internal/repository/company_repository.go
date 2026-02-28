package repository

import (
	"database/sql"
	"job-statistics-api/internal/models"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// GetAll возвращает все компании
func (r *CompanyRepository) GetAll() ([]models.Company, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM companies ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}

	return companies, nil
}

// GetByID возвращает компанию по ID
func (r *CompanyRepository) GetByID(id int) (*models.Company, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM companies WHERE id = ?`

	var c models.Company
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Create создает новую компанию
func (r *CompanyRepository) Create(c *models.Company) error {
	query := `INSERT INTO companies (name, description) VALUES (?, ?)`

	result, err := r.db.Exec(query, c.Name, c.Description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = int(id)
	return nil
}

// Update обновляет компанию
func (r *CompanyRepository) Update(c *models.Company) error {
	query := `UPDATE companies SET name = ?, description = ? WHERE id = ?`

	_, err := r.db.Exec(query, c.Name, c.Description, c.ID)
	return err
}

// Delete удаляет компанию
func (r *CompanyRepository) Delete(id int) error {
	query := `DELETE FROM companies WHERE id = ?`

	_, err := r.db.Exec(query, id)
	return err
}
