package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/models"
	"github.com/yhonda-ohishi/dtako_mod/repositories"
)

// DtakoRowsService handles business logic for dtako_rows
type DtakoRowsService struct {
	repo *repositories.DtakoRowsRepository
}

// NewDtakoRowsService creates a new service instance
func NewDtakoRowsService() *DtakoRowsService {
	return &DtakoRowsService{
		repo: repositories.NewDtakoRowsRepository(),
	}
}

// GetRows retrieves rows within date range
func (s *DtakoRowsService) GetRows(from, to string) ([]models.DtakoRow, error) {
	// Parse dates if provided
	var fromDate, toDate time.Time
	var err error

	if from != "" {
		fromDate, err = time.Parse("2006-01-02", from)
		if err != nil {
			return nil, fmt.Errorf("invalid from date: %v", err)
		}
	} else {
		fromDate = time.Now().AddDate(0, -1, 0)
	}

	if to != "" {
		toDate, err = time.Parse("2006-01-02", to)
		if err != nil {
			return nil, fmt.Errorf("invalid to date: %v", err)
		}
	} else {
		toDate = time.Now()
	}

	return s.repo.GetByDateRange(fromDate, toDate)
}

// GetRowByID retrieves a specific row by ID
func (s *DtakoRowsService) GetRowByID(id string) (*models.DtakoRow, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("row not found: %s", id)
		}
		return nil, err
	}
	return row, nil
}

// ImportFromProduction imports data from production database
func (s *DtakoRowsService) ImportFromProduction(fromDate, toDate string) (*models.ImportResult, error) {
	// Parse dates
	from, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return nil, fmt.Errorf("invalid from date: %v", err)
	}

	to, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return nil, fmt.Errorf("invalid to date: %v", err)
	}

	// Validate date range
	if from.After(to) {
		return nil, fmt.Errorf("from_date cannot be after to_date")
	}

	// Fetch from production
	rows, err := s.repo.FetchFromProduction(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from production: %v", err)
	}

	// Import to local database
	imported := 0
	var errors []string

	for _, row := range rows {
		if err := s.repo.Insert(&row); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to import row %s: %v", row.ID, err))
		} else {
			imported++
		}
	}

	result := &models.ImportResult{
		Success:      imported > 0,
		ImportedRows: imported,
		Message:      fmt.Sprintf("Imported %d rows from %s to %s", imported, fromDate, toDate),
		ImportedAt:   time.Now(),
		Errors:       errors,
	}

	return result, nil
}