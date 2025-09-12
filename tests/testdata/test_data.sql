-- Test data for dtako_rows
INSERT INTO dtako_rows (id, date, vehicle_no, driver_code, route_code, distance, fuel_amount) VALUES
('ROW001', '2024-01-15', 'V001', 'D001', 'R001', 150.5, 20.3),
('ROW002', '2024-01-15', 'V002', 'D002', 'R002', 200.8, 25.7),
('ROW003', '2024-01-16', 'V001', 'D001', 'R003', 175.2, 22.1);

-- Test data for dtako_events
INSERT INTO dtako_events (id, event_date, event_type, vehicle_no, driver_code, description, latitude, longitude) VALUES
('EVENT001', '2024-01-15 08:30:00', 'START', 'V001', 'D001', 'Trip started', 35.6762, 139.6503),
('EVENT002', '2024-01-15 12:15:00', 'STOP', 'V001', 'D001', 'Lunch break', 35.6895, 139.6917),
('EVENT003', '2024-01-15 16:45:00', 'END', 'V001', 'D001', 'Trip ended', 35.6762, 139.6503);

-- Test data for dtako_ferry
INSERT INTO dtako_ferry (id, date, route, vehicle_no, driver_code, departure_time, arrival_time, passengers, vehicles) VALUES
('FERRY001', '2024-01-15', 'ROUTE_A', 'F001', 'D001', '08:00:00', '10:00:00', 150, 25),
('FERRY002', '2024-01-15', 'ROUTE_B', 'F002', 'D002', '09:30:00', '11:30:00', 200, 30),
('FERRY003', '2024-01-16', 'ROUTE_A', 'F001', 'D001', '08:00:00', '10:00:00', 175, 28);