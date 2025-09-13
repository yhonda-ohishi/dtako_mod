-- データベース作成（存在しない場合）
CREATE DATABASE IF NOT EXISTS dtako_local CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS dtako_test_prod CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- dtako_local データベースのテーブル作成
USE dtako_local;

-- dtako_rows テーブル (実際のDBスキーマに合わせて日本語カラム名を使用)
CREATE TABLE IF NOT EXISTS dtako_rows (
    id VARCHAR(24) PRIMARY KEY,
    運行NO VARCHAR(23) NOT NULL UNIQUE,
    読取日 DATE NOT NULL,
    運行日 DATE NOT NULL,
    車輌CD INT NOT NULL,
    車輌CC VARCHAR(6) NOT NULL,
    乗務員CD1 INT,
    対象乗務員区分 INT NOT NULL DEFAULT 0,
    対象乗務員CD INT NOT NULL DEFAULT 0,
    出社日時 DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
    退社日時 DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
    出庫日時 DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
    帰庫日時 DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
    出庫メーター DOUBLE NOT NULL DEFAULT 0,
    帰庫メーター DOUBLE NOT NULL DEFAULT 0,
    総走行距離 DOUBLE NOT NULL DEFAULT 0,
    実車走行距離 DOUBLE,
    行先市町村名 VARCHAR(40),
    行先場所名 VARCHAR(40),
    一般道運転時間 INT NOT NULL DEFAULT 0,
    高速道運転時間 INT NOT NULL DEFAULT 0,
    バイパス運転時間 INT NOT NULL DEFAULT 0,
    実車走行時間 INT NOT NULL DEFAULT 0,
    空車走行時間 INT NOT NULL DEFAULT 0,
    作業１時間 INT NOT NULL DEFAULT 0,
    作業２時間 INT NOT NULL DEFAULT 0,
    作業３時間 INT NOT NULL DEFAULT 0,
    作業４時間 INT NOT NULL DEFAULT 0,
    状態１距離 DOUBLE NOT NULL DEFAULT 0,
    状態１時間 INT NOT NULL DEFAULT 0,
    状態２距離 DOUBLE NOT NULL DEFAULT 0,
    状態２時間 INT NOT NULL DEFAULT 0,
    状態３距離 DOUBLE NOT NULL DEFAULT 0,
    状態３時間 INT NOT NULL DEFAULT 0,
    状態４距離 DOUBLE NOT NULL DEFAULT 0,
    状態４時間 INT NOT NULL DEFAULT 0,
    状態５距離 DOUBLE NOT NULL DEFAULT 0,
    状態５時間 INT NOT NULL DEFAULT 0,
    自社主燃料 DOUBLE NOT NULL DEFAULT 0,
    自社主添加剤 DOUBLE NOT NULL DEFAULT 0,
    他社主燃料 DOUBLE NOT NULL DEFAULT 0,
    他社主添加剤 DOUBLE NOT NULL DEFAULT 0,
    アイドリング時間 BIGINT NOT NULL DEFAULT 0,
    アイドリング時間回数 INT NOT NULL DEFAULT 0,
    総合評価点 INT,
    安全評価点 INT,
    経済評価点 INT,
    INDEX idx_運行日 (運行日),
    INDEX idx_運行NO (運行NO),
    INDEX idx_車輌CD (車輌CD)
);

-- dtako_events テーブル
CREATE TABLE IF NOT EXISTS dtako_events (
    id VARCHAR(50) PRIMARY KEY,
    運行NO VARCHAR(23),
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
    INDEX idx_運行NO (運行NO),
    FOREIGN KEY (運行NO) REFERENCES dtako_rows(運行NO) ON DELETE CASCADE
);

-- dtako_ferry_rows テーブル
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

-- テストデータの投入（既存データがある場合は置き換え）
-- dtako_rows のテストデータ
-- 最小限の必須カラムでテストデータを追加
REPLACE INTO dtako_rows (
    id, 運行NO, 読取日, 運行日, 車輌CD, 車輌CC,
    対象乗務員区分, 対象乗務員CD,
    出社日時, 退社日時, 出庫日時, 帰庫日時,
    出庫メーター, 帰庫メーター, 総走行距離,
    行先市町村名, 自社主燃料
) VALUES
('ROW001', '2024011501', '2024-01-15', '2024-01-15', 1, '001100',
 0, 1,
 '2024-01-15 08:00:00', '2024-01-15 17:00:00', '2024-01-15 08:30:00', '2024-01-15 16:30:00',
 1000, 1150.5, 150.5,
 '東京都', 20.3),
('ROW002', '2024011502', '2024-01-15', '2024-01-15', 2, '002200',
 0, 2,
 '2024-01-15 08:00:00', '2024-01-15 17:00:00', '2024-01-15 08:30:00', '2024-01-15 16:30:00',
 2000, 2200.8, 200.8,
 '大阪府', 25.7),
('ROW003', '2024011601', '2024-01-16', '2024-01-16', 1, '001100',
 0, 1,
 '2024-01-16 08:00:00', '2024-01-16 17:00:00', '2024-01-16 08:30:00', '2024-01-16 16:30:00',
 1150.5, 1325.7, 175.2,
 '名古屋市', 22.1);

-- dtako_events のテストデータ
REPLACE INTO dtako_events (id, 運行NO, event_date, event_type, vehicle_no, driver_code, description, latitude, longitude) VALUES
('EVENT001', '2024011501', '2024-01-15 08:30:00', 'START', 'V001', 'D001', 'Trip started', 35.6762, 139.6503),
('EVENT002', '2024011501', '2024-01-15 12:15:00', 'STOP', 'V001', 'D001', 'Lunch break', 35.6895, 139.6917),
('EVENT003', '2024011501', '2024-01-15 16:45:00', 'END', 'V001', 'D001', 'Trip ended', 35.6762, 139.6503);

-- dtako_ferry_rows のテストデータ（IDは自動生成なので指定しない）
INSERT IGNORE INTO dtako_ferry_rows (運行NO, 運行日, 読取日, 事業所CD, 事業所名, 車輌CD, 車輌名, 乗務員CD1, 乗務員名１, 対象乗務員区分, 開始日時, 終了日時, フェリー会社CD, フェリー会社名, 乗場CD, 乗場名, 便, 降場CD, 降場名, 精算区分, 精算区分名, 標準料金, 契約料金, 航送車種区分, 航送車種区分名, 見なし距離, ferry_srch) VALUES
('2024011501', '2024-01-15', '2024-01-15', 1, '東京事業所', 101, 'トラック1号', 1001, '山田太郎', 1, '2024-01-15 08:00:00', '2024-01-15 10:00:00', 1, '東京フェリー', 1, '東京港', '1便', 2, '大阪港', 1, '現金', 10000, 8000, 1, '大型車', 500, '東京-大阪'),
('2024011502', '2024-01-15', '2024-01-15', 1, '東京事業所', 102, 'トラック2号', 1002, '佐藤次郎', 1, '2024-01-15 09:30:00', '2024-01-15 11:30:00', 1, '東京フェリー', 1, '東京港', '2便', 2, '大阪港', 1, '現金', 10000, 8000, 1, '大型車', 500, '東京-大阪'),
('2024011601', '2024-01-16', '2024-01-16', 1, '東京事業所', 101, 'トラック1号', 1001, '山田太郎', 1, '2024-01-16 08:00:00', '2024-01-16 10:00:00', 2, '関西フェリー', 3, '神戸港', '1便', 4, '高松港', 2, 'クレジット', 8000, 7000, 2, '中型車', 300, '神戸-高松');

-- dtako_test_prod データベースにも同じ構造とデータを作成
USE dtako_test_prod;

-- dtako_rows テーブル (実際のDBスキーマに合わせて日本語カラム名を使用)
CREATE TABLE IF NOT EXISTS dtako_rows (
    id VARCHAR(24) PRIMARY KEY,
    運行NO VARCHAR(23) NOT NULL UNIQUE,
    読取日 DATE NOT NULL,
    運行日 DATE NOT NULL,
    車輌CD INT NOT NULL,
    車輌CC VARCHAR(6) NOT NULL,
    乗務員CD1 INT,
    対象乗務員区分 INT NOT NULL DEFAULT 0,
    対象乗務員CD INT NOT NULL DEFAULT 0,
    出社日時 DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
    退社日時 DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
    出庫日時 DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
    帰庫日時 DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
    出庫メーター DOUBLE NOT NULL DEFAULT 0,
    帰庫メーター DOUBLE NOT NULL DEFAULT 0,
    総走行距離 DOUBLE NOT NULL DEFAULT 0,
    実車走行距離 DOUBLE,
    行先市町村名 VARCHAR(40),
    行先場所名 VARCHAR(40),
    一般道運転時間 INT NOT NULL DEFAULT 0,
    高速道運転時間 INT NOT NULL DEFAULT 0,
    バイパス運転時間 INT NOT NULL DEFAULT 0,
    実車走行時間 INT NOT NULL DEFAULT 0,
    空車走行時間 INT NOT NULL DEFAULT 0,
    作業１時間 INT NOT NULL DEFAULT 0,
    作業２時間 INT NOT NULL DEFAULT 0,
    作業３時間 INT NOT NULL DEFAULT 0,
    作業４時間 INT NOT NULL DEFAULT 0,
    状態１距離 DOUBLE NOT NULL DEFAULT 0,
    状態１時間 INT NOT NULL DEFAULT 0,
    状態２距離 DOUBLE NOT NULL DEFAULT 0,
    状態２時間 INT NOT NULL DEFAULT 0,
    状態３距離 DOUBLE NOT NULL DEFAULT 0,
    状態３時間 INT NOT NULL DEFAULT 0,
    状態４距離 DOUBLE NOT NULL DEFAULT 0,
    状態４時間 INT NOT NULL DEFAULT 0,
    状態５距離 DOUBLE NOT NULL DEFAULT 0,
    状態５時間 INT NOT NULL DEFAULT 0,
    自社主燃料 DOUBLE NOT NULL DEFAULT 0,
    自社主添加剤 DOUBLE NOT NULL DEFAULT 0,
    他社主燃料 DOUBLE NOT NULL DEFAULT 0,
    他社主添加剤 DOUBLE NOT NULL DEFAULT 0,
    アイドリング時間 BIGINT NOT NULL DEFAULT 0,
    アイドリング時間回数 INT NOT NULL DEFAULT 0,
    総合評価点 INT,
    安全評価点 INT,
    経済評価点 INT,
    INDEX idx_運行日 (運行日),
    INDEX idx_運行NO (運行NO),
    INDEX idx_車輌CD (車輌CD)
);

-- dtako_events テーブル
CREATE TABLE IF NOT EXISTS dtako_events (
    id VARCHAR(50) PRIMARY KEY,
    運行NO VARCHAR(23),
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
    INDEX idx_運行NO (運行NO),
    FOREIGN KEY (運行NO) REFERENCES dtako_rows(運行NO) ON DELETE CASCADE
);

-- dtako_ferry_rows テーブル
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

-- テストデータの投入（既存データがある場合は置き換え）
-- dtako_rows のテストデータ
-- 最小限の必須カラムでテストデータを追加
REPLACE INTO dtako_rows (
    id, 運行NO, 読取日, 運行日, 車輌CD, 車輌CC,
    対象乗務員区分, 対象乗務員CD,
    出社日時, 退社日時, 出庫日時, 帰庫日時,
    出庫メーター, 帰庫メーター, 総走行距離,
    行先市町村名, 自社主燃料
) VALUES
('ROW001', '2024011501', '2024-01-15', '2024-01-15', 1, '001100',
 0, 1,
 '2024-01-15 08:00:00', '2024-01-15 17:00:00', '2024-01-15 08:30:00', '2024-01-15 16:30:00',
 1000, 1150.5, 150.5,
 '東京都', 20.3),
('ROW002', '2024011502', '2024-01-15', '2024-01-15', 2, '002200',
 0, 2,
 '2024-01-15 08:00:00', '2024-01-15 17:00:00', '2024-01-15 08:30:00', '2024-01-15 16:30:00',
 2000, 2200.8, 200.8,
 '大阪府', 25.7),
('ROW003', '2024011601', '2024-01-16', '2024-01-16', 1, '001100',
 0, 1,
 '2024-01-16 08:00:00', '2024-01-16 17:00:00', '2024-01-16 08:30:00', '2024-01-16 16:30:00',
 1150.5, 1325.7, 175.2,
 '名古屋市', 22.1);

-- dtako_events のテストデータ
REPLACE INTO dtako_events (id, 運行NO, event_date, event_type, vehicle_no, driver_code, description, latitude, longitude) VALUES
('EVENT001', '2024011501', '2024-01-15 08:30:00', 'START', 'V001', 'D001', 'Trip started', 35.6762, 139.6503),
('EVENT002', '2024011501', '2024-01-15 12:15:00', 'STOP', 'V001', 'D001', 'Lunch break', 35.6895, 139.6917),
('EVENT003', '2024011501', '2024-01-15 16:45:00', 'END', 'V001', 'D001', 'Trip ended', 35.6762, 139.6503);

-- dtako_ferry_rows のテストデータ（IDは自動生成なので指定しない）
INSERT IGNORE INTO dtako_ferry_rows (運行NO, 運行日, 読取日, 事業所CD, 事業所名, 車輌CD, 車輌名, 乗務員CD1, 乗務員名１, 対象乗務員区分, 開始日時, 終了日時, フェリー会社CD, フェリー会社名, 乗場CD, 乗場名, 便, 降場CD, 降場名, 精算区分, 精算区分名, 標準料金, 契約料金, 航送車種区分, 航送車種区分名, 見なし距離, ferry_srch) VALUES
('2024011501', '2024-01-15', '2024-01-15', 1, '東京事業所', 101, 'トラック1号', 1001, '山田太郎', 1, '2024-01-15 08:00:00', '2024-01-15 10:00:00', 1, '東京フェリー', 1, '東京港', '1便', 2, '大阪港', 1, '現金', 10000, 8000, 1, '大型車', 500, '東京-大阪'),
('2024011502', '2024-01-15', '2024-01-15', 1, '東京事業所', 102, 'トラック2号', 1002, '佐藤次郎', 1, '2024-01-15 09:30:00', '2024-01-15 11:30:00', 1, '東京フェリー', 1, '東京港', '2便', 2, '大阪港', 1, '現金', 10000, 8000, 1, '大型車', 500, '東京-大阪'),
('2024011601', '2024-01-16', '2024-01-16', 1, '東京事業所', 101, 'トラック1号', 1001, '山田太郎', 1, '2024-01-16 08:00:00', '2024-01-16 10:00:00', 2, '関西フェリー', 3, '神戸港', '1便', 4, '高松港', 2, 'クレジット', 8000, 7000, 2, '中型車', 300, '神戸-高松');