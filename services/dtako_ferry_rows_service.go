package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/models"
	"github.com/yhonda-ohishi/dtako_mod/repositories"
)

// DtakoFerryRowsService handles business logic for dtako_ferry_rows
type DtakoFerryRowsService struct {
	repo *repositories.DtakoFerryRowsRepository
}

// NewDtakoFerryRowsService creates a new service instance
func NewDtakoFerryRowsService() *DtakoFerryRowsService {
	return &DtakoFerryRowsService{
		repo: repositories.NewDtakoFerryRowsRepository(),
	}
}

// GetFerryRows retrieves ferry row records within date range and optional ferry company filter
func (s *DtakoFerryRowsService) GetFerryRows(from, to, ferryCompany string) ([]models.DtakoFerryRow, error) {
	// Parse dates if provided
	var fromDate, toDate time.Time
	var err error

	// JSTタイムゾーンを取得
	jst, _ := time.LoadLocation("Asia/Tokyo")

	if from != "" {
		fromDate, err = time.ParseInLocation("2006-01-02", from, jst)
		if err != nil {
			return nil, fmt.Errorf("invalid from date: %v", err)
		}
	} else {
		fromDate = time.Now().In(jst).AddDate(0, -1, 0)
	}

	if to != "" {
		toDate, err = time.ParseInLocation("2006-01-02", to, jst)
		if err != nil {
			return nil, fmt.Errorf("invalid to date: %v", err)
		}
	} else {
		toDate = time.Now().In(jst)
	}

	return s.repo.GetByDateRange(fromDate, toDate, ferryCompany)
}

// GetFerryRowByID retrieves a specific ferry row record by ID
func (s *DtakoFerryRowsService) GetFerryRowByID(id string) (*models.DtakoFerryRow, error) {
	record, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ferry row record not found: %s", id)
		}
		return nil, err
	}
	return record, nil
}

// ImportFromProduction imports ferry row data from production database
func (s *DtakoFerryRowsService) ImportFromProduction(fromDate, toDate, ferryCompany string) (*models.ImportResult, error) {
	// JSTタイムゾーンを取得
	jst, _ := time.LoadLocation("Asia/Tokyo")

	// Parse dates
	from, err := time.ParseInLocation("2006-01-02", fromDate, jst)
	if err != nil {
		return nil, fmt.Errorf("invalid from date: %v", err)
	}

	to, err := time.ParseInLocation("2006-01-02", toDate, jst)
	if err != nil {
		return nil, fmt.Errorf("invalid to date: %v", err)
	}

	// Validate date range
	if from.After(to) {
		return nil, fmt.Errorf("from_date cannot be after to_date")
	}

	// Fetch from production
	records, err := s.repo.FetchFromProduction(from, to, ferryCompany)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from production: %v", err)
	}

	// Import to local database
	imported := 0
	var errors []string

	for _, record := range records {
		if err := s.repo.Insert(&record); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to import ferry row record %d: %v", record.ID, err))
		} else {
			imported++
		}
	}

	result := &models.ImportResult{
		Success:      imported > 0,
		ImportedRows: imported,
		Message:      fmt.Sprintf("Imported %d ferry row records from %s to %s", imported, fromDate, toDate),
		ImportedAt:   time.Now(),
		Errors:       errors,
	}

	if ferryCompany != "" {
		result.Message += fmt.Sprintf(" (ferry company: %s)", ferryCompany)
	}

	return result, nil
}