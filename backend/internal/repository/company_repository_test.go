package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"job-statistics-api/internal/models"
)

func TestCompanyRepository_GetAll(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		wantLen   int
		wantErr   bool
	}{
		{
			name: "returns two companies",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
					AddRow(1, "Yandex", "Поисковик", now, now).
					AddRow(2, "Sber", "Банк", now, now)
				mock.ExpectQuery("SELECT id, name, description, created_at, updated_at FROM companies ORDER BY name").
					WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns empty list",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, name, description, created_at, updated_at FROM companies ORDER BY name").
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, description, created_at, updated_at FROM companies ORDER BY name").
					WillReturnError(errors.New("connection lost"))
			},
			wantLen: 0,
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
			repo := NewCompanyRepository(db)

			companies, err := repo.GetAll()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(companies) != tt.wantLen {
				t.Errorf("GetAll() returned %d companies, want %d", len(companies), tt.wantLen)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestCompanyRepository_GetByID(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		id        int
		mockSetup func(mock sqlmock.Sqlmock)
		wantName  string
		wantErr   bool
	}{
		{
			name: "found company",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
					AddRow(1, "Yandex", "Поисковик", now, now)
				mock.ExpectQuery("SELECT id, name, description, created_at, updated_at FROM companies WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantName: "Yandex",
			wantErr:  false,
		},
		{
			name: "not found",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, description, created_at, updated_at FROM companies WHERE id = ?").
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
			repo := NewCompanyRepository(db)

			company, err := repo.GetByID(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && company.Name != tt.wantName {
				t.Errorf("GetByID() name = %q, want %q", company.Name, tt.wantName)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestCompanyRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		company   models.Company
		mockSetup func(mock sqlmock.Sqlmock)
		wantID    int
		wantErr   bool
	}{
		{
			name:    "create successfully",
			company: models.Company{Name: "VK", Description: "Соцсеть"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO companies").
					WithArgs("VK", "Соцсеть").
					WillReturnResult(sqlmock.NewResult(42, 1))
			},
			wantID:  42,
			wantErr: false,
		},
		{
			name:    "duplicate name error",
			company: models.Company{Name: "Yandex", Description: ""},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO companies").
					WithArgs("Yandex", "").
					WillReturnError(errors.New("Duplicate entry 'Yandex' for key 'name'"))
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
			repo := NewCompanyRepository(db)

			c := tt.company
			err = repo.Create(&c)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && c.ID != tt.wantID {
				t.Errorf("Create() ID = %d, want %d", c.ID, tt.wantID)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestCompanyRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		company   models.Company
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:    "update successfully",
			company: models.Company{ID: 1, Name: "Яндекс", Description: "Технологическая компания"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE companies SET name = \\?, description = \\? WHERE id = \\?").
					WithArgs("Яндекс", "Технологическая компания", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:    "db error",
			company: models.Company{ID: 99, Name: "Ghost"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE companies SET name = \\?, description = \\? WHERE id = \\?").
					WithArgs("Ghost", "", 99).
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
			repo := NewCompanyRepository(db)

			c := tt.company
			err = repo.Update(&c)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestCompanyRepository_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "delete successfully",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM companies WHERE id = \\?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM companies WHERE id = \\?").
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
			repo := NewCompanyRepository(db)

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
