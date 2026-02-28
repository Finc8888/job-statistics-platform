package handlers

import (
	"encoding/json"
	"job-statistics-api/internal/repository"
	"net/http"
	"strconv"
)

type StatsHandler struct {
	repo *repository.StatsRepository
}

func NewStatsHandler(repo *repository.StatsRepository) *StatsHandler {
	return &StatsHandler{repo: repo}
}

// GetTopSkills получает топ навыков
func (h *StatsHandler) GetTopSkills(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // по умолчанию
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	skills, err := h.repo.GetTopSkills(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}

// GetSkillSalaries получает зарплаты по навыкам
func (h *StatsHandler) GetSkillSalaries(w http.ResponseWriter, r *http.Request) {
	minVacStr := r.URL.Query().Get("min_vacancies")
	minVac := 1 // по умолчанию
	if minVacStr != "" {
		if mv, err := strconv.Atoi(minVacStr); err == nil && mv > 0 {
			minVac = mv
		}
	}

	salaries, err := h.repo.GetSkillSalaries(minVac)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(salaries)
}

// GetSkillsByLevel получает навыки по уровню
func (h *StatsHandler) GetSkillsByLevel(w http.ResponseWriter, r *http.Request) {
	skills, err := h.repo.GetSkillsByLevel()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}

// GetCompanyStats получает статистику по компаниям
func (h *StatsHandler) GetCompanyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.GetCompanyStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetDatabaseStats получает статистику по базам данных
func (h *StatsHandler) GetDatabaseStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.GetDatabaseStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetLanguageStats получает статистику по языкам программирования
func (h *StatsHandler) GetLanguageStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.repo.GetLanguageStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
