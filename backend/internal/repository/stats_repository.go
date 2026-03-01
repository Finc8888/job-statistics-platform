package repository

import (
	"database/sql"
	"job-statistics-api/internal/models"
)

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetTopSkills - топ-10 самых востребованных навыков
func (r *StatsRepository) GetTopSkills(limit int) ([]models.TopSkill, error) {
	query := `
		SELECT
			s.name,
			s.category,
			COUNT(*) as vacancy_count,
			SUM(CASE WHEN js.is_required THEN 1 ELSE 0 END) as required_count,
			SUM(CASE WHEN js.is_nice_to_have THEN 1 ELSE 0 END) as nice_to_have_count
		FROM skills s
		JOIN job_skills js ON s.id = js.skill_id
		JOIN jobs j ON js.job_id = j.id
		WHERE j.is_active = TRUE
		GROUP BY s.id, s.name, s.category
		ORDER BY vacancy_count DESC
		LIMIT ?`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.TopSkill
	for rows.Next() {
		var s models.TopSkill
		if err := rows.Scan(&s.Name, &s.Category, &s.VacancyCount, &s.RequiredCount, &s.NiceToHaveCount); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	return skills, nil
}

// GetSkillSalaries - средняя зарплата по навыкам
func (r *StatsRepository) GetSkillSalaries(minVacancies int) ([]models.SkillSalary, error) {
	query := `
		SELECT
			s.name,
			COALESCE(AVG((j.salary_min + COALESCE(j.salary_max, j.salary_min)) / 2), 0) as avg_salary,
			COUNT(*) as vacancy_count
		FROM skills s
		JOIN job_skills js ON s.id = js.skill_id
		JOIN jobs j ON js.job_id = j.id
		WHERE j.is_active = TRUE
			AND j.salary_min IS NOT NULL
		GROUP BY s.id, s.name
		HAVING vacancy_count >= ?
		ORDER BY avg_salary DESC`

	rows, err := r.db.Query(query, minVacancies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var salaries []models.SkillSalary
	for rows.Next() {
		var ss models.SkillSalary
		if err := rows.Scan(&ss.Name, &ss.AvgSalary, &ss.VacancyCount); err != nil {
			return nil, err
		}
		salaries = append(salaries, ss)
	}

	return salaries, nil
}

// GetSkillsByLevel - востребованность навыков по уровню
func (r *StatsRepository) GetSkillsByLevel() ([]models.SkillByLevel, error) {
	query := `
		SELECT
			j.level,
			s.name,
			COUNT(*) as count
		FROM jobs j
		JOIN job_skills js ON j.id = js.job_id
		JOIN skills s ON js.skill_id = s.id
		WHERE j.is_active = TRUE
		GROUP BY j.level, s.name
		ORDER BY j.level, count DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SkillByLevel
	for rows.Next() {
		var sbl models.SkillByLevel
		if err := rows.Scan(&sbl.Level, &sbl.Name, &sbl.Count); err != nil {
			return nil, err
		}
		results = append(results, sbl)
	}

	return results, nil
}

// GetCompanyStats - статистика по компаниям
func (r *StatsRepository) GetCompanyStats() ([]models.CompanyStats, error) {
	query := `
		SELECT
			c.name as company,
			COUNT(j.id) as vacancies_count,
			GROUP_CONCAT(DISTINCT j.level ORDER BY j.level SEPARATOR ', ') as levels,
			COALESCE(MIN(j.salary_min), 0) as min_salary,
			COALESCE(MAX(COALESCE(j.salary_max, j.salary_min)), 0) as max_salary,
			COUNT(DISTINCT l.city) as locations_count,
			SUM(CASE WHEN j.remote_available THEN 1 ELSE 0 END) as remote_vacancies
		FROM companies c
		LEFT JOIN jobs j ON c.id = j.company_id
		LEFT JOIN locations l ON j.id = l.job_id
		WHERE j.is_active = TRUE
		GROUP BY c.id, c.name
		ORDER BY vacancies_count DESC, max_salary DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.CompanyStats
	for rows.Next() {
		var cs models.CompanyStats
		if err := rows.Scan(&cs.Company, &cs.VacanciesCount, &cs.Levels, &cs.MinSalary, &cs.MaxSalary, &cs.LocationsCount, &cs.RemoteVacancies); err != nil {
			return nil, err
		}
		stats = append(stats, cs)
	}

	return stats, nil
}

// GetDatabaseStats - топ баз данных
func (r *StatsRepository) GetDatabaseStats() ([]models.DatabaseStats, error) {
	query := `
		SELECT
			s.name as database_name,
			COUNT(DISTINCT j.id) as vacancies,
			SUM(CASE WHEN js.is_required THEN 1 ELSE 0 END) as required,
			SUM(CASE WHEN js.is_nice_to_have THEN 1 ELSE 0 END) as nice_to_have,
			GROUP_CONCAT(DISTINCT c.name ORDER BY c.name SEPARATOR ', ') as companies,
			COALESCE(AVG((j.salary_min + COALESCE(j.salary_max, j.salary_min)) / 2), 0) as avg_salary
		FROM skills s
		JOIN job_skills js ON s.id = js.skill_id
		JOIN jobs j ON js.job_id = j.id
		JOIN companies c ON j.company_id = c.id
		WHERE s.category = 'База данных'
			AND j.is_active = TRUE
		GROUP BY s.id, s.name
		ORDER BY vacancies DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.DatabaseStats
	for rows.Next() {
		var ds models.DatabaseStats
		if err := rows.Scan(&ds.Database, &ds.Vacancies, &ds.Required, &ds.NiceToHave, &ds.Companies, &ds.AvgSalary); err != nil {
			return nil, err
		}
		stats = append(stats, ds)
	}

	return stats, nil
}

// GetLanguageStats - статистика по языкам программирования
func (r *StatsRepository) GetLanguageStats() ([]models.LanguageStats, error) {
	query := `
		SELECT
			s.name as language,
			COUNT(DISTINCT j.id) as vacancy_count,
			GROUP_CONCAT(DISTINCT c.name ORDER BY c.name SEPARATOR ', ') as companies,
			GROUP_CONCAT(DISTINCT j.level ORDER BY j.level SEPARATOR ', ') as levels,
			COALESCE(ROUND(AVG((j.salary_min + COALESCE(j.salary_max, j.salary_min)) / 2), 0), 0) as avg_salary_rub
		FROM skills s
		JOIN job_skills js ON s.id = js.skill_id
		JOIN jobs j ON js.job_id = j.id
		JOIN companies c ON j.company_id = c.id
		WHERE s.category = 'Язык программирования'
			AND j.is_active = TRUE
			AND j.salary_currency = 'RUB'
		GROUP BY s.id, s.name
		ORDER BY vacancy_count DESC, avg_salary_rub DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.LanguageStats
	for rows.Next() {
		var ls models.LanguageStats
		if err := rows.Scan(&ls.Language, &ls.VacancyCount, &ls.Companies, &ls.Levels, &ls.AvgSalaryRub); err != nil {
			return nil, err
		}
		stats = append(stats, ls)
	}

	return stats, nil
}
