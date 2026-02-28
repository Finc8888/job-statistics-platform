package repository

import (
	"database/sql"
	"job-statistics-api/internal/models"
)

type SkillRepository struct {
	db *sql.DB
}

func NewSkillRepository(db *sql.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

// GetAll возвращает все навыки
func (r *SkillRepository) GetAll() ([]models.Skill, error) {
	query := `SELECT id, name, category, created_at FROM skills ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var s models.Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.Category, &s.CreatedAt); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	return skills, nil
}

// GetByID возвращает навык по ID
func (r *SkillRepository) GetByID(id int) (*models.Skill, error) {
	query := `SELECT id, name, category, created_at FROM skills WHERE id = ?`

	var s models.Skill
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.Name, &s.Category, &s.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Create создает новый навык
func (r *SkillRepository) Create(s *models.Skill) error {
	query := `INSERT INTO skills (name, category) VALUES (?, ?)`

	result, err := r.db.Exec(query, s.Name, s.Category)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = int(id)
	return nil
}

// Update обновляет навык
func (r *SkillRepository) Update(s *models.Skill) error {
	query := `UPDATE skills SET name = ?, category = ? WHERE id = ?`

	_, err := r.db.Exec(query, s.Name, s.Category, s.ID)
	return err
}

// Delete удаляет навык
func (r *SkillRepository) Delete(id int) error {
	query := `DELETE FROM skills WHERE id = ?`

	_, err := r.db.Exec(query, id)
	return err
}
