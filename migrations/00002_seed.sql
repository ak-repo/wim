-- +goose Up
-- Seed data for warehouse inventory system

-- 1. WAREHOUSES
INSERT INTO warehouses (ref_code, code, name, address_line1, address_line2, city, state, postal_code, country, is_active)
VALUES
    ('WH-NYC-01', 'NYC-01', 'NYC Distribution Center', '450 Warehouse Ave', 'Building A', 'New York', 'NY', '10001', 'US', true),
    ('WH-DXB-01', 'DXB-01', 'Dubai Logistics Hub', 'Al Quoz Industrial Area 3', 'Warehouse Complex Block B', 'Dubai', 'Dubai', '00000', 'AE', true),
    ('WH-LON-01', 'LON-01', 'London Fulfillment Center', '45 Logistics Way', 'Unit 12', 'London', 'Greater London', 'SW1A 1AA', 'GB', true),
    ('WH-SIN-01', 'SIN-01', 'Singapore Regional Hub', '25 Changi North Crescent', 'Level 3', 'Singapore', 'Singapore', '498997', 'SG', true),
    ('WH-AMS-01', 'AMS-01', 'Amsterdam Distribution', 'Schiphol Logistics Park', 'Hangar 45', 'Amsterdam', 'North Holland', '1118 CP', 'NL', true)
ON CONFLICT DO NOTHING;

-- 2. PRODUCTS
INSERT INTO products (ref_code, sku, name, description, category, unit_of_measure, weight, length, width, height, barcode, is_active)
VALUES
    ('PRD-IPH-15PRO-256', 'IPH-15PRO-256', 'iPhone 15 Pro 256GB', 'Apple iPhone 15 Pro with 256GB storage, Natural Titanium', 'Electronics', 'EA', 0.187, 0.147, 0.072, 0.008, '1234567890123', true),
    ('PRD-SAM-TV-65QLED', 'SAM-TV-65QLED', 'Samsung 65" QLED 4K Smart TV', 'Samsung 65-inch QLED 4K Smart TV with Quantum Processor', 'Electronics', 'EA', 22.500, 144.5, 83.0, 5.5, '9876543210987', true),
    ('PRD-NES-ESP-ORIG', 'NES-ESP-ORIG', 'Nescafe Espresso Original', 'Nescafe Espresso Original Coffee Capsules, Pack of 30', 'FMCG', 'EA', 0.200, 0.15, 0.10, 0.12, '4567890123456', true),
    ('PRD-PG-TIDE-PODS', 'PG-TIDE-PODS', 'Tide Pods Laundry Detergent', 'Tide Pods Original Scent, 81 Count, HE Compatible', 'FMCG', 'EA', 2.100, 0.28, 0.18, 0.22, '7890123456789', true),
    ('PRD-DEL-XPS-15', 'DEL-XPS-15', 'Dell XPS 15 Laptop', 'Dell XPS 15 9530, 15.6" FHD+, Intel Core i7, 16GB RAM, 512GB SSD', 'Electronics', 'EA', 1.860, 0.344, 0.230, 0.018, '3456789012345', true)
ON CONFLICT DO NOTHING;

-- 3. LOCATIONS
INSERT INTO locations (ref_code, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active)
VALUES
    ('LOC-NYC-01-A01-R1-B1', (SELECT id FROM warehouses WHERE code = 'NYC-01'), 'A', '01', 'R1', 'B1', 'A01-R1-B1', 'STORAGE', true, 500.00, true),
    ('LOC-NYC-01-A01-R1-B2', (SELECT id FROM warehouses WHERE code = 'NYC-01'), 'A', '01', 'R1', 'B2', 'A01-R1-B2', 'STORAGE', false, 500.00, true),
    ('LOC-DXB-01-B02-R2-B1', (SELECT id FROM warehouses WHERE code = 'DXB-01'), 'B', '02', 'R2', 'B1', 'B02-R2-B1', 'STORAGE', true, 1000.00, true),
    ('LOC-DXB-01-B02-R2-B2', (SELECT id FROM warehouses WHERE code = 'DXB-01'), 'B', '02', 'R2', 'B2', 'B02-R2-B2', 'BULK', false, 2000.00, true),
    ('LOC-LON-01-C03-R3-B1', (SELECT id FROM warehouses WHERE code = 'LON-01'), 'C', '03', 'R3', 'B1', 'C03-R3-B1', 'STORAGE', true, 750.00, true)
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM locations;
DELETE FROM products;
DELETE FROM warehouses;
