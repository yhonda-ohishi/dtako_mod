package models

import "time"

// ImportRequest represents an import request
type ImportRequest struct {
	FromDate  string `json:"from_date" example:"2025-01-01"`
	ToDate    string `json:"to_date" example:"2025-01-31"`
	EventType string `json:"event_type,omitempty" example:"運転"` // For events
	Route     string `json:"route,omitempty" example:"Tokyo-Osaka"`      // For ferry
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
	ID          string    `json:"id" example:"row-123"`
	UnkoNo      string    `json:"unko_no" example:"2025010101"`      // 運行NO
	Date        time.Time `json:"date" example:"2025-01-13T00:00:00Z"`
	VehicleNo   string    `json:"vehicle_no" example:"vehicle-001"`
	DriverCode  string    `json:"driver_code" example:"driver-123"`
	RouteCode   string    `json:"route_code" example:"route-A"`
	Distance    float64   `json:"distance" example:"123.45"`
	FuelAmount  float64   `json:"fuel_amount" example:"45.67"`
	CreatedAt   time.Time `json:"created_at" example:"2025-01-13T15:04:05Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-01-13T15:04:05Z"`
}

// DtakoEvent represents an event record from production
type DtakoEvent struct {
	ID          string    `json:"id" example:"event-456"`
	UnkoNo      string    `json:"unko_no,omitempty" example:"2025010101"`      // 運行NO - links to DtakoRow
	EventDate   time.Time `json:"event_date" example:"2025-01-13T10:30:00Z"`
	EventType   string    `json:"event_type" example:"運転"`
	VehicleNo   string    `json:"vehicle_no" example:"vehicle-001"`
	DriverCode  string    `json:"driver_code" example:"driver-123"`
	Description string    `json:"description" example:"Started driving from depot"`
	Latitude    float64   `json:"latitude,omitempty" example:"35.6762"`
	Longitude   float64   `json:"longitude,omitempty" example:"139.6503"`
	CreatedAt   time.Time `json:"created_at" example:"2025-01-13T15:04:05Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-01-13T15:04:05Z"`
}

// DtakoFerry represents a ferry record from production
type DtakoFerry struct {
	ID            string    `json:"id" example:"ferry-789"`
	Date          time.Time `json:"date" example:"2025-01-13T00:00:00Z"`
	Route         string    `json:"route" example:"Tokyo-Osaka"`
	VehicleNo     string    `json:"vehicle_no" example:"vehicle-001"`
	DriverCode    string    `json:"driver_code" example:"driver-123"`
	DepartureTime time.Time `json:"departure_time" example:"2025-01-13T08:00:00Z"`
	ArrivalTime   time.Time `json:"arrival_time" example:"2025-01-13T12:00:00Z"`
	Passengers    int       `json:"passengers" example:"150"`
	Vehicles      int       `json:"vehicles" example:"30"`
	CreatedAt     time.Time `json:"created_at" example:"2025-01-13T15:04:05Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2025-01-13T15:04:05Z"`
}