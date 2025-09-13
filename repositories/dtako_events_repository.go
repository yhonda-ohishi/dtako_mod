package repositories

import (
	"database/sql"
	"os"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// DtakoEventsRepository handles database operations for dtako_events
type DtakoEventsRepository struct {
	prodDB  *sql.DB
	localDB *sql.DB
}

// NewDtakoEventsRepository creates a new repository instance
func NewDtakoEventsRepository() *DtakoEventsRepository {
	prodDB, _ := GetProductionDB()
	localDB, _ := GetLocalDB()

	return &DtakoEventsRepository{
		prodDB:  prodDB,
		localDB: localDB,
	}
}

// GetByDateRange retrieves events within a date range from local database
func (r *DtakoEventsRepository) GetByDateRange(from, to time.Time, eventType, unkoNo string) ([]models.DtakoEvent, error) {
	query := `
		SELECT id, COALESCE(運行NO, ''), event_date, event_type, vehicle_no, driver_code,
		       description, latitude, longitude, created_at, updated_at
		FROM dtako_events
		WHERE event_date BETWEEN ? AND ?
	`
	args := []interface{}{from, to}

	if eventType != "" {
		query += " AND event_type = ?"
		args = append(args, eventType)
	}

	if unkoNo != "" {
		query += " AND 運行NO = ?"
		args = append(args, unkoNo)
	}

	query += " ORDER BY event_date DESC"

	rows, err := r.localDB.Query(query, args...)
	if err != nil {
		return []models.DtakoEvent{}, err
	}
	defer rows.Close()

	results := []models.DtakoEvent{}
	for rows.Next() {
		var event models.DtakoEvent
		err := rows.Scan(
			&event.ID, &event.UnkoNo, &event.EventDate, &event.EventType, &event.VehicleNo,
			&event.DriverCode, &event.Description, &event.Latitude, &event.Longitude,
			&event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return []models.DtakoEvent{}, err
		}
		results = append(results, event)
	}

	return results, nil
}

// GetByID retrieves a specific event by ID from local database
func (r *DtakoEventsRepository) GetByID(id string) (*models.DtakoEvent, error) {
	query := `
		SELECT id, COALESCE(運行NO, ''), event_date, event_type, vehicle_no, driver_code,
		       description, latitude, longitude, created_at, updated_at
		FROM dtako_events
		WHERE id = ?
	`

	var event models.DtakoEvent
	err := r.localDB.QueryRow(query, id).Scan(
		&event.ID, &event.UnkoNo, &event.EventDate, &event.EventType, &event.VehicleNo,
		&event.DriverCode, &event.Description, &event.Latitude, &event.Longitude,
		&event.CreatedAt, &event.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

// FetchFromProduction fetches event data from production database
func (r *DtakoEventsRepository) FetchFromProduction(from, to time.Time, eventType string) ([]models.DtakoEvent, error) {
	if r.prodDB == nil {
		return []models.DtakoEvent{}, nil
	}

	// テスト環境のdtako_test_prodは英語カラム名を使用
	// 本番環境は日本語カラム名を使用
	query := ``
	if os.Getenv("PROD_DB_NAME") == "dtako_test_prod" {
		// テスト用プロダクションDB（英語カラム名）
		query = `
			SELECT id, COALESCE(unko_no, ''), event_date, event_type, vehicle_no, driver_code,
			       description, latitude, longitude, created_at, updated_at
			FROM dtako_events
			WHERE event_date BETWEEN ? AND ?
		`
	} else {
		// 本番DB（日本語カラム名）
		query = `
			SELECT id, COALESCE(運行NO, ''), event_date, event_type, vehicle_no, driver_code,
			       description, latitude, longitude, created_at, updated_at
			FROM dtako_events
			WHERE event_date BETWEEN ? AND ?
		`
	}
	args := []interface{}{from, to}

	if eventType != "" {
		query += " AND event_type = ?"
		args = append(args, eventType)
	}

	query += " ORDER BY event_date DESC"

	rows, err := r.prodDB.Query(query, args...)
	if err != nil {
		return []models.DtakoEvent{}, err
	}
	defer rows.Close()

	results := []models.DtakoEvent{}
	for rows.Next() {
		var event models.DtakoEvent
		err := rows.Scan(
			&event.ID, &event.UnkoNo, &event.EventDate, &event.EventType, &event.VehicleNo,
			&event.DriverCode, &event.Description, &event.Latitude, &event.Longitude,
			&event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return []models.DtakoEvent{}, err
		}
		results = append(results, event)
	}

	return results, nil
}

// Insert inserts an event into local database
func (r *DtakoEventsRepository) Insert(event *models.DtakoEvent) error {
	query := `
		INSERT INTO dtako_events (id, 運行NO, event_date, event_type, vehicle_no, driver_code,
		                         description, latitude, longitude, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		    運行NO = VALUES(運行NO),
		    event_date = VALUES(event_date),
		    event_type = VALUES(event_type),
		    vehicle_no = VALUES(vehicle_no),
		    driver_code = VALUES(driver_code),
		    description = VALUES(description),
		    latitude = VALUES(latitude),
		    longitude = VALUES(longitude),
		    updated_at = VALUES(updated_at)
	`

	unkoNo := sql.NullString{String: event.UnkoNo, Valid: event.UnkoNo != ""}
	_, err := r.localDB.Exec(query,
		event.ID, unkoNo, event.EventDate, event.EventType, event.VehicleNo,
		event.DriverCode, event.Description, event.Latitude, event.Longitude,
		event.CreatedAt, event.UpdatedAt,
	)

	return err
}