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
	// デフォルトで本番DB構造（日本語カラム名）を使用
	var query string
	var db *sql.DB
	if os.Getenv("DTAKO_ENV") == "test_english" {
		// テスト用英語カラム環境
		db = r.localDB
		query = `
			SELECT id, COALESCE(運行NO, ''), 開始日時 as event_date, イベント名 as event_type,
			       CAST(車輌CD AS CHAR) as vehicle_no, CAST(対象乗務員CD AS CHAR) as driver_code,
			       COALESCE(備考, '') as description, 開始GPS緯度 as latitude, 開始GPS経度 as longitude,
			       NULL as created_at, NULL as updated_at
			FROM dtako_events
			WHERE DATE(開始日時) BETWEEN ? AND ?
		`
	} else {
		// 本番・デフォルト環境（日本語カラム名）
		db = r.prodDB
		if db == nil {
			db = r.localDB // フォールバック
		}
		query = `
			SELECT id, COALESCE(運行NO, ''), 開始日時 as event_date, イベント名 as event_type,
			       CAST(車輌CD AS CHAR) as vehicle_no, CAST(対象乗務員CD AS CHAR) as driver_code,
			       COALESCE(備考, '') as description, 開始GPS緯度 as latitude, 開始GPS経度 as longitude,
			       NULL as created_at, NULL as updated_at
			FROM dtako_events
			WHERE DATE(開始日時) BETWEEN ? AND ?
		`
	}

	args := []interface{}{from, to}

	if eventType != "" {
		query += " AND イベント名 = ?"
		args = append(args, eventType)
	}

	if unkoNo != "" {
		query += " AND 運行NO = ?"
		args = append(args, unkoNo)
	}

	query += " ORDER BY 開始日時 DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return []models.DtakoEvent{}, err
	}
	defer rows.Close()

	results := []models.DtakoEvent{}
	for rows.Next() {
		var event models.DtakoEvent

		// すべての環境で日本語カラムを使用するため、bigint型GPS座標の変換が必要
		var latBigint, lngBigint sql.NullInt64

		err := rows.Scan(
			&event.ID, &event.UnkoNo, &event.EventDate, &event.EventType, &event.VehicleNo,
			&event.DriverCode, &event.Description, &latBigint, &lngBigint,
			&event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return []models.DtakoEvent{}, err
		}

		// 緯度経度の型変換（bigint → float64）
		if latBigint.Valid {
			lat := float64(latBigint.Int64) / 1000000.0
			event.Latitude = &lat
		}
		if lngBigint.Valid {
			lng := float64(lngBigint.Int64) / 1000000.0
			event.Longitude = &lng
		}

		results = append(results, event)
	}

	return results, nil
}

// GetByID retrieves a specific event by ID from local database
func (r *DtakoEventsRepository) GetByID(id string) (*models.DtakoEvent, error) {
	// デフォルトで本番DB構造（日本語カラム名）を使用
	var query string
	var db *sql.DB
	isEnglish := os.Getenv("DTAKO_ENV") == "test_english"

	if isEnglish {
		// テスト用英語カラム環境
		db = r.localDB
		query = `
			SELECT id, COALESCE(運行NO, ''), 開始日時 as event_date, イベント名 as event_type,
			       CAST(車輌CD AS CHAR) as vehicle_no, CAST(対象乗務員CD AS CHAR) as driver_code,
			       COALESCE(備考, '') as description, 開始GPS緯度 as latitude, 開始GPS経度 as longitude,
			       NULL as created_at, NULL as updated_at
			FROM dtako_events
			WHERE id = ?
		`
	} else {
		// 本番・デフォルト環境（日本語カラム名）
		db = r.prodDB
		if db == nil {
			db = r.localDB // フォールバック
		}
		query = `
			SELECT id, COALESCE(運行NO, ''), 開始日時 as event_date, イベント名 as event_type,
			       CAST(車輌CD AS CHAR) as vehicle_no, CAST(対象乗務員CD AS CHAR) as driver_code,
			       COALESCE(備考, '') as description, 開始GPS緯度 as latitude, 開始GPS経度 as longitude,
			       NULL as created_at, NULL as updated_at
			FROM dtako_events
			WHERE id = ?
		`
	}

	var event models.DtakoEvent

	// すべての環境で日本語カラムを使用するため、bigint型GPS座標の変換が必要
	var latBigint, lngBigint sql.NullInt64

	err := db.QueryRow(query, id).Scan(
		&event.ID, &event.UnkoNo, &event.EventDate, &event.EventType, &event.VehicleNo,
		&event.DriverCode, &event.Description, &latBigint, &lngBigint,
		&event.CreatedAt, &event.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// 緯度経度の型変換（bigint → float64）
	if latBigint.Valid {
		lat := float64(latBigint.Int64) / 1000000.0
		event.Latitude = &lat
	}
	if lngBigint.Valid {
		lng := float64(lngBigint.Int64) / 1000000.0
		event.Longitude = &lng
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
	var eventTypeColumn string
	var dateColumn string

	isEnglish := os.Getenv("DTAKO_ENV") == "test_english"

	if isEnglish {
		// テスト用英語カラム環境
		eventTypeColumn = "イベント名"
		dateColumn = "開始日時"
		query = `
			SELECT id, COALESCE(運行NO, ''), 開始日時 as event_date, イベント名 as event_type,
			       CAST(車輌CD AS CHAR) as vehicle_no, CAST(対象乗務員CD AS CHAR) as driver_code,
			       COALESCE(備考, '') as description, 開始GPS緯度 as latitude, 開始GPS経度 as longitude,
			       NULL as created_at, NULL as updated_at
			FROM dtako_events
			WHERE DATE(開始日時) BETWEEN ? AND ?
		`
	} else {
		// デフォルト（日本語カラム名）
		eventTypeColumn = "イベント名"
		dateColumn = "開始日時"
		query = `
			SELECT id, COALESCE(運行NO, ''), 開始日時 as event_date, イベント名 as event_type,
			       CAST(車輌CD AS CHAR) as vehicle_no, CAST(対象乗務員CD AS CHAR) as driver_code,
			       COALESCE(備考, '') as description, 開始GPS緯度 as latitude, 開始GPS経度 as longitude,
			       NULL as created_at, NULL as updated_at
			FROM dtako_events
			WHERE DATE(開始日時) BETWEEN ? AND ?
		`
	}
	args := []interface{}{from, to}

	if eventType != "" {
		query += " AND " + eventTypeColumn + " = ?"
		args = append(args, eventType)
	}

	query += " ORDER BY " + dateColumn + " DESC"

	rows, err := r.prodDB.Query(query, args...)
	if err != nil {
		return []models.DtakoEvent{}, err
	}
	defer rows.Close()

	results := []models.DtakoEvent{}
	for rows.Next() {
		var event models.DtakoEvent
		var latBigint, lngBigint sql.NullInt64

		// すべての環境で日本語カラムを使用するため、bigint型GPS座標の変換が必要
		err := rows.Scan(
			&event.ID, &event.UnkoNo, &event.EventDate, &event.EventType, &event.VehicleNo,
			&event.DriverCode, &event.Description, &latBigint, &lngBigint,
			&event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return []models.DtakoEvent{}, err
		}

		// 緯度経度の型変換（bigint → float64）
		if latBigint.Valid {
			lat := float64(latBigint.Int64) / 1000000.0
			event.Latitude = &lat
		}
		if lngBigint.Valid {
			lng := float64(lngBigint.Int64) / 1000000.0
			event.Longitude = &lng
		}

		results = append(results, event)
	}

	return results, nil
}

// Insert inserts an event into local database
func (r *DtakoEventsRepository) Insert(event *models.DtakoEvent) error {
	// ローカルDBの実際のカラム構造に合わせる
	// 必須カラム: id, 運行NO, 読取日, 車輌CD, 車輌CC, 開始日時, 終了日時, イベント名
	query := `
		INSERT INTO dtako_events (id, 運行NO, 読取日, 車輌CD, 車輌CC, 開始日時, 終了日時,
		                         イベント名, 対象乗務員CD, 対象乗務員区分, 乗務員CD1,
		                         開始走行距離, 終了走行距離, 区間時間, 区間距離,
		                         開始市町村名, 終了市町村名, 開始場所名, 終了場所名,
		                         開始GPS緯度, 開始GPS経度, 備考)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		    運行NO = VALUES(運行NO),
		    読取日 = VALUES(読取日),
		    開始日時 = VALUES(開始日時),
		    終了日時 = VALUES(終了日時),
		    イベント名 = VALUES(イベント名),
		    備考 = VALUES(備考)
	`

	// デフォルト値の設定
	readDate := event.EventDate                 // 読取日は EventDate を使用
	vehicleCD := 1                              // 車輌CD
	vehicleCC := "001100"                       // 車輌CC
	driverCode := 0                             // 対象乗務員CD
	driverKubun := 0                            // 対象乗務員区分
	driverCD1 := 0                              // 乗務員CD1
	startDistance := 0.0                        // 開始走行距離
	endDistance := 0.0                          // 終了走行距離
	sectionTime := 0                            // 区間時間
	sectionDistance := 0.0                      // 区間距離
	startCity := ""                             // 開始市町村名
	endCity := ""                               // 終了市町村名
	startPlace := ""                            // 開始場所名
	endPlace := ""                              // 終了場所名

	if event.VehicleNo != "" {
		// VehicleNoから変換
		vehicleCD = 1
	}
	if event.DriverCode != "" {
		driverCode = 1
	}

	// イベントの終了日時（開始日時と同じにする）
	endDateTime := event.EventDate

	// 備考欄にdescriptionを設定
	var description sql.NullString
	if event.Description != "" {
		description = sql.NullString{String: event.Description, Valid: true}
	}

	// NULLable な緯度経度の処理
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