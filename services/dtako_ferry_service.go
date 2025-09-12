package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/models"
	"github.com/yhonda-ohishi/dtako_mod/repositories"
)

// DtakoFerryService handles business logic for dtako_ferry
type DtakoFerryService struct {
	repo *repositories.DtakoFerryRepository
}

// NewDtakoFerryService creates a new service instance
func NewDtakoFerryService() *DtakoFerryService {
	return &DtakoFerryService{
		repo: repositories.NewDtakoFerryRepository(),
	}
}

// GetFerryRecords retrieves ferry records within date range and optional route filter
func (s *DtakoFerryService) GetFerryRecords(from, to, route string) ([]models.DtakoFerry, error) {
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

	return s.repo.GetByDateRange(fromDate, toDate, route)
}

// GetFerryRecordByID retrieves a specific ferry record by ID
func (s *DtakoFerryService) GetFerryRecordByID(id string) (*models.DtakoFerry, error) {
	record, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ferry record not found: %s", id)
		}
		return nil, err
	}
	return record, nil
}

// ImportFromProduction imports ferry data from production database
func (s *DtakoFerryService) ImportFromProduction(fromDate, toDate, route string) (*models.ImportResult, error) {
	// Parse dates
	from, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return nil, fmt.Errorf("invalid from date: %v", err)
	}

	to, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return nil, fmt.Errorf("invalid to date: %v", err)
	}

	// Fetch from production
	records, err := s.repo.FetchFromProduction(from, to, route)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from production: %v", err)
	}

	// Import to local database
	imported := 0
	var errors []string

	for _, record := range records {
		if err := s.repo.Insert(&record); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to import ferry record %s: %v", record.ID, err))
		} else {
			imported++
		}
	}

	result := &models.ImportResult{
		Success:      imported > 0,
		ImportedRows: imported,
		Message:      fmt.Sprintf("Imported %d ferry records from %s to %s", imported, fromDate, toDate),
		ImportedAt:   time.Now(),
		Errors:       errors,
	}

	if route != "" {
		result.Message += fmt.Sprintf(" (route: %s)", route)
	}

	return result, nil
}