package models

import "time"

// ImportRequest represents an import request
type ImportRequest struct {
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
	EventType string `json:"event_type,omitempty"` // For events
	Route     string `json:"route,omitempty"`      // For ferry
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	Success      bool      `json:"success"`
	ImportedRows int       `json:"imported_rows"`
	Message      string    `json:"message"`
	ImportedAt   time.Time `json:"imported_at"`
	Errors       []string  `json:"errors,omitempty"`
}

// DtakoRow represents a row record from production
type DtakoRow struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	VehicleNo   string    `json:"vehicle_no"`
	DriverCode  string    `json:"driver_code"`
	RouteCode   string    `json:"route_code"`
	Distance    float64   `json:"distance"`
	FuelAmount  float64   `json:"fuel_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DtakoEvent represents an event record from production
type DtakoEvent struct {
	ID          string    `json:"id"`
	EventDate   time.Time `json:"event_date"`
	EventType   string    `json:"event_type"`
	VehicleNo   string    `json:"vehicle_no"`
	DriverCode  string    `json:"driver_code"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude,omitempty"`
	Longitude   float64   `json:"longitude,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DtakoFerry represents a ferry record from production
type DtakoFerry struct {
	ID            string    `json:"id"`
	Date          time.Time `json:"date"`
	Route         string    `json:"route"`
	VehicleNo     string    `json:"vehicle_no"`
	DriverCode    string    `json:"driver_code"`
	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`
	Passengers    int       `json:"passengers"`
	Vehicles      int       `json:"vehicles"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}