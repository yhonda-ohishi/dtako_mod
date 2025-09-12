-- Database schema for dtako_mod

-- Create database if not exists
CREATE DATABASE IF NOT EXISTS dtako_local DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE dtako_local;

-- dtako_rows table
CREATE TABLE IF NOT EXISTS dtako_rows (
    id VARCHAR(255) PRIMARY KEY,
    date DATE NOT NULL,
    vehicle_no VARCHAR(50),
    driver_code VARCHAR(50),
    route_code VARCHAR(50),
    distance DECIMAL(10,2),
    fuel_amount DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_date (date),
    INDEX idx_vehicle (vehicle_no),
    INDEX idx_driver (driver_code),
    INDEX idx_route (route_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- dtako_events table
CREATE TABLE IF NOT EXISTS dtako_events (
    id VARCHAR(255) PRIMARY KEY,
    event_date DATETIME NOT NULL,
    event_type VARCHAR(50),
    vehicle_no VARCHAR(50),
    driver_code VARCHAR(50),
    description TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_event_date (event_date),
    INDEX idx_event_type (event_type),
    INDEX idx_vehicle (vehicle_no),
    INDEX idx_driver (driver_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- dtako_ferry table
CREATE TABLE IF NOT EXISTS dtako_ferry (
    id VARCHAR(255) PRIMARY KEY,
    date DATE NOT NULL,
    route VARCHAR(100),
    vehicle_no VARCHAR(50),
    driver_code VARCHAR(50),
    departure_time DATETIME,
    arrival_time DATETIME,
    passengers INT,
    vehicles INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_date (date),
    INDEX idx_route (route),
    INDEX idx_vehicle (vehicle_no),
    INDEX idx_driver (driver_code),
    INDEX idx_departure (departure_time),
    INDEX idx_arrival (arrival_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;