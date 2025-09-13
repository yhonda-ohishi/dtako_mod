-- Schema for dtako_rows table
CREATE TABLE IF NOT EXISTS dtako_rows (
    id VARCHAR(50) PRIMARY KEY,
    unko_no VARCHAR(255) UNIQUE NOT NULL,  -- 運行NO
    date DATE NOT NULL,
    vehicle_no VARCHAR(20),
    driver_code VARCHAR(20),
    route_code VARCHAR(20),
    distance DECIMAL(10, 2),
    fuel_amount DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_date (date),
    INDEX idx_unko_no (unko_no),
    INDEX idx_vehicle (vehicle_no),
    INDEX idx_driver (driver_code)
);

-- Schema for dtako_events table
CREATE TABLE IF NOT EXISTS dtako_events (
    id VARCHAR(50) PRIMARY KEY,
    unko_no VARCHAR(255),  -- 運行NO - links to dtako_rows
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
    INDEX idx_vehicle (vehicle_no),
    INDEX idx_unko_no (unko_no),
    FOREIGN KEY (unko_no) REFERENCES dtako_rows(unko_no) ON DELETE CASCADE
);

-- Schema for dtako_ferry_rows table
CREATE TABLE IF NOT EXISTS dtako_ferry_rows (
    id INT PRIMARY KEY AUTO_INCREMENT,
    運行NO VARCHAR(23) NOT NULL,
    運行日 DATE NOT NULL,
    読取日 DATE NOT NULL,
    事業所CD INT NOT NULL,
    事業所名 VARCHAR(20) NOT NULL,
    車輌CD INT NOT NULL,
    車輌名 VARCHAR(20) NOT NULL,
    乗務員CD1 INT NOT NULL,
    乗務員名１ VARCHAR(20) NOT NULL,
    対象乗務員区分 INT NOT NULL,
    開始日時 DATETIME NOT NULL,
    終了日時 DATETIME NOT NULL,
    フェリー会社CD INT NOT NULL,
    フェリー会社名 VARCHAR(20) NOT NULL,
    乗場CD INT NOT NULL,
    乗場名 VARCHAR(20) NOT NULL,
    便 VARCHAR(10) NOT NULL,
    降場CD INT NOT NULL,
    降場名 VARCHAR(20) NOT NULL,
    精算区分 INT NOT NULL,
    精算区分名 VARCHAR(20) NOT NULL,
    標準料金 INT NOT NULL,
    契約料金 INT NOT NULL,
    航送車種区分 INT NOT NULL,
    航送車種区分名 VARCHAR(20) NOT NULL,
    見なし距離 INT NOT NULL,
    ferry_srch VARCHAR(60) DEFAULT NULL,
    INDEX idx_unko_no (運行NO),
    INDEX idx_unko_date (運行日),
    INDEX idx_ferry_company (フェリー会社名)
);