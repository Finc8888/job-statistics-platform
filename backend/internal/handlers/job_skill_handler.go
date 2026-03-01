package handlers

import (
	"encoding/json"
	"job-statistics-api/internal/repository"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type JobSkillHandler struct {
	repo repository.JobSkillRepositoryInterface
}

func NewJobSkillHandler(repo repository.JobSkillRepositoryInterface) *JobSkillHandler {
	return &JobSkillHandler{repo: repo}
}

// GetByJobID returns all skills for a given job: GET /api/v1/jobs/{id}/skills
func (h *JobSkillHandler) GetByJobID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Job ID", http.StatusBadRequest)
		return
	}

	skills, err := h.repo.GetSkillsByJobID(jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}

type setJobSkillsRequest struct {
	SkillIDs []int `json:"skill_ids"`
}

// SetJobSkills replaces skills for a job: POST /api/v1/jobs/{id}/skills
func (h *JobSkillHandler) SetJobSkills(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Job ID", http.StatusBadRequest)
		return
	}

	var req setJobSkillsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.SkillIDs == nil {
		req.SkillIDs = []int{}
	}

	if err := h.repo.SetJobSkills(jobID, req.SkillIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	skills, err := h.repo.GetSkillsByJobID(jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}
