package repositories

import (
	"context"
	"database/sql"
	"log"
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
	log.Printf("🔍 DEBUG: GetByDateRange START - from=%v, to=%v, eventType=%s, unkoNo=%s", from, to, eventType, unkoNo)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var db *sql.DB = r.prodDB
	if db == nil {
		db = r.localDB
	}

	// 最初にテーブル存在確認
	log.Printf("🔍 DEBUG: Testing table access")
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM dtako_events").Scan(&count)
	if err != nil {
		log.Printf("❌ ERROR: Table access failed: %v", err)
		return []models.DtakoEvent{}, err
	}
	log.Printf("✅ SUCCESS: Table has %d rows", count)

	// 根本問題修正: 実際のテーブル構造に合わせたクエリ
	// - created_at, updated_at カラムを除外
	// - DATE()関数を使わず直接日時比較
	// - 実際のカラム型に合わせたスキャン
	query := `
		SELECT
			id,
			COALESCE(運行NO, '') as unko_no,
			開始日時 as event_date,
			イベント名 as event_type,
			CAST(車輌CD AS CHAR) as vehicle_no,
			CAST(対象乗務員CD AS CHAR) as driver_code,
			COALESCE(備考, '') as description,
			開始GPS緯度,
			開始GPS経度
		FROM dtako_events
		WHERE 開始日時 >= ? AND 開始日時 < DATE_ADD(?, INTERVAL 1 DAY)
	`

	args := []interface{}{from.Format("2006-01-02"), to.Format("2006-01-02")}

	if eventType != "" {
		query += " AND イベント名 = ?"
		args = append(args, eventType)
	}

	if unkoNo != "" {
		query += " AND 運行NO = ?"
		args = append(args, unkoNo)
	}

	query += " ORDER BY 開始日時 DESC LIMIT 100"

	log.Printf("🔍 DEBUG: Executing optimized query")
	log.Printf("🔍 DEBUG: Query: %s", query)
	log.Printf("🔍 DEBUG: Args: %v", args)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("❌ ERROR: Query failed: %v", err)
		return []models.DtakoEvent{}, err
	}
	defer rows.Close()

	log.Printf("🔍 DEBUG: Query executed successfully, processing rows")

	results := []models.DtakoEvent{}
	rowCount := 0

	for rows.Next() {
		rowCount++
		log.Printf("🔍 DEBUG: Processing row %d", rowCount)

		if rowCount > 100 {
			log.Printf("⚠️ WARNING: Too many rows (%d), breaking loop", rowCount)
			break
		}

		var event models.DtakoEvent
		var latBigint, lngBigint sql.NullInt64

		// 根本修正: created_at, updated_at を除外
		err := rows.Scan(
			&event.ID,
			&event.UnkoNo,
			&event.EventDate,
			&event.EventType,
			&event.VehicleNo,
			&event.DriverCode,
			&event.Description,
			&latBigint,
			&lngBigint,
		)
		if err != nil {
			log.Printf("❌ ERROR: Row scan failed at row %d: %v", rowCount, err)
			return []models.DtakoEvent{}, err
		}

		// GPS座標変換
		if latBigint.Valid {
			lat := float64(latBigint.Int64) / 1000000.0
			event.Latitude = &lat
		}
		if lngBigint.Valid {
			lng := float64(lngBigint.Int64) / 1000000.0
			event.Longitude = &lng
		}

		// created_at, updated_at はnilのままにする（実際のテーブルには存在しない）
		event.CreatedAt = nil
		event.UpdatedAt = nil

		results = append(results, event)
		log.Printf("🔍 DEBUG: Row %d processed successfully", rowCount)
	}

	log.Printf("✅ SUCCESS: GetByDateRange completed - %d rows processed", rowCount)
	return results, nil
}

// GetByID retrieves a specific event by ID from local database
func (r *DtakoEventsRepository) GetByID(id string) (*models.DtakoEvent, error) {
	log.Printf("🔍 DEBUG: GetByID START - id=%s", id)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var db *sql.DB = r.prodDB
	if db == nil {
		db = r.localDB
	}

	// 根本修正: created_at, updated_at を除外したクエリ
	query := `
		SELECT
			id,
			COALESCE(運行NO, '') as unko_no,
			開始日時 as event_date,
			イベント名 as event_type,
			CAST(車輌CD AS CHAR) as vehicle_no,
			CAST(対象乗務員CD AS CHAR) as driver_code,
			COALESCE(備考, '') as description,
			開始GPS緯度,
			開始GPS経度
		FROM dtako_events
		WHERE id = ?
	`

	var event models.DtakoEvent
	var latBigint, lngBigint sql.NullInt64

	log.Printf("🔍 DEBUG: Executing GetByID query")
	err := db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.UnkoNo,
		&event.EventDate,
		&event.EventType,
		&event.VehicleNo,
		&event.DriverCode,
		&event.Description,
		&latBigint,
		&lngBigint,
	)

	if err != nil {
		log.Printf("❌ ERROR: GetByID query failed: %v", err)
		return nil, err
	}

	// GPS座標変換
	if latBigint.Valid {
		lat := float64(latBigint.Int64) / 1000000.0
		event.Latitude = &lat
	}
	if lngBigint.Valid {
		lng := float64(lngBigint.Int64) / 1000000.0
		event.Longitude = &lng
	}

	// created_at, updated_at はnilのままにする
	event.CreatedAt = nil
	event.UpdatedAt = nil

	log.Printf("✅ SUCCESS: GetByID completed")
	return &event, nil
}

// FetchFromProduction fetches event data from production database
func (r *DtakoEventsRepository) FetchFromProduction(from, to time.Time, eventType string) ([]models.DtakoEvent, error) {
	if r.prodDB == nil {
		return []models.DtakoEvent{}, nil
	}

	log.Printf("🔍 DEBUG: FetchFromProduction START")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT
			id,
			COALESCE(運行NO, '') as unko_no,
			開始日時 as event_date,
			イベント名 as event_type,
			CAST(車輌CD AS CHAR) as vehicle_no,
			CAST(対象乗務員CD AS CHAR) as driver_code,
			COALESCE(備考, '') as description,
			開始GPS緯度,
			開始GPS経度
		FROM dtako_events
		WHERE 開始日時 >= ? AND 開始日時 < DATE_ADD(?, INTERVAL 1 DAY)
	`
	args := []interface{}{from.Format("2006-01-02"), to.Format("2006-01-02")}

	if eventType != "" {
		query += " AND イベント名 = ?"
		args = append(args, eventType)
	}

	query += " ORDER BY 開始日時 DESC LIMIT 100"

	rows, err := r.prodDB.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("❌ ERROR: FetchFromProduction query failed: %v", err)
		return []models.DtakoEvent{}, err
	}
	defer rows.Close()

	results := []models.DtakoEvent{}
	rowCount := 0

	for rows.Next() {
		rowCount++
		if rowCount > 100 {
			log.Printf("⚠️ WARNING: Too many rows in FetchFromProduction, breaking")
			break
		}

		var event models.DtakoEvent
		var latBigint, lngBigint sql.NullInt64

		err := rows.Scan(
			&event.ID,
			&event.UnkoNo,
			&event.EventDate,
			&event.EventType,
			&event.VehicleNo,
			&event.DriverCode,
			&event.Description,
			&latBigint,
			&lngBigint,
		)
		if err != nil {
			log.Printf("❌ ERROR: FetchFromProduction scan failed: %v", err)
			return []models.DtakoEvent{}, err
		}

		// GPS座標変換
		if latBigint.Valid {
			lat := float64(latBigint.Int64) / 1000000.0
			event.Latitude = &lat
		}
		if lngBigint.Valid {
			lng := float64(lngBigint.Int64) / 1000000.0
			event.Longitude = &lng
		}

		// created_at, updated_at はnilのまま
		event.CreatedAt = nil
		event.UpdatedAt = nil

		results = append(results, event)
	}

	log.Printf("✅ SUCCESS: FetchFromProduction completed - %d rows", rowCount)
	return results, nil
}

// Insert inserts an event into local database
func (r *DtakoEventsRepository) Insert(event *models.DtakoEvent) error {
	// 実際のテーブル構造に合わせたINSERT
	query := `
		INSERT INTO dtako_events (
			id, 運行NO, 読取日, 車輌CD, 車輌CC, 開始日時, 終了日時,
			イベント名, 対象乗務員CD, 対象乗務員区分, 乗務員CD1,
			開始走行距離, 終了走行距離, 区間時間, 区間距離,
			開始市町村名, 終了市町村名, 開始場所名, 終了場所名,
			開始GPS緯度, 開始GPS経度, 備考
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		    運行NO = VALUES(運行NO),
		    読取日 = VALUES(読取日),
		    開始日時 = VALUES(開始日時),
		    終了日時 = VALUES(終了日時),
		    イベント名 = VALUES(イベント名),
		    備考 = VALUES(備考)
	`

	// デフォルト値
	readDate := event.EventDate
	vehicleCD := 1
	vehicleCC := "001100"
	driverCode := 0
	driverKubun := 0
	driverCD1 := 0
	startDistance := 0.0
	endDistance := 0.0
	sectionTime := 0
	sectionDistance := 0.0
	startCity := ""
	endCity := ""
	startPlace := ""
	endPlace := ""

	if event.VehicleNo != "" {
		vehicleCD = 1
	}
	if event.DriverCode != "" {
		driverCode = 1
	}

	endDateTime := event.EventDate

	var description sql.NullString
	if event.Description != "" {
		description = sql.NullString{String: event.Description, Valid: true}
	}

	var latitude, longitude sql.NullInt64
	if event.Latitude != nil {
		latitude = sql.NullInt64{Int64: int64(*event.Latitude * 1000000), Valid: true}
	}
	if event.Longitude != nil {
		longitude = sql.NullInt64{Int64: int64(*event.Longitude * 1000000), Valid: true}
	}

	_, err := r.localDB.Exec(query,
		event.ID, event.UnkoNo, readDate, vehicleCD, vehicleCC,
		event.EventDate, endDateTime, event.EventType,
		driverCode, driverKubun, driverCD1,
		startDistance, endDistance, sectionTime, sectionDistance,
		startCity, endCity, startPlace, endPlace,
		latitude, longitude, description,
	)

	return err
}