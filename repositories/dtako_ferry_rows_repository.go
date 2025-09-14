package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// DtakoFerryRowsRepository handles database operations for dtako_ferry_rows
type DtakoFerryRowsRepository struct {
	prodDB  *sql.DB
	localDB *sql.DB
}

// NewDtakoFerryRowsRepository creates a new repository instance
func NewDtakoFerryRowsRepository() *DtakoFerryRowsRepository {
	prodDB, _ := GetProductionDB()
	localDB, _ := GetLocalDB()

	return &DtakoFerryRowsRepository{
		prodDB:  prodDB,
		localDB: localDB,
	}
}

// GetByDateRange retrieves ferry row records within a date range from local database
func (r *DtakoFerryRowsRepository) GetByDateRange(from, to time.Time, ferryCompany string) ([]models.DtakoFerryRow, error) {
	query := `
		SELECT id, 運行NO, 運行日, 読取日, 事業所CD, 事業所名,
		       車輌CD, 車輌名, 乗務員CD1, 乗務員名１, 対象乗務員区分,
		       開始日時, 終了日時, フェリー会社CD, フェリー会社名,
		       乗場CD, 乗場名, 便, 降場CD, 降場名,
		       精算区分, 精算区分名, 標準料金, 契約料金,
		       航送車種区分, 航送車種区分名, 見なし距離,
		       COALESCE(ferry_srch, '')
		FROM dtako_ferry_rows
		WHERE 運行日 BETWEEN ? AND ?
	`
	args := []interface{}{from, to}

	if ferryCompany != "" {
		query += " AND フェリー会社名 = ?"
		args = append(args, ferryCompany)
	}

	query += " ORDER BY 運行日 DESC, 開始日時 DESC LIMIT 100"

	// 本番DBのみ使用（ローカルは無視）
	var db *sql.DB = r.prodDB
	if db == nil {
		return []models.DtakoFerryRow{}, fmt.Errorf("production database not available")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return []models.DtakoFerryRow{}, err
	}
	defer rows.Close()

	results := []models.DtakoFerryRow{}
	for rows.Next() {
		var record models.DtakoFerryRow
		err := rows.Scan(
			&record.ID, &record.UnkoNo, &record.UnkoDate, &record.ReadDate,
			&record.OfficeCode, &record.OfficeName, &record.VehicleCode, &record.VehicleName,
			&record.DriverCode1, &record.DriverName1, &record.TargetDriverClass,
			&record.StartTime, &record.EndTime, &record.FerryCompanyCode, &record.FerryCompanyName,
			&record.BoardingCode, &record.BoardingName, &record.ShipNumber,
			&record.LandingCode, &record.LandingName, &record.SettlementClass, &record.SettlementName,
			&record.StandardFare, &record.ContractFare, &record.ShipVehicleClass, &record.ShipVehicleName,
			&record.EstimatedDistance, &record.FerrySearch,
		)
		if err != nil {
			return []models.DtakoFerryRow{}, err
		}
		results = append(results, record)
	}

	return results, nil
}

// GetByID retrieves a specific ferry row record by ID from production database
func (r *DtakoFerryRowsRepository) GetByID(id string) (*models.DtakoFerryRow, error) {
	// 本番DBのみ使用（ローカルは無視）
	var db *sql.DB = r.prodDB
	if db == nil {
		return nil, fmt.Errorf("production database not available")
	}

	query := `
		SELECT id, 運行NO, 運行日, 読取日, 事業所CD, 事業所名,
		       車輌CD, 車輌名, 乗務員CD1, 乗務員名１, 対象乗務員区分,
		       開始日時, 終了日時, フェリー会社CD, フェリー会社名,
		       乗場CD, 乗場名, 便, 降場CD, 降場名,
		       精算区分, 精算区分名, 標準料金, 契約料金,
		       航送車種区分, 航送車種区分名, 見なし距離,
		       COALESCE(ferry_srch, '')
		FROM dtako_ferry_rows
		WHERE id = ?
	`

	var record models.DtakoFerryRow
	err := db.QueryRow(query, id).Scan(
		&record.ID, &record.UnkoNo, &record.UnkoDate, &record.ReadDate,
		&record.OfficeCode, &record.OfficeName, &record.VehicleCode, &record.VehicleName,
		&record.DriverCode1, &record.DriverName1, &record.TargetDriverClass,
		&record.StartTime, &record.EndTime, &record.FerryCompanyCode, &record.FerryCompanyName,
		&record.BoardingCode, &record.BoardingName, &record.ShipNumber,
		&record.LandingCode, &record.LandingName, &record.SettlementClass, &record.SettlementName,
		&record.StandardFare, &record.ContractFare, &record.ShipVehicleClass, &record.ShipVehicleName,
		&record.EstimatedDistance, &record.FerrySearch,
	)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

// FetchFromProduction fetches ferry row data from production database
func (r *DtakoFerryRowsRepository) FetchFromProduction(from, to time.Time, ferryCompany string) ([]models.DtakoFerryRow, error) {
	if r.prodDB == nil {
		return []models.DtakoFerryRow{}, fmt.Errorf("production database not connected")
	}

	query := `
		SELECT id, 運行NO, 運行日, 読取日, 事業所CD, 事業所名,
		       車輌CD, 車輌名, 乗務員CD1, 乗務員名１, 対象乗務員区分,
		       開始日時, 終了日時, フェリー会社CD, フェリー会社名,
		       乗場CD, 乗場名, 便, 降場CD, 降場名,
		       精算区分, 精算区分名, 標準料金, 契約料金,
		       航送車種区分, 航送車種区分名, 見なし距離,
		       COALESCE(ferry_srch, '')
		FROM dtako_ferry_rows
		WHERE 運行日 BETWEEN ? AND ?
	`
	args := []interface{}{from, to}

	if ferryCompany != "" {
		query += " AND フェリー会社名 = ?"
		args = append(args, ferryCompany)
	}

	query += " ORDER BY 運行日 DESC, 開始日時 DESC"

	rows, err := r.prodDB.Query(query, args...)
	if err != nil {
		return []models.DtakoFerryRow{}, err
	}
	defer rows.Close()

	results := []models.DtakoFerryRow{}
	for rows.Next() {
		var record models.DtakoFerryRow
		err := rows.Scan(
			&record.ID, &record.UnkoNo, &record.UnkoDate, &record.ReadDate,
			&record.OfficeCode, &record.OfficeName, &record.VehicleCode, &record.VehicleName,
			&record.DriverCode1, &record.DriverName1, &record.TargetDriverClass,
			&record.StartTime, &record.EndTime, &record.FerryCompanyCode, &record.FerryCompanyName,
			&record.BoardingCode, &record.BoardingName, &record.ShipNumber,
			&record.LandingCode, &record.LandingName, &record.SettlementClass, &record.SettlementName,
			&record.StandardFare, &record.ContractFare, &record.ShipVehicleClass, &record.ShipVehicleName,
			&record.EstimatedDistance, &record.FerrySearch,
		)
		if err != nil {
			return []models.DtakoFerryRow{}, err
		}
		results = append(results, record)
	}

	return results, nil
}

// Insert inserts a ferry row record into local database
func (r *DtakoFerryRowsRepository) Insert(record *models.DtakoFerryRow) error {
	query := `
		INSERT INTO dtako_ferry_rows (
			id, 運行NO, 運行日, 読取日, 事業所CD, 事業所名,
			車輌CD, 車輌名, 乗務員CD1, 乗務員名１, 対象乗務員区分,
			開始日時, 終了日時, フェリー会社CD, フェリー会社名,
			乗場CD, 乗場名, 便, 降場CD, 降場名,
			精算区分, 精算区分名, 標準料金, 契約料金,
			航送車種区分, 航送車種区分名, 見なし距離, ferry_srch
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			運行NO = VALUES(運行NO),
			運行日 = VALUES(運行日),
			読取日 = VALUES(読取日),
			事業所CD = VALUES(事業所CD),
			事業所名 = VALUES(事業所名),
			車輌CD = VALUES(車輌CD),
			車輌名 = VALUES(車輌名),
			乗務員CD1 = VALUES(乗務員CD1),
			乗務員名１ = VALUES(乗務員名１),
			対象乗務員区分 = VALUES(対象乗務員区分),
			開始日時 = VALUES(開始日時),
			終了日時 = VALUES(終了日時),
			フェリー会社CD = VALUES(フェリー会社CD),
			フェリー会社名 = VALUES(フェリー会社名),
			乗場CD = VALUES(乗場CD),
			乗場名 = VALUES(乗場名),
			便 = VALUES(便),
			降場CD = VALUES(降場CD),
			降場名 = VALUES(降場名),
			精算区分 = VALUES(精算区分),
			精算区分名 = VALUES(精算区分名),
			標準料金 = VALUES(標準料金),
			契約料金 = VALUES(契約料金),
			航送車種区分 = VALUES(航送車種区分),
			航送車種区分名 = VALUES(航送車種区分名),
			見なし距離 = VALUES(見なし距離),
			ferry_srch = VALUES(ferry_srch)
	`

	_, err := r.localDB.Exec(query,
		record.ID, record.UnkoNo, record.UnkoDate, record.ReadDate,
		record.OfficeCode, record.OfficeName, record.VehicleCode, record.VehicleName,
		record.DriverCode1, record.DriverName1, record.TargetDriverClass,
		record.StartTime, record.EndTime, record.FerryCompanyCode, record.FerryCompanyName,
		record.BoardingCode, record.BoardingName, record.ShipNumber,
		record.LandingCode, record.LandingName, record.SettlementClass, record.SettlementName,
		record.StandardFare, record.ContractFare, record.ShipVehicleClass, record.ShipVehicleName,
		record.EstimatedDistance, record.FerrySearch,
	)

	return err
}