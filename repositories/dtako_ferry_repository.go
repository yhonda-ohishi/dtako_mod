package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// DtakoFerryRepository handles database operations for dtako_ferry
type DtakoFerryRepository struct {
	prodDB  *sql.DB
	localDB *sql.DB
}

// NewDtakoFerryRepository creates a new repository instance
func NewDtakoFerryRepository() *DtakoFerryRepository {
	prodDB, _ := GetProductionDB()
	localDB, _ := GetLocalDB()
	
	return &DtakoFerryRepository{
		prodDB:  prodDB,
		localDB: localDB,
	}
}

// GetByDateRange retrieves ferry records within a date range from local database
func (r *DtakoFerryRepository) GetByDateRange(from, to time.Time, route string) ([]models.DtakoFerry, error) {
	query := `
		SELECT id, date, route, vehicle_no, driver_code, 
		       departure_time, arrival_time, passengers, vehicles,
		       created_at, updated_at
		FROM dtako_ferry
		WHERE date BETWEEN ? AND ?
	`
	args := []interface{}{from, to}

	if route != "" {
		query += " AND route = ?"
		args = append(args, route)
	}

	query += " ORDER BY date DESC, departure_time DESC"

	rows, err := r.localDB.Query(query, args...)
	if err != nil {
		return []models.DtakoFerry{}, err
	}
	defer rows.Close()

	results := []models.DtakoFerry{}
	for rows.Next() {
		var record models.DtakoFerry
		err := rows.Scan(
			&record.ID, &record.Date, &record.Route, &record.VehicleNo,
			&record.DriverCode, &record.DepartureTime, &record.ArrivalTime,
			&record.Passengers, &record.Vehicles,
			&record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return []models.DtakoFerry{}, err
		}
		results = append(results, record)
	}

	return results, nil
}

// GetByID retrieves a specific ferry record by ID from local database
func (r *DtakoFerryRepository) GetByID(id string) (*models.DtakoFerry, error) {
	query := `
		SELECT id, date, route, vehicle_no, driver_code, 
		       departure_time, arrival_time, passengers, vehicles,
		       created_at, updated_at
		FROM dtako_ferry
		WHERE id = ?
	`

	var record models.DtakoFerry
	err := r.localDB.QueryRow(query, id).Scan(
		&record.ID, &record.Date, &record.Route, &record.VehicleNo,
		&record.DriverCode, &record.DepartureTime, &record.ArrivalTime,
		&record.Passengers, &record.Vehicles,
		&record.CreatedAt, &record.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

// FetchFromProduction fetches ferry data from production database
func (r *DtakoFerryRepository) FetchFromProduction(from, to time.Time, route string) ([]models.DtakoFerry, error) {
	if r.prodDB == nil {
		return []models.DtakoFerry{}, fmt.Errorf("production database not connected")
	}

	query := `
		SELECT id, date, route, vehicle_no, driver_code, 
		       departure_time, arrival_time, passengers, vehicles,
		       created_at, updated_at
		FROM dtako_ferry
		WHERE date BETWEEN ? AND ?
	`
	args := []interface{}{from, to}

	if route != "" {
		query += " AND route = ?"
		args = append(args, route)
	}

	query += " ORDER BY date DESC, departure_time DESC"

	rows, err := r.prodDB.Query(query, args...)
	if err != nil {
		return []models.DtakoFerry{}, err
	}
	defer rows.Close()

	results := []models.DtakoFerry{}
	for rows.Next() {
		var record models.DtakoFerry
		err := rows.Scan(
			&record.ID, &record.Date, &record.Route, &record.VehicleNo,
			&record.DriverCode, &record.DepartureTime, &record.ArrivalTime,
			&record.Passengers, &record.Vehicles,
			&record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return []models.DtakoFerry{}, err
		}
		results = append(results, record)
	}

	return results, nil
}

// Insert inserts a ferry record into local database
func (r *DtakoFerryRepository) Insert(record *models.DtakoFerry) error {
	query := `
		INSERT INTO dtako_ferry (id, date, route, vehicle_no, driver_code, 
		                        departure_time, arrival_time, passengers, vehicles,
		                        created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		    date = VALUES(date),
		    route = VALUES(route),
		    vehicle_no = VALUES(vehicle_no),
		    driver_code = VALUES(driver_code),
		    departure_time = VALUES(departure_time),
		    arrival_time = VALUES(arrival_time),
		    passengers = VALUES(passengers),
		    vehicles = VALUES(vehicles),
		    updated_at = VALUES(updated_at)
	`

	_, err := r.localDB.Exec(query,
		record.ID, record.Date, record.Route, record.VehicleNo,
		record.DriverCode, record.DepartureTime, record.ArrivalTime,
		record.Passengers, record.Vehicles,
		record.CreatedAt, record.UpdatedAt,
	)

	return err
}