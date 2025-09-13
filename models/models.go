package models

import "time"

// ImportRequest represents an import request
type ImportRequest struct {
	FromDate      string `json:"from_date" example:"2025-01-01"`
	ToDate        string `json:"to_date" example:"2025-01-31"`
	EventType     string `json:"event_type,omitempty" example:"運転"`        // For events
	FerryCompany  string `json:"ferry_company,omitempty" example:"東京フェリー"` // For ferry rows
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	Success      bool      `json:"success" example:"true"`
	ImportedRows int       `json:"imported_rows" example:"150"`
	Message      string    `json:"message" example:"Imported 150 rows successfully"`
	ImportedAt   time.Time `json:"imported_at" example:"2025-01-13T15:04:05Z"`
	Errors       []string  `json:"errors,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid request parameters"`
}

// DtakoRow represents a row record from production
type DtakoRow struct {
	ID          string     `json:"id" example:"row-123"`
	UnkoNo      string     `json:"unko_no" example:"2025010101"`      // 運行NO
	Date        time.Time  `json:"date" example:"2025-01-13T00:00:00Z"`
	VehicleNo   string     `json:"vehicle_no" example:"vehicle-001"`
	DriverCode  string     `json:"driver_code" example:"driver-123"`
	RouteCode   string     `json:"route_code" example:"route-A"`
	Distance    float64    `json:"distance" example:"123.45"`
	FuelAmount  float64    `json:"fuel_amount" example:"45.67"`
	CreatedAt   *time.Time `json:"created_at,omitempty" example:"2025-01-13T15:04:05Z"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" example:"2025-01-13T15:04:05Z"`
}

// DtakoEvent represents an event record from production
type DtakoEvent struct {
	ID          string     `json:"id" example:"event-456"`
	UnkoNo      string     `json:"unko_no,omitempty" example:"2025010101"`      // 運行NO - links to DtakoRow
	EventDate   time.Time  `json:"event_date" example:"2025-01-13T10:30:00Z"`
	EventType   string     `json:"event_type" example:"運転"`
	VehicleNo   string     `json:"vehicle_no" example:"vehicle-001"`
	DriverCode  string     `json:"driver_code" example:"driver-123"`
	Description string     `json:"description" example:"Started driving from depot"`
	Latitude    *float64   `json:"latitude,omitempty" example:"35.6762"`
	Longitude   *float64   `json:"longitude,omitempty" example:"139.6503"`
	CreatedAt   *time.Time `json:"created_at,omitempty" example:"2025-01-13T15:04:05Z"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" example:"2025-01-13T15:04:05Z"`
}

// DtakoFerryRow represents a ferry row record from production
type DtakoFerryRow struct {
	ID                int       `json:"id" example:"1"`
	UnkoNo            string    `json:"unko_no" example:"2025010101"`                      // 運行NO
	UnkoDate          time.Time `json:"unko_date" example:"2025-01-13T00:00:00Z"`         // 運行日
	ReadDate          time.Time `json:"read_date" example:"2025-01-13T00:00:00Z"`         // 読取日
	OfficeCode        int       `json:"office_code" example:"1"`                          // 事業所CD
	OfficeName        string    `json:"office_name" example:"東京事業所"`                      // 事業所名
	VehicleCode       int       `json:"vehicle_code" example:"101"`                       // 車輌CD
	VehicleName       string    `json:"vehicle_name" example:"トラック1号"`                    // 車輌名
	DriverCode1       int       `json:"driver_code_1" example:"1001"`                     // 乗務員CD1
	DriverName1       string    `json:"driver_name_1" example:"山田太郎"`                     // 乗務員名１
	TargetDriverClass int       `json:"target_driver_class" example:"1"`                  // 対象乗務員区分
	StartTime         time.Time `json:"start_time" example:"2025-01-13T08:00:00Z"`        // 開始日時
	EndTime           time.Time `json:"end_time" example:"2025-01-13T12:00:00Z"`          // 終了日時
	FerryCompanyCode  int       `json:"ferry_company_code" example:"1"`                   // フェリー会社CD
	FerryCompanyName  string    `json:"ferry_company_name" example:"東京フェリー"`              // フェリー会社名
	BoardingCode      int       `json:"boarding_code" example:"1"`                        // 乗場CD
	BoardingName      string    `json:"boarding_name" example:"東京港"`                      // 乗場名
	ShipNumber        string    `json:"ship_number" example:"1便"`                         // 便
	LandingCode       int       `json:"landing_code" example:"2"`                         // 降場CD
	LandingName       string    `json:"landing_name" example:"大阪港"`                       // 降場名
	SettlementClass  int       `json:"settlement_class" example:"1"`                     // 精算区分
	SettlementName    string    `json:"settlement_name" example:"現金"`                     // 精算区分名
	StandardFare      int       `json:"standard_fare" example:"10000"`                    // 標準料金
	ContractFare      int       `json:"contract_fare" example:"8000"`                     // 契約料金
	ShipVehicleClass  int       `json:"ship_vehicle_class" example:"1"`                   // 航送車種区分
	ShipVehicleName   string    `json:"ship_vehicle_name" example:"大型車"`                  // 航送車種区分名
	EstimatedDistance int       `json:"estimated_distance" example:"500"`                 // 見なし距離
	FerrySearch       string    `json:"ferry_search,omitempty" example:"東京-大阪"`           // ferry_srch
}