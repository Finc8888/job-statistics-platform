package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"job-statistics-api/internal/dto"
	"job-statistics-api/internal/models"
)

// mockJobRepo — тестовая реализация JobRepositoryInterface
type mockJobRepo struct {
	jobs []models.Job
	err  error
}

func (m *mockJobRepo) GetAll() ([]models.Job, error) {
	return m.jobs, m.err
}

func (m *mockJobRepo) GetByID(id int) (*models.Job, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, j := range m.jobs {
		if j.ID == id {
			return &j, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockJobRepo) Create(j *models.Job) error {
	if m.err != nil {
		return m.err
	}
	j.ID = 55
	return nil
}

func (m *mockJobRepo) Update(j *models.Job) error {
	return m.err
}

func (m *mockJobRepo) Delete(id int) error {
	return m.err
}

func makeJobRouter(repo *mockJobRepo) *mux.Router {
	h := NewJobHandler(repo)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/jobs", h.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/jobs/{id}", h.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/jobs", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/jobs/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/jobs/{id}", h.Delete).Methods(http.MethodDelete)
	return r
}

func TestJobHandler_GetAll(t *testing.T) {
	tests := []struct {
		name       string
		repo       *mockJobRepo
		wantStatus int
		wantLen    int
	}{
		{
			name: "returns jobs",
			repo: &mockJobRepo{
				jobs: []models.Job{
					{ID: 1, Title: "Go Developer", Level: "Middle"},
					{ID: 2, Title: "Python Developer", Level: "Senior"},
				},
			},
			wantStatus: http.StatusOK,
			wantLen:    2,
		},
		{
			name:       "empty list",
			repo:       &mockJobRepo{jobs: []models.Job{}},
			wantStatus: http.StatusOK,
			wantLen:    0,
		},
		{
			name:       "db error returns 500",
			repo:       &mockJobRepo{err: errors.New("db down")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
			rec := httptest.NewRecorder()

			makeJobRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				var jobs []dto.JobResponse
				if err := json.NewDecoder(rec.Body).Decode(&jobs); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(jobs) != tt.wantLen {
					t.Errorf("got %d jobs, want %d", len(jobs), tt.wantLen)
				}
			}
		})
	}
}

func TestJobHandler_GetByID(t *testing.T) {
	sampleJobs := []models.Job{
		{ID: 1, Title: "Go Developer", Level: "Middle"},
	}

	tests := []struct {
		name       string
		url        string
		repo       *mockJobRepo
		wantStatus int
		wantTitle  string
	}{
		{
			name:       "found job",
			url:        "/api/v1/jobs/1",
			repo:       &mockJobRepo{jobs: sampleJobs},
			wantStatus: http.StatusOK,
			wantTitle:  "Go Developer",
		},
		{
			name:       "not found returns 404",
			url:        "/api/v1/jobs/999",
			repo:       &mockJobRepo{jobs: sampleJobs},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id returns 400",
			url:        "/api/v1/jobs/abc",
			repo:       &mockJobRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "db error returns 500",
			url:        "/api/v1/jobs/1",
			repo:       &mockJobRepo{err: errors.New("timeout")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			makeJobRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK && tt.wantTitle != "" {
				var job dto.JobResponse
				if err := json.NewDecoder(rec.Body).Decode(&job); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if job.Title != tt.wantTitle {
					t.Errorf("title = %q, want %q", job.Title, tt.wantTitle)
				}
			}
		})
	}
}

func TestJobHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		repo       *mockJobRepo
		wantStatus int
		wantID     int
	}{
		{
			name: "create successfully",
			body: dto.JobRequest{
				CompanyID: 1,
				Title:     "Rust Developer",
				Level:     "Senior",
			},
			repo:       &mockJobRepo{},
			wantStatus: http.StatusCreated,
			wantID:     55,
		},
		{
			name:       "invalid json returns 400",
			body:       "not-json",
			repo:       &mockJobRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "db error returns 500",
			body:       dto.JobRequest{Title: "Broken Job"},
			repo:       &mockJobRepo{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			makeJobRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				var job dto.JobResponse
				if err := json.NewDecoder(rec.Body).Decode(&job); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if job.ID != tt.wantID {
					t.Errorf("ID = %d, want %d", job.ID, tt.wantID)
				}
			}
		})
	}
}

func TestJobHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		body       interface{}
		repo       *mockJobRepo
		wantStatus int
	}{
		{
			name:       "update successfully",
			url:        "/api/v1/jobs/1",
			body:       dto.JobRequest{Title: "Senior Go Dev", Level: "Senior", CompanyID: 1},
			repo:       &mockJobRepo{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id returns 400",
			url:        "/api/v1/jobs/abc",
			body:       dto.JobRequest{},
			repo:       &mockJobRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "db error returns 500",
			url:        "/api/v1/jobs/1",
			body:       dto.JobRequest{Title: "Broken"},
			repo:       &mockJobRepo{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, tt.url, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			makeJobRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestJobHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		repo       *mockJobRepo
		wantStatus int
	}{
		{
			name:       "delete successfully",
			url:        "/api/v1/jobs/1",
			repo:       &mockJobRepo{},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id returns 400",
			url:        "/api/v1/jobs/abc",
			repo:       &mockJobRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "db error returns 500",
			url:        "/api/v1/jobs/1",
			repo:       &mockJobRepo{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			rec := httptest.NewRecorder()

			makeJobRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
