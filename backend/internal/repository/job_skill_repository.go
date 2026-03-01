package repository

import (
	"database/sql"
	"fmt"
	"job-statistics-api/internal/models"
	"strings"
)

type JobSkillRepository struct {
	db *sql.DB
}

func NewJobSkillRepository(db *sql.DB) *JobSkillRepository {
	return &JobSkillRepository{db: db}
}

// GetSkillsByJobID returns all skills associated with a given job
func (r *JobSkillRepository) GetSkillsByJobID(jobID int) ([]models.Skill, error) {
	query := `
		SELECT s.id, s.name, s.category, s.created_at
		FROM skills s
		INNER JOIN job_skills js ON s.id = js.skill_id
		WHERE js.job_id = ?
		ORDER BY s.name
	`
	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skills := []models.Skill{}
	for rows.Next() {
		var s models.Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.Category, &s.CreatedAt); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, nil
}

// SetJobSkills replaces all skills for a job (delete existing + insert new).
// Passing an empty slice removes all skill associations.
func (r *JobSkillRepository) SetJobSkills(jobID int, skillIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM job_skills WHERE job_id = ?`, jobID); err != nil {
		return err
	}

	if len(skillIDs) > 0 {
		placeholders := make([]string, len(skillIDs))
		args := make([]interface{}, 0, len(skillIDs)*2)
		for i, skillID := range skillIDs {
			placeholders[i] = "(?, ?, 1, 0)"
			args = append(args, jobID, skillID)
		}
		query := fmt.Sprintf(
			`INSERT INTO job_skills (job_id, skill_id, is_required, is_nice_to_have) VALUES %s`,
			strings.Join(placeholders, ", "),
		)
		if _, err := tx.Exec(query, args...); err != nil {
			return err
		}
	}

	return tx.Commit()
}
