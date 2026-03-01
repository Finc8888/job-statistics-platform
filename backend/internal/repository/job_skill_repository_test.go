package repository

import (
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var skillColumns = []string{"id", "name", "category", "created_at"}

func newSkillRow(id int, name, category string) []driver.Value {
	return []driver.Value{int64(id), name, category, time.Now()}
}

func TestJobSkillRepository_GetSkillsByJobID(t *testing.T) {
	tests := []struct {
		name      string
		jobID     int
		mockSetup func(mock sqlmock.Sqlmock)
		wantLen   int
		wantErr   bool
	}{
		{
			name:  "returns skills for a job",
			jobID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(skillColumns).
					AddRow(newSkillRow(10, "Go", "Язык программирования")...).
					AddRow(newSkillRow(11, "PostgreSQL", "База данных")...)
				mock.ExpectQuery("SELECT s.id, s.name, s.category, s.created_at").
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:  "returns empty slice when no skills associated",
			jobID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(skillColumns)
				mock.ExpectQuery("SELECT s.id, s.name, s.category, s.created_at").
					WithArgs(2).
					WillReturnRows(rows)
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:  "database error",
			jobID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT s.id, s.name, s.category, s.created_at").
					WithArgs(3).
					WillReturnError(errors.New("connection reset"))
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
			repo := NewJobSkillRepository(db)

			skills, err := repo.GetSkillsByJobID(tt.jobID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSkillsByJobID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(skills) != tt.wantLen {
				t.Errorf("GetSkillsByJobID() returned %d skills, want %d", len(skills), tt.wantLen)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}

func TestJobSkillRepository_SetJobSkills(t *testing.T) {
	tests := []struct {
		name      string
		jobID     int
		skillIDs  []int
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name:     "set two skills successfully",
			jobID:    5,
			skillIDs: []int{10, 20},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM job_skills WHERE job_id = \\?").
					WithArgs(5).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectExec("INSERT INTO job_skills").
					WithArgs(5, 10, 5, 20).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:     "clear all skills (empty slice)",
			jobID:    6,
			skillIDs: []int{},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM job_skills WHERE job_id = \\?").
					WithArgs(6).
					WillReturnResult(sqlmock.NewResult(0, 3))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:     "begin transaction error",
			jobID:    7,
			skillIDs: []int{1},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("connection lost"))
			},
			wantErr: true,
		},
		{
			name:     "delete error triggers rollback",
			jobID:    8,
			skillIDs: []int{1},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM job_skills WHERE job_id = \\?").
					WithArgs(8).
					WillReturnError(errors.New("lock timeout"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:     "insert error triggers rollback",
			jobID:    9,
			skillIDs: []int{1, 2},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM job_skills WHERE job_id = \\?").
					WithArgs(9).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("INSERT INTO job_skills").
					WithArgs(9, 1, 9, 2).
					WillReturnError(errors.New("foreign key constraint"))
				mock.ExpectRollback()
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
			repo := NewJobSkillRepository(db)

			err = repo.SetJobSkills(tt.jobID, tt.skillIDs)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetJobSkills() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled mock expectations: %v", err)
			}
		})
	}
}
