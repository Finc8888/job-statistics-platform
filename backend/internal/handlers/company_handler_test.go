package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"job-statistics-api/internal/models"
)

// mockCompanyRepo — тестовая реализация CompanyRepositoryInterface
type mockCompanyRepo struct {
	companies []models.Company
	err       error
}

func (m *mockCompanyRepo) GetAll() ([]models.Company, error) {
	return m.companies, m.err
}

func (m *mockCompanyRepo) GetByID(id int) (*models.Company, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, c := range m.companies {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockCompanyRepo) Create(c *models.Company) error {
	if m.err != nil {
		return m.err
	}
	c.ID = 99
	return nil
}

func (m *mockCompanyRepo) Update(c *models.Company) error {
	return m.err
}

func (m *mockCompanyRepo) Delete(id int) error {
	return m.err
}

// helpers

func now() time.Time { return time.Now() }

func makeCompanyRouter(repo *mockCompanyRepo) *mux.Router {
	h := NewCompanyHandler(repo)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/companies", h.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/companies/{id}", h.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/companies", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/companies/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/companies/{id}", h.Delete).Methods(http.MethodDelete)
	return r
}

func TestCompanyHandler_GetAll(t *testing.T) {
	tests := []struct {
		name       string
		repo       *mockCompanyRepo
		wantStatus int
		wantLen    int
	}{
		{
			name: "returns companies",
			repo: &mockCompanyRepo{
				companies: []models.Company{
					{ID: 1, Name: "Yandex", CreatedAt: now(), UpdatedAt: now()},
					{ID: 2, Name: "Sber", CreatedAt: now(), UpdatedAt: now()},
				},
			},
			wantStatus: http.StatusOK,
			wantLen:    2,
		},
		{
			name:       "empty list",
			repo:       &mockCompanyRepo{companies: []models.Company{}},
			wantStatus: http.StatusOK,
			wantLen:    0,
		},
		{
			name:       "db error returns 500",
			repo:       &mockCompanyRepo{err: errors.New("db down")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/companies", nil)
			rec := httptest.NewRecorder()

			makeCompanyRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				var companies []models.Company
				if err := json.NewDecoder(rec.Body).Decode(&companies); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(companies) != tt.wantLen {
					t.Errorf("got %d companies, want %d", len(companies), tt.wantLen)
				}
			}
		})
	}
}

func TestCompanyHandler_GetByID(t *testing.T) {
	sampleCompanies := []models.Company{
		{ID: 1, Name: "Yandex", CreatedAt: now(), UpdatedAt: now()},
	}

	tests := []struct {
		name       string
		url        string
		repo       *mockCompanyRepo
		wantStatus int
		wantName   string
	}{
		{
			name:       "found company",
			url:        "/api/v1/companies/1",
			repo:       &mockCompanyRepo{companies: sampleCompanies},
			wantStatus: http.StatusOK,
			wantName:   "Yandex",
		},
		{
			name:       "not found returns 404",
			url:        "/api/v1/companies/999",
			repo:       &mockCompanyRepo{companies: sampleCompanies},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id returns 400",
			url:        "/api/v1/companies/abc",
			repo:       &mockCompanyRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "db error returns 500",
			url:        "/api/v1/companies/1",
			repo:       &mockCompanyRepo{err: errors.New("timeout")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			makeCompanyRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK && tt.wantName != "" {
				var company models.Company
				if err := json.NewDecoder(rec.Body).Decode(&company); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if company.Name != tt.wantName {
					t.Errorf("name = %q, want %q", company.Name, tt.wantName)
				}
			}
		})
	}
}

func TestCompanyHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		repo       *mockCompanyRepo
		wantStatus int
		wantID     int
	}{
		{
			name:       "create successfully",
			body:       models.Company{Name: "VK", Description: "Соцсеть"},
			repo:       &mockCompanyRepo{},
			wantStatus: http.StatusCreated,
			wantID:     99,
		},
		{
			name:       "invalid json returns 400",
			body:       "not-json",
			repo:       &mockCompanyRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "db error returns 500",
			body:       models.Company{Name: "VK"},
			repo:       &mockCompanyRepo{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/companies", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			makeCompanyRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				var company models.Company
				if err := json.NewDecoder(rec.Body).Decode(&company); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if company.ID != tt.wantID {
					t.Errorf("ID = %d, want %d", company.ID, tt.wantID)
				}
			}
		})
	}
}

func TestCompanyHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		repo       *mockCompanyRepo
		wantStatus int
	}{
		{
			name:       "delete successfully",
			url:        "/api/v1/companies/1",
			repo:       &mockCompanyRepo{},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id returns 400",
			url:        "/api/v1/companies/abc",
			repo:       &mockCompanyRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "db error returns 500",
			url:        "/api/v1/companies/1",
			repo:       &mockCompanyRepo{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			rec := httptest.NewRecorder()

			makeCompanyRouter(tt.repo).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
