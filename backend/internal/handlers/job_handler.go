package handlers

import (
	"database/sql"
	"encoding/json"
	"job-statistics-api/internal/dto"
	"job-statistics-api/internal/repository"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type JobHandler struct {
	repo repository.JobRepositoryInterface
}

func NewJobHandler(repo repository.JobRepositoryInterface) *JobHandler {
	return &JobHandler{repo: repo}
}

// GetAll получает все вакансии
func (h *JobHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.JobResponseList(jobs))
}

// GetByID получает вакансию по ID
func (h *JobHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	job, err := h.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.JobResponseFromModel(*job))
}

// Create создает новую вакансию
func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	job := req.ToModel()
	if err := h.repo.Create(&job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.JobResponseFromModel(job))
}

// Update обновляет вакансию
func (h *JobHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req dto.JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	job := req.ToModel()
	job.ID = id
	if err := h.repo.Update(&job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.JobResponseFromModel(job))
}

// Delete удаляет вакансию
func (h *JobHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
