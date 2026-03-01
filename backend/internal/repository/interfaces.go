package repository

import "job-statistics-api/internal/models"

// JobRepositoryInterface определяет контракт для работы с вакансиями.
// Позволяет подменять реализацию в тестах mock-объектами.
type JobRepositoryInterface interface {
	GetAll() ([]models.Job, error)
	GetByID(id int) (*models.Job, error)
	Create(j *models.Job) error
	Update(j *models.Job) error
	Delete(id int) error
}

// CompanyRepositoryInterface определяет контракт для работы с компаниями.
type CompanyRepositoryInterface interface {
	GetAll() ([]models.Company, error)
	GetByID(id int) (*models.Company, error)
	Create(c *models.Company) error
	Update(c *models.Company) error
	Delete(id int) error
}

// JobSkillRepositoryInterface определяет контракт для работы с навыками вакансии.
type JobSkillRepositoryInterface interface {
	GetSkillsByJobID(jobID int) ([]models.Skill, error)
	SetJobSkills(jobID int, skillIDs []int) error
}
