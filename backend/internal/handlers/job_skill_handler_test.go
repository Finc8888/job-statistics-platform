package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"job-statistics-api/internal/models"
)

// mockJobSkillRepo — тестовая реализация JobSkillRepositoryInterface
type mockJobSkillRepo struct {
	skills []models.Skill
	getErr error
	setErr error
}

func (m *mockJobSkillRepo) GetSkillsByJobID(_ int) ([]models.Skill, error) {
	return m.skills, m.getErr
}

func (m *mockJobSkillRepo) SetJobSkills(_ int, _ []int) error {
	return m.setErr
}

func makeJobSkillRouter(repo *mockJobSkillRepo) *mux.Router {
	h := NewJobSkillHandler(repo)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/jobs/{id}/skills", h.GetByJobID).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/jobs/{id}/skills", h.SetJobSkills).Methods(http.MethodPost)
	return r
}

func sampleSkills() []models.Skill {
	return []models.Skill{
		{ID: 10, Name: "Go", Category: "Язык программирования", CreatedAt: time.Now()},
		{ID: 11, Name: "PostgreSQL", Category: "База данных", CreatedAt: time.Now()},
	}
}

func TestJobSkillHandler_GetByJobID(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		repo       *mockJobSkillRepo
		wantStatus int
		wantLen    int
	}{
		{
			name:       "returns skills for a job",
			url:        "/api/v1/jobs/1/skills",
			repo:       &mockJobSkillRepo{skills: sampleSkills()},
			wantStatus: http.StatusOK,
			wantLen:    2,
		},
		{
			name:       "returns empty array when no skills",
			url:        "/api/v1/jobs/2/skills",
			repo:       &mockJobSkillRepo{skills: []models.Skill{}},
			wantStatus: http.StatusOK,
			wantLen:    0,
		},
		{
			name:       "invalid id returns 400",
			url:        "/api/v1/jobs/abc/skills",
			repo:       &mockJobSkillRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "database error returns 500",
			url:        "/api/v1/jobs/1/skills",
			repo:       &mockJobSkillRepo{getErr: errors.New("db down")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			makeJobSkillRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				var skills []models.Skill
				if err := json.NewDecoder(rec.Body).Decode(&skills); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(skills) != tt.wantLen {
					t.Errorf("got %d skills, want %d", len(skills), tt.wantLen)
				}
			}
		})
	}
}

func TestJobSkillHandler_SetJobSkills(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		body       interface{}
		repo       *mockJobSkillRepo
		wantStatus int
		wantLen    int
	}{
		{
			name:       "set skills successfully",
			url:        "/api/v1/jobs/1/skills",
			body:       map[string]interface{}{"skill_ids": []int{10, 11}},
			repo:       &mockJobSkillRepo{skills: sampleSkills()},
			wantStatus: http.StatusOK,
			wantLen:    2,
		},
		{
			name:       "clear all skills with empty array",
			url:        "/api/v1/jobs/1/skills",
			body:       map[string]interface{}{"skill_ids": []int{}},
			repo:       &mockJobSkillRepo{skills: []models.Skill{}},
			wantStatus: http.StatusOK,
			wantLen:    0,
		},
		{
			name:       "invalid job id returns 400",
			url:        "/api/v1/jobs/abc/skills",
			body:       map[string]interface{}{"skill_ids": []int{1}},
			repo:       &mockJobSkillRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json body returns 400",
			url:        "/api/v1/jobs/1/skills",
			body:       "not-json",
			repo:       &mockJobSkillRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "set error returns 500",
			url:        "/api/v1/jobs/1/skills",
			body:       map[string]interface{}{"skill_ids": []int{1}},
			repo:       &mockJobSkillRepo{setErr: errors.New("constraint violation")},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "get after set error returns 500",
			url:  "/api/v1/jobs/1/skills",
			body: map[string]interface{}{"skill_ids": []int{1}},
			repo: &mockJobSkillRepo{
				setErr: nil,
				getErr: errors.New("db down"),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			makeJobSkillRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				var skills []models.Skill
				if err := json.NewDecoder(rec.Body).Decode(&skills); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(skills) != tt.wantLen {
					t.Errorf("got %d skills, want %d", len(skills), tt.wantLen)
				}
			}
		})
	}
}
