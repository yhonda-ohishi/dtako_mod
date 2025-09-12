-- Schema for dtako_rows table
CREATE TABLE IF NOT EXISTS dtako_rows (
    id VARCHAR(50) PRIMARY KEY,
    date DATE NOT NULL,
    vehicle_no VARCHAR(20),
    driver_code VARCHAR(20),
    route_code VARCHAR(20),
    distance DECIMAL(10, 2),
    fuel_amount DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_date (date),
    INDEX idx_vehicle (vehicle_no),
    INDEX idx_driver (driver_code)
);

-- Schema for dtako_events table
CREATE TABLE IF NOT EXISTS dtako_events (
    id VARCHAR(50) PRIMARY KEY,
    event_date DATETIME NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    vehicle_no VARCHAR(20),
    driver_code VARCHAR(20),
    description TEXT,
    latitude DECIMAL(10, 6),
    longitude DECIMAL(10, 6),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_event_date (event_date),
    INDEX idx_event_type (event_type),
    INDEX idx_vehicle (vehicle_no)
);

-- Schema for dtako_ferry table
CREATE TABLE IF NOT EXISTS dtako_ferry (
    id VARCHAR(50) PRIMARY KEY,
    date DATE NOT NULL,
    route VARCHAR(50) NOT NULL,
    vehicle_no VARCHAR(20),
    driver_code VARCHAR(20),
    departure_time TIME,
    arrival_time TIME,
    passengers INT,
    vehicles INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_date (date),
    INDEX idx_route (route),
    INDEX idx_vehicle (vehicle_no)
);