package dto

import (
	"database/sql"
	"job-statistics-api/internal/models"
)

// LocationResponse — DTO для отправки локации клиенту.
type LocationResponse struct {
	ID           int     `json:"id"`
	JobID        int     `json:"job_id"`
	City         string  `json:"city"`
	MetroStation *string `json:"metro_station"`
	IsPrimary    bool    `json:"is_primary"`
}

// LocationRequest — DTO для приёма данных от клиента.
type LocationRequest struct {
	JobID        int     `json:"job_id"`
	City         string  `json:"city"`
	MetroStation *string `json:"metro_station"`
	IsPrimary    bool    `json:"is_primary"`
}

// LocationResponseFromModel конвертирует models.Location → LocationResponse
func LocationResponseFromModel(l models.Location) LocationResponse {
	r := LocationResponse{
		ID:        l.ID,
		JobID:     l.JobID,
		City:      l.City,
		IsPrimary: l.IsPrimary,
	}

	if l.MetroStation.Valid {
		r.MetroStation = &l.MetroStation.String
	}

	return r
}

// LocationResponseList конвертирует []models.Location → []LocationResponse
func LocationResponseList(locations []models.Location) []LocationResponse {
	result := make([]LocationResponse, len(locations))
	for i, l := range locations {
		result[i] = LocationResponseFromModel(l)
	}
	return result
}

// ToModel конвертирует LocationRequest → models.Location
func (req LocationRequest) ToModel() models.Location {
	l := models.Location{
		JobID:     req.JobID,
		City:      req.City,
		IsPrimary: req.IsPrimary,
	}

	if req.MetroStation != nil && *req.MetroStation != "" {
		l.MetroStation = sql.NullString{String: *req.MetroStation, Valid: true}
	}

	return l
}
