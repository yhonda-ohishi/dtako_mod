package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// DtakoRowsRepository handles database operations for dtako_rows
type DtakoRowsRepository struct {
	prodDB  *sql.DB
	localDB *sql.DB
}

// NewDtakoRowsRepository creates a new repository instance
func NewDtakoRowsRepository() *DtakoRowsRepository {
	prodDB, _ := GetProductionDB()
	localDB, _ := GetLocalDB()
	
	return &DtakoRowsRepository{
		prodDB:  prodDB,
		localDB: localDB,
	}
}

// GetByDateRange retrieves rows within a date range from local database
func (r *DtakoRowsRepository) GetByDateRange(from, to time.Time) ([]models.DtakoRow, error) {
	query := `
		SELECT id, date, vehicle_no, driver_code, route_code, 
		       distance, fuel_amount, created_at, updated_at
		FROM dtako_rows
		WHERE date BETWEEN ? AND ?
		ORDER BY date DESC
	`

	rows, err := r.localDB.Query(query, from, to)
	if err != nil {
		return []models.DtakoRow{}, err
	}
	defer rows.Close()

	results := []models.DtakoRow{}
	for rows.Next() {
		var row models.DtakoRow
		err := rows.Scan(
			&row.ID, &row.Date, &row.VehicleNo, &row.DriverCode,
			&row.RouteCode, &row.Distance, &row.FuelAmount,
			&row.CreatedAt, &row.UpdatedAt,
		)
		if err != nil {
			return []models.DtakoRow{}, err
		}
		results = append(results, row)
	}

	return results, nil
}

// GetByID retrieves a specific row by ID from local database
func (r *DtakoRowsRepository) GetByID(id string) (*models.DtakoRow, error) {
	query := `
		SELECT id, date, vehicle_no, driver_code, route_code, 
		       distance, fuel_amount, created_at, updated_at
		FROM dtako_rows
		WHERE id = ?
	`

	var row models.DtakoRow
	err := r.localDB.QueryRow(query, id).Scan(
		&row.ID, &row.Date, &row.VehicleNo, &row.DriverCode,
		&row.RouteCode, &row.Distance, &row.FuelAmount,
		&row.CreatedAt, &row.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &row, nil
}

// FetchFromProduction fetches data from production database
func (r *DtakoRowsRepository) FetchFromProduction(from, to time.Time) ([]models.DtakoRow, error) {
	if r.prodDB == nil {
		return []models.DtakoRow{}, fmt.Errorf("production database not connected")
	}

	query := `
		SELECT id, date, vehicle_no, driver_code, route_code, 
		       distance, fuel_amount, created_at, updated_at
		FROM dtako_rows
		WHERE date BETWEEN ? AND ?
		ORDER BY date DESC
	`

	rows, err := r.prodDB.Query(query, from, to)
	if err != nil {
		return []models.DtakoRow{}, err
	}
	defer rows.Close()

	results := []models.DtakoRow{}
	for rows.Next() {
		var row models.DtakoRow
		err := rows.Scan(
			&row.ID, &row.Date, &row.VehicleNo, &row.DriverCode,
			&row.RouteCode, &row.Distance, &row.FuelAmount,
			&row.CreatedAt, &row.UpdatedAt,
		)
		if err != nil {
			return []models.DtakoRow{}, err
		}
		results = append(results, row)
	}

	return results, nil
}

// Insert inserts a row into local database
func (r *DtakoRowsRepository) Insert(row *models.DtakoRow) error {
	query := `
		INSERT INTO dtako_rows (id, date, vehicle_no, driver_code, route_code, 
		                        distance, fuel_amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		    date = VALUES(date),
		    vehicle_no = VALUES(vehicle_no),
		    driver_code = VALUES(driver_code),
		    route_code = VALUES(route_code),
		    distance = VALUES(distance),
		    fuel_amount = VALUES(fuel_amount),
		    updated_at = VALUES(updated_at)
	`

	_, err := r.localDB.Exec(query,
		row.ID, row.Date, row.VehicleNo, row.DriverCode,
		row.RouteCode, row.Distance, row.FuelAmount,
		row.CreatedAt, row.UpdatedAt,
	)

	return err
}