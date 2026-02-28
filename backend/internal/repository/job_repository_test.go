package repository

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"job-statistics-api/internal/models"
)

// jobColumns — список колонок для scan в тестах
var jobColumns = []string{
	"id", "company_id", "title", "level", "specialization",
	"salary_min", "salary_max", "salary_currency", "experience_years",
	"location", "remote_available", "description", "responsibilities",
	"benefits", "posted_date", "is_active", "source_url",
	"created_at", "updated_at",
}

func newJobRow(id, companyID int, title, level string) []driver.Value {
	now := time.Now()
	return []driver.Value{
		int64(id), int64(companyID), title, level,
		nil,   // specialization
		nil, nil, // salary_min, salary_max
		"RUB",
		nil,      // experience_years
		"Москва", // location
		false,
		nil, nil, nil, // description, responsibilities, benefits
		nil,  // posted_date
		true, // is_active
		nil,  // source_url
		now, now,
	}
}

func TestJobRepository_GetAll(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantLen   int
		wantErr   bool
	}{
		{
			name: "returns two jobs",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(jobColumns).
					AddRow(newJobRow(1, 1, "Go Developer", "Middle")...).
					AddRow(newJobRow(2, 2, "Python Developer", "Senior")...)
				mock.ExpectQuery("SELECT id, company_id, title").WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns empty list",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(jobColumns)
				mock.ExpectQuery("SELECT id, company_id, title").WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, company_id, title").
					WillReturnError(errors.New("timeout"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewJobRepository(db)

			jobs, err := repo.GetAll()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(jobs) != tt.wantLen {
				t.Errorf("GetAll() returned %d jobs, want %d", len(jobs), tt.wantLen)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestJobRepository_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(mock sqlmock.Sqlmock)
		wantTitle string
		wantErr   bool
	}{
		{
			name: "found job",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(jobColumns).
					AddRow(newJobRow(1, 1, "Go Developer", "Middle")...)
				mock.ExpectQuery("SELECT id, company_id, title").
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantTitle: "Go Developer",
			wantErr:   false,
		},
		{
			name: "not found",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, company_id, title").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewJobRepository(db)

			job, err := repo.GetByID(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && job.Title != tt.wantTitle {
				t.Errorf("GetByID() title = %q, want %q", job.Title, tt.wantTitle)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestJobRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		job       models.Job
		mockSetup func(mock sqlmock.Sqlmock)
		wantID    int
		wantErr   bool
	}{
		{
			name: "create successfully",
			job: models.Job{
				CompanyID: 1,
				Title:     "Go Developer",
				Level:     "Middle",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO jobs").
					WillReturnResult(sqlmock.NewResult(10, 1))
			},
			wantID:  10,
			wantErr: false,
		},
		{
			name: "db error on insert",
			job:  models.Job{CompanyID: 1, Title: "Bad Job", Level: "Junior"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO jobs").
					WillReturnError(errors.New("constraint violation"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewJobRepository(db)

			j := tt.job
			err = repo.Create(&j)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && j.ID != tt.wantID {
				t.Errorf("Create() ID = %d, want %d", j.ID, tt.wantID)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestJobRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		job       models.Job
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "update successfully",
			job:  models.Job{ID: 1, CompanyID: 1, Title: "Senior Go Dev", Level: "Senior"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE jobs SET").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			job:  models.Job{ID: 99, Title: "Ghost"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE jobs SET").
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewJobRepository(db)

			j := tt.job
			err = repo.Update(&j)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestJobRepository_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "delete successfully",
			id:   5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM jobs WHERE id = \\?").
					WithArgs(5).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM jobs WHERE id = \\?").
					WithArgs(999).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)
			repo := NewJobRepository(db)

			err = repo.Delete(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}
