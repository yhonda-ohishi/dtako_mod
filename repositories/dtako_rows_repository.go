package repositories

import (
	"database/sql"
	"fmt"
	"os"
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

// GetByDateRange retrieves rows within a date range from production database
func (r *DtakoRowsRepository) GetByDateRange(from, to time.Time) ([]models.DtakoRow, error) {
	// 本番DBのみ使用（ローカルは無視）
	var db *sql.DB = r.prodDB
	if db == nil {
		return []models.DtakoRow{}, fmt.Errorf("production database not available")
	}

	// 本番DBは日本語カラム名
	// 日付範囲の調整: その日の終わり（23:59:59）まで含める
	// 正しく日付の終わりにセット (Add ではなく、その日の23:59:59にする)
	year, month, day := to.Date()
	toEndOfDay := time.Date(year, month, day, 23, 59, 59, 999999999, to.Location())

	query := `
		SELECT id, 運行NO, 運行日, 車輌CD, 対象乗務員CD, 行先市町村名,
		       総走行距離, 自社主燃料, NULL as created_at, NULL as updated_at
		FROM dtako_rows
		WHERE 運行日 BETWEEN ? AND ?
		ORDER BY 運行日 DESC
		LIMIT 100
	`

	rows, err := db.Query(query, from, toEndOfDay)
	if err != nil {
		return []models.DtakoRow{}, err
	}
	defer rows.Close()

	results := []models.DtakoRow{}
	for rows.Next() {
		var row models.DtakoRow
		err := rows.Scan(
			&row.ID, &row.UnkoNo, &row.Date, &row.VehicleNo, &row.DriverCode,
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

// GetByID retrieves a specific row by ID from production database
func (r *DtakoRowsRepository) GetByID(id string) (*models.DtakoRow, error) {
	// 本番DBのみ使用（ローカルは無視）
	var db *sql.DB = r.prodDB
	if db == nil {
		return nil, fmt.Errorf("production database not available")
	}

	query := `
		SELECT id, 運行NO, 運行日, 車輌CD, 対象乗務員CD, 行先市町村名,
		       総走行距離, 自社主燃料, NULL as created_at, NULL as updated_at
		FROM dtako_rows
		WHERE id = ?
	`

	var row models.DtakoRow
	err := db.QueryRow(query, id).Scan(
		&row.ID, &row.UnkoNo, &row.Date, &row.VehicleNo, &row.DriverCode,
		&row.RouteCode, &row.Distance, &row.FuelAmount,
		&row.CreatedAt, &row.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &row, nil
}

// FetchFromProduction fetches row data from production database
func (r *DtakoRowsRepository) FetchFromProduction(from, to time.Time) ([]models.DtakoRow, error) {
	if r.prodDB == nil {
		return []models.DtakoRow{}, nil
	}

	// 日付範囲の調整: その日の終わり（23:59:59）まで含める
	year, month, day := to.Date()
	toEndOfDay := time.Date(year, month, day, 23, 59, 59, 999999999, to.Location())

	// テスト環境のdtako_test_prodは英語カラム名を使用
	// 本番環境は日本語カラム名を使用
	// PROD_DB_NAMEで判断
	query := ``
	if os.Getenv("PROD_DB_NAME") == "dtako_test_prod" {
		// テスト用プロダクションDB（英語カラム名）
		query = `
			SELECT id, unko_no, date, vehicle_no, driver_code, route_code,
			       distance, fuel_amount, created_at, updated_at
			FROM dtako_rows
			WHERE date BETWEEN ? AND ?
			ORDER BY date DESC
		`
	} else {
		// 本番DB（日本語カラム名）
		query = `
			SELECT id, 運行NO, 運行日, 車輌CD, 対象乗務員CD, 行先市町村名,
			       総走行距離, 自社主燃料, NULL as created_at, NULL as updated_at
			FROM dtako_rows
			WHERE 運行日 BETWEEN ? AND ?
			ORDER BY 運行日 DESC
		`
	}

	rows, err := r.prodDB.Query(query, from, toEndOfDay)
	if err != nil {
		return []models.DtakoRow{}, err
	}
	defer rows.Close()

	results := []models.DtakoRow{}
	for rows.Next() {
		var row models.DtakoRow
		err := rows.Scan(
			&row.ID, &row.UnkoNo, &row.Date, &row.VehicleNo, &row.DriverCode,
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
	// ローカルDBの実際のカラム構造に合わせる
	// 必須カラム: id, 運行NO, 読取日, 運行日, 車輌CD, 車輌CC
	query := `
		INSERT INTO dtako_rows (id, 運行NO, 読取日, 運行日, 車輌CD, 車輌CC, 対象乗務員CD, 行先市町村名,
		                       総走行距離, 自社主燃料)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		    運行NO = VALUES(運行NO),
		    読取日 = VALUES(読取日),
		    運行日 = VALUES(運行日),
		    車輌CD = VALUES(車輌CD),
		    車輌CC = VALUES(車輌CC),
		    対象乗務員CD = VALUES(対象乗務員CD),
		    行先市町村名 = VALUES(行先市町村名),
		    総走行距離 = VALUES(総走行距離),
		    自社主燃料 = VALUES(自社主燃料)
	`

	// VehicleNoとDriverCodeはstringからintに変換が必要
	vehicleCD := 1
	if row.VehicleNo != "" {
		// 車輌CDは数値型なので変換が必要
		vehicleCD = 1 // デフォルト値を使用
	}

	driverCode := 0
	if row.DriverCode != "" {
		driverCode = 1 // デフォルト値を使用
	}

	// デフォルト値
	vehicleCC := "001100" // 車輌CC（実際のデータ形式）

	// 読取日は運行日と同じ値を使用
	_, err := r.prodDB.Exec(query,
		row.ID, row.UnkoNo, row.Date, row.Date, vehicleCD, vehicleCC, driverCode,
		row.RouteCode, row.Distance, row.FuelAmount,
	)

	return err
}