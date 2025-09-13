package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/models"
	"github.com/yhonda-ohishi/dtako_mod/repositories"
)

// DtakoEventsService handles business logic for dtako_events
type DtakoEventsService struct {
	repo *repositories.DtakoEventsRepository
}

// NewDtakoEventsService creates a new service instance
func NewDtakoEventsService() *DtakoEventsService {
	return &DtakoEventsService{
		repo: repositories.NewDtakoEventsRepository(),
	}
}

// GetEvents retrieves events within date range and optional type filter
func (s *DtakoEventsService) GetEvents(from, to, eventType, unkoNo string) ([]models.DtakoEvent, error) {
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

	return s.repo.GetByDateRange(fromDate, toDate, eventType, unkoNo)
}

// GetEventByID retrieves a specific event by ID
func (s *DtakoEventsService) GetEventByID(id string) (*models.DtakoEvent, error) {
	event, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found: %s", id)
		}
		return nil, err
	}
	return event, nil
}

// ImportFromProduction imports event data from production database
func (s *DtakoEventsService) ImportFromProduction(fromDate, toDate, eventType string) (*models.ImportResult, error) {
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

	// Validate event type if specified
	validEventTypes := []string{"START", "STOP", "END", "運転", "休憩", "作業"}
	if eventType != "" {
		isValid := false
		for _, validType := range validEventTypes {
			if eventType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, fmt.Errorf("invalid event_type: %s", eventType)
		}
	}

	// Fetch from production
	events, err := s.repo.FetchFromProduction(from, to, eventType)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from production: %v", err)
	}

	// Import to local database
	imported := 0
	var errors []string

	for _, event := range events {
		if err := s.repo.Insert(&event); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to import event %s: %v", event.ID, err))
		} else {
			imported++
		}
	}

	result := &models.ImportResult{
		Success:      imported > 0,
		ImportedRows: imported,
		Message:      fmt.Sprintf("Imported %d events from %s to %s", imported, fromDate, toDate),
		ImportedAt:   time.Now(),
		Errors:       errors,
	}

	if eventType != "" {
		result.Message += fmt.Sprintf(" (type: %s)", eventType)
	}

	return result, nil
}