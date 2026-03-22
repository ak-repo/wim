-- +goose Up
-- Seed data for warehouse inventory system
-- Generated: production-quality test data

-- 1. WAREHOUSES
INSERT INTO warehouses (id, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at)
VALUES
    ('a1b2c3d4-1234-5678-9abc-def012345678', 'NYC-01', 'NYC Distribution Center', '450 Warehouse Ave', 'Building A', 'New York', 'NY', '10001', 'US', true, NOW(), NOW()),
    ('b2c3d4e5-2345-6789-abcd-ef0123456789', 'DXB-01', 'Dubai Logistics Hub', 'Al Quoz Industrial Area 3', 'Warehouse Complex Block B', 'Dubai', 'Dubai', '00000', 'AE', true, NOW(), NOW()),
    ('c3d4e5f6-3456-789a-bcde-f01234567890', 'LON-01', 'London Fulfillment Center', '45 Logistics Way', 'Unit 12', 'London', 'Greater London', 'SW1A 1AA', 'GB', true, NOW(), NOW()),
    ('d4e5f6a7-4567-89ab-cdef-012345678901', 'SIN-01', 'Singapore Regional Hub', '25 Changi North Crescent', 'Level 3', 'Singapore', 'Singapore', '498997', 'SG', true, NOW(), NOW()),
    ('e5f6a7b8-5678-9abc-def0-123456789012', 'AMS-01', 'Amsterdam Distribution', 'Schiphol Logistics Park', 'Hangar 45', 'Amsterdam', 'North Holland', '1118 CP', 'NL', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- 2. PRODUCTS
INSERT INTO products (id, sku, name, description, category, unit_of_measure, weight, length, width, height, barcode, is_active, created_at, updated_at)
VALUES
    ('f6a7b8c9-6789-abcd-ef01-234567890123', 'IPH-15PRO-256', 'iPhone 15 Pro 256GB', 'Apple iPhone 15 Pro with 256GB storage, Natural Titanium', 'Electronics', 'EA', 0.187, 0.147, 0.072, 0.008, '1234567890123', true, NOW(), NOW()),
    ('a7b8c9d0-789a-bcde-f012-345678901234', 'SAM-TV-65QLED', 'Samsung 65" QLED 4K Smart TV', 'Samsung 65-inch QLED 4K Smart TV with Quantum Processor', 'Electronics', 'EA', 22.500, 144.5, 83.0, 5.5, '9876543210987', true, NOW(), NOW()),
    ('b8c9d0e1-89ab-cdef-0123-456789012345', 'NES-ESP-ORIG', 'Nescafe Espresso Original', 'Nescafe Espresso Original Coffee Capsules, Pack of 30', 'FMCG', 'EA', 0.200, 0.15, 0.10, 0.12, '4567890123456', true, NOW(), NOW()),
    ('c9d0e1f2-9abc-def0-1234-567890123456', 'PG-TIDE-PODS', 'Tide Pods Laundry Detergent', 'Tide Pods Original Scent, 81 Count, HE Compatible', 'FMCG', 'EA', 2.100, 0.28, 0.18, 0.22, '7890123456789', true, NOW(), NOW()),
    ('d0e1f2a3-abcd-ef01-2345-678901234567', 'DEL-XPS-15', 'Dell XPS 15 Laptop', 'Dell XPS 15 9530, 15.6" FHD+, Intel Core i7, 16GB RAM, 512GB SSD', 'Electronics', 'EA', 1.860, 0.344, 0.230, 0.018, '3456789012345', true, NOW(), NOW())
ON CONFLICT (sku) DO NOTHING;

-- 3. LOCATIONS
INSERT INTO locations (id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at)
VALUES
    ('e1f2a3b4-bcde-f012-3456-789012345678', 'a1b2c3d4-1234-5678-9abc-def012345678', 'A', '01', 'R1', 'B1', 'A01-R1-B1', 'STORAGE', true, 500.00, true, NOW(), NOW()),
    ('f2a3b4c5-cdef-0123-4567-890123456789', 'a1b2c3d4-1234-5678-9abc-def012345678', 'A', '01', 'R1', 'B2', 'A01-R1-B2', 'STORAGE', false, 500.00, true, NOW(), NOW()),
    ('a3b4c5d6-def0-1234-5678-901234567890', 'b2c3d4e5-2345-6789-abcd-ef0123456789', 'B', '02', 'R2', 'B1', 'B02-R2-B1', 'STORAGE', true, 1000.00, true, NOW(), NOW()),
    ('b4c5d6e7-ef01-2345-6789-012345678901', 'b2c3d4e5-2345-6789-abcd-ef0123456789', 'B', '02', 'R2', 'B2', 'B02-R2-B2', 'BULK', false, 2000.00, true, NOW(), NOW()),
    ('c5d6e7f8-f012-3456-7890-123456789012', 'c3d4e5f6-3456-789a-bcde-f01234567890', 'C', '03', 'R3', 'B1', 'C03-R3-B1', 'STORAGE', true, 750.00, true, NOW(), NOW())
ON CONFLICT (warehouse_id, location_code) DO NOTHING;

-- 4. BATCHES
INSERT INTO batches (id, batch_number, product_id, supplier_id, manufacturing_date, expiry_date, origin_country, quantity_initial, quantity_remaining, is_active, created_at, updated_at)
VALUES
    ('d6e7f8a9-0123-4567-8901-234567890123', 'IP24-001-A', 'f6a7b8c9-6789-abcd-ef01-234567890123', '11111111-1111-1111-1111-111111111111', '2024-01-15', NULL, 'CN', 500, 500, true, NOW(), NOW()),
    ('e7f8a9b0-1234-5678-9012-345678901234', 'IP24-002-B', 'f6a7b8c9-6789-abcd-ef01-234567890123', '11111111-1111-1111-1111-111111111111', '2024-02-20', NULL, 'CN', 300, 300, true, NOW(), NOW()),
    ('f8a9b0c1-2345-6789-0123-456789012345', 'STV24-001', 'a7b8c9d0-789a-bcde-f012-345678901234', '22222222-2222-2222-2222-222222222222', '2024-01-10', NULL, 'MX', 200, 200, true, NOW(), NOW()),
    ('a9b0c1d2-3456-7890-1234-567890123456', 'NCF24-001', 'b8c9d0e1-89ab-cdef-0123-456789012345', '33333333-3333-3333-3333-333333333333', '2024-01-01', '2025-01-01', 'CH', 1000, 1000, true, NOW(), NOW()),
    ('b0c1d2e3-4567-8901-2345-678901234567', 'TID24-001', 'c9d0e1f2-9abc-def0-1234-567890123456', '44444444-4444-4444-4444-444444444444', '2024-02-01', '2026-02-01', 'US', 800, 800, true, NOW(), NOW())
ON CONFLICT (product_id, batch_number) DO NOTHING;

-- 5. INVENTORY
INSERT INTO inventory (id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_quantity, version, last_movement_id, created_at, updated_at)
VALUES
    ('c1d2e3f4-5678-9012-3456-789012345678', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'a1b2c3d4-1234-5678-9abc-def012345678', 'e1f2a3b4-bcde-f012-3456-789012345678', 'd6e7f8a9-0123-4567-8901-234567890123', 250, 50, 1, NULL, NOW(), NOW()),
    ('d2e3f4a5-6789-0123-4567-890123456789', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'a1b2c3d4-1234-5678-9abc-def012345678', 'f2a3b4c5-cdef-0123-4567-890123456789', 'e7f8a9b0-1234-5678-9012-345678901234', 150, 20, 1, NULL, NOW(), NOW()),
    ('e3f4a5b6-7890-1234-5678-901234567890', 'a7b8c9d0-789a-bcde-f012-345678901234', 'b2c3d4e5-2345-6789-abcd-ef0123456789', 'a3b4c5d6-def0-1234-5678-901234567890', 'f8a9b0c1-2345-6789-0123-456789012345', 100, 25, 1, NULL, NOW(), NOW()),
    ('f4a5b6c7-8901-2345-6789-012345678901', 'b8c9d0e1-89ab-cdef-0123-456789012345', 'b2c3d4e5-2345-6789-abcd-ef0123456789', 'b4c5d6e7-ef01-2345-6789-012345678901', 'a9b0c1d2-3456-7890-1234-567890123456', 500, 100, 1, NULL, NOW(), NOW()),
    ('a5b6c7d8-9012-3456-7890-123456789012', 'c9d0e1f2-9abc-def0-1234-567890123456', 'c3d4e5f6-3456-789a-bcde-f01234567890', 'c5d6e7f8-f012-3456-7890-123456789012', 'b0c1d2e3-4567-8901-2345-678901234567', 300, 75, 1, NULL, NOW(), NOW())
ON CONFLICT (product_id, warehouse_id, location_id, batch_id) DO NOTHING;

-- 6. STOCK_MOVEMENTS
INSERT INTO stock_movements (id, movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id, quantity, reference_type, reference_id, performed_by, notes, created_at)
VALUES
    ('b6c7d8e9-0123-4567-8901-234567890123', 'RECEIPT', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'a1b2c3d4-1234-5678-9abc-def012345678', NULL, 'e1f2a3b4-bcde-f012-3456-789012345678', 'd6e7f8a9-0123-4567-8901-234567890123', 500, 'PURCHASE_ORDER', '55555555-5555-5555-5555-555555555555', '66666666-6666-6666-6666-666666666666', 'Initial receipt - iPhone 15 Pro batch IP24-001-A', NOW() - INTERVAL '30 days'),
    ('c7d8e9f0-1234-5678-9012-345678901234', 'RECEIPT', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'a1b2c3d4-1234-5678-9abc-def012345678', NULL, 'f2a3b4c5-cdef-0123-4567-890123456789', 'e7f8a9b0-1234-5678-9012-345678901234', 300, 'PURCHASE_ORDER', '55555555-5555-5555-5555-555555555555', '66666666-6666-6666-6666-666666666666', 'Initial receipt - iPhone 15 Pro batch IP24-002-B', NOW() - INTERVAL '15 days'),
    ('d8e9f0a1-2345-6789-0123-456789012345', 'TRANSFER', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'a1b2c3d4-1234-5678-9abc-def012345678', 'e1f2a3b4-bcde-f012-3456-789012345678', 'f2a3b4c5-cdef-0123-4567-890123456789', 'd6e7f8a9-0123-4567-8901-234567890123', 100, 'ADJUSTMENT', NULL, '66666666-6666-6666-6666-666666666666', 'Stock rebalancing', NOW() - INTERVAL '10 days'),
    ('e9f0a1b2-3456-7890-1234-567890123456', 'SHIPMENT', 'a7b8c9d0-789a-bcde-f012-345678901234', 'b2c3d4e5-2345-6789-abcd-ef0123456789', 'a3b4c5d6-def0-1234-5678-901234567890', NULL, 'f8a9b0c1-2345-6789-0123-456789012345', 100, 'SALES_ORDER', '77777777-7777-7777-7777-777777777777', '66666666-6666-6666-6666-666666666666', 'Customer shipment - Samsung TVs', NOW() - INTERVAL '5 days')
ON CONFLICT DO NOTHING;

-- 7. PURCHASE_ORDERS
INSERT INTO purchase_orders (id, po_number, supplier_id, warehouse_id, order_date, expected_date, received_date, status, total_amount, notes, created_by, created_at, updated_at)
VALUES
    ('55555555-5555-5555-5555-555555555555', 'PO-2024-0001', '11111111-1111-1111-1111-111111111111', 'a1b2c3d4-1234-5678-9abc-def012345678', '2024-01-01', '2024-01-15', '2024-01-14', 'COMPLETED', 249500.00, 'iPhone 15 Pro initial stock order', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '45 days', NOW()),
    ('88888888-8888-8888-8888-888888888888', 'PO-2024-0002', '22222222-2222-2222-2222-222222222222', 'b2c3d4e5-2345-6789-abcd-ef0123456789', '2024-01-05', '2024-01-20', '2024-01-18', 'COMPLETED', 179800.00, 'Samsung TV order for Dubai market', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '40 days', NOW()),
    ('99999999-9999-9999-9999-999999999999', 'PO-2024-0003', '33333333-3333-3333-3333-333333333333', 'c3d4e5f6-3456-789a-bcde-f01234567890', '2024-03-01', '2024-03-15', NULL, 'PENDING', 25000.00, 'Nescafe capsules pending delivery', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '5 days', NOW()),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'PO-2024-0004', '44444444-4444-4444-4444-444444444444', 'd4e5f6a7-4567-89ab-cdef-012345678901', '2024-03-10', '2024-03-25', NULL, 'PENDING', 32000.00, 'Tide Pods order for Singapore', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '2 days', NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'PO-2024-0005', '55555555-5555-5555-5555-555555555555', 'e5f6a7b8-5678-9abc-def0-123456789012', '2024-03-15', '2024-04-01', NULL, 'PENDING', 95000.00, 'Dell XPS 15 laptops for Amsterdam', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '1 day', NOW())
ON CONFLICT (po_number) DO NOTHING;

-- 8. PURCHASE_ORDER_ITEMS
INSERT INTO purchase_order_items (id, purchase_order_id, product_id, batch_number, quantity_ordered, quantity_received, unit_price, created_at, updated_at)
VALUES
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '55555555-5555-5555-5555-555555555555', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'IP24-001-A', 500, 500, 499.00, NOW() - INTERVAL '45 days', NOW()),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', '55555555-5555-5555-5555-555555555555', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'IP24-002-B', 300, 300, 499.00, NOW() - INTERVAL '45 days', NOW()),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '88888888-8888-8888-8888-888888888888', 'a7b8c9d0-789a-bcde-f012-345678901234', 'STV24-001', 200, 200, 899.00, NOW() - INTERVAL '40 days', NOW()),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', '99999999-9999-9999-9999-999999999999', 'b8c9d0e1-89ab-cdef-0123-456789012345', NULL, 1000, 0, 25.00, NOW() - INTERVAL '5 days', NOW()),
    ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'c9d0e1f2-9abc-def0-1234-567890123456', NULL, 800, 0, 40.00, NOW() - INTERVAL '2 days', NOW())
ON CONFLICT DO NOTHING;

-- 9. SALES_ORDERS
INSERT INTO sales_orders (id, order_number, customer_id, warehouse_id, order_date, required_date, shipped_date, status, allocation_status, shipping_method, shipping_address, billing_address, subtotal, tax_amount, shipping_amount, total_amount, notes, created_by, created_at, updated_at)
VALUES
    ('77777777-7777-7777-7777-777777777777', 'SO-2024-0001', '22222222-2222-2222-2222-222222222222', 'b2c3d4e5-2345-6789-abcd-ef0123456789', NOW() - INTERVAL '7 days', '2024-03-25', NOW() - INTERVAL '2 days', 'COMPLETED', 'ALLOCATED', 'EXPRESS', 'Dubai Mall, Downtown Dubai, Dubai, UAE', 'Dubai Mall, Downtown Dubai, Dubai, UAE', 89900.00, 4500.00, 500.00, 94900.00, 'Wholesale order for retailer', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '7 days', NOW()),
    ('33333333-3333-3333-3333-333333333333', 'SO-2024-0002', '44444444-4444-4444-4444-444444444444', 'a1b2c3d4-1234-5678-9abc-def012345678', NOW() - INTERVAL '3 days', '2024-03-30', NULL, 'PENDING', 'UNALLOCATED', 'STANDARD', '123 Broadway, New York, NY 10001', '123 Broadway, New York, NY 10001', 124750.00, 10000.00, 250.00, 135000.00, 'Corporate iPhone order', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '3 days', NOW()),
    ('44444444-4444-4444-4444-444444444444', 'SO-2024-0003', '55555555-5555-5555-5555-555555555555', 'c3d4e5f6-3456-789a-bcde-f01234567890', NOW() - INTERVAL '1 day', '2024-04-05', NULL, 'PENDING', 'PARTIALLY_ALLOCATED', 'STANDARD', '10 Downing Street, London SW1A 2AA', '10 Downing Street, London SW1A 2AA', 12000.00, 2400.00, 150.00, 14550.00, 'Government office supply', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '1 day', NOW()),
    ('12121212-1212-1212-1212-121212121212', 'SO-2024-0004', '66666666-6666-6666-6666-666666666666', 'e5f6a7b8-5678-9abc-def0-123456789012', NOW() - INTERVAL '12 hours', '2024-03-28', NULL, 'PENDING', 'UNALLOCATED', 'EXPRESS', 'Herengracht 386, 1016 CJ Amsterdam', 'Herengracht 386, 1016 CJ Amsterdam', 15000.00, 3000.00, 200.00, 18200.00, 'Tech startup laptop order', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '12 hours', NOW()),
    ('13131313-1313-1313-1313-131313131313', 'SO-2024-0005', '77777777-7777-7777-7777-777777777777', 'd4e5f6a7-4567-89ab-cdef-012345678901', NOW() - INTERVAL '2 days', '2024-04-01', NULL, 'PROCESSING', 'ALLOCATED', 'STANDARD', '50 Raffles Place, Singapore 048623', '50 Raffles Place, Singapore 048623', 5000.00, 400.00, 100.00, 5500.00, 'Coffee shop supply', '66666666-6666-6666-6666-666666666666', NOW() - INTERVAL '2 days', NOW())
ON CONFLICT (order_number) DO NOTHING;

-- 10. SALES_ORDER_ITEMS
INSERT INTO sales_order_items (id, sales_order_id, product_id, batch_id, location_id, quantity_ordered, quantity_allocated, quantity_picked, quantity_shipped, unit_price, created_at, updated_at)
VALUES
    ('14141414-1414-1414-1414-141414141414', '77777777-7777-7777-7777-777777777777', 'a7b8c9d0-789a-bcde-f012-345678901234', 'f8a9b0c1-2345-6789-0123-456789012345', 'a3b4c5d6-def0-1234-5678-901234567890', 100, 100, 100, 100, 899.00, NOW() - INTERVAL '7 days', NOW()),
    ('15151515-1515-1515-1515-151515151515', '33333333-3333-3333-3333-333333333333', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'd6e7f8a9-0123-4567-8901-234567890123', 'e1f2a3b4-bcde-f012-3456-789012345678', 250, 0, 0, 0, 499.00, NOW() - INTERVAL '3 days', NOW()),
    ('16161616-1616-1616-1616-161616161616', '44444444-4444-4444-4444-444444444444', 'c9d0e1f2-9abc-def0-1234-567890123456', 'b0c1d2e3-4567-8901-2345-678901234567', 'c5d6e7f8-f012-3456-7890-123456789012', 300, 150, 0, 0, 40.00, NOW() - INTERVAL '1 day', NOW()),
    ('17171717-1717-1717-1717-171717171717', '12121212-1212-1212-1212-121212121212', 'd0e1f2a3-abcd-ef01-2345-678901234567', NULL, NULL, 10, 0, 0, 0, 1500.00, NOW() - INTERVAL '12 hours', NOW()),
    ('18181818-1818-1818-1818-181818181818', '13131313-1313-1313-1313-131313131313', 'b8c9d0e1-89ab-cdef-0123-456789012345', 'a9b0c1d2-3456-7890-1234-567890123456', 'b4c5d6e7-ef01-2345-6789-012345678901', 200, 200, 0, 0, 25.00, NOW() - INTERVAL '2 days', NOW())
ON CONFLICT DO NOTHING;

-- 11. BARCODES
INSERT INTO barcodes (id, product_id, barcode_value, barcode_type, is_primary, created_at)
VALUES
    ('19191919-1919-1919-1919-191919191919', 'f6a7b8c9-6789-abcd-ef01-234567890123', '1234567890123', 'UPC', true, NOW()),
    ('20202020-2020-2020-2020-202020202020', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'BATCH-IP24-001-A', 'CODE128', false, NOW()),
    ('21212121-2121-2121-2121-212121212121', 'a7b8c9d0-789a-bcde-f012-345678901234', '9876543210987', 'UPC', true, NOW()),
    ('22222222-2222-2222-2222-222222222222', 'b8c9d0e1-89ab-cdef-0123-456789012345', '4567890123456', 'EAN', true, NOW()),
    ('23232323-2323-2323-2323-232323232323', 'c9d0e1f2-9abc-def0-1234-567890123456', '7890123456789', 'UPC', true, NOW())
ON CONFLICT (barcode_value) DO NOTHING;

-- 12. AUDIT_LOGS
INSERT INTO audit_logs (id, entity_type, entity_id, action, user_id, old_values, new_values, ip_address, user_agent, created_at)
VALUES
    ('24242424-2424-2424-2424-242424242424', 'WAREHOUSE', 'a1b2c3d4-1234-5678-9abc-def012345678', 'CREATE', '66666666-6666-6666-6666-666666666666', NULL, '{"code": "NYC-01", "name": "NYC Distribution Center"}'::jsonb, '192.168.1.100', 'Mozilla/5.0', NOW() - INTERVAL '60 days'),
    ('25252525-2525-2525-2525-252525252525', 'PRODUCT', 'f6a7b8c9-6789-abcd-ef01-234567890123', 'CREATE', '66666666-6666-6666-6666-666666666666', NULL, '{"sku": "IPH-15PRO-256", "name": "iPhone 15 Pro 256GB"}'::jsonb, '192.168.1.100', 'Mozilla/5.0', NOW() - INTERVAL '60 days'),
    ('26262626-2626-2626-2626-262626262626', 'PURCHASE_ORDER', '55555555-5555-5555-5555-555555555555', 'UPDATE', '66666666-6666-6666-6666-666666666666', '{"status": "PENDING"}'::jsonb, '{"status": "COMPLETED", "received_date": "2024-01-14"}'::jsonb, '192.168.1.100', 'Mozilla/5.0', NOW() - INTERVAL '30 days'),
    ('27272727-2727-2727-2727-272727272727', 'INVENTORY', 'c1d2e3f4-5678-9012-3456-789012345678', 'UPDATE', '66666666-6666-6666-6666-666666666666', '{"quantity": 500, "reserved_quantity": 0}'::jsonb, '{"quantity": 250, "reserved_quantity": 50}'::jsonb, '192.168.1.100', 'Mozilla/5.0', NOW() - INTERVAL '5 days'),
    ('28282828-2828-2828-2828-282828282828', 'SALES_ORDER', '77777777-7777-7777-7777-777777777777', 'UPDATE', '66666666-6666-6666-6666-666666666666', '{"status": "PENDING"}'::jsonb, '{"status": "COMPLETED", "shipped_date": "2024-03-18"}'::jsonb, '192.168.1.100', 'Mozilla/5.0', NOW() - INTERVAL '2 days')
ON CONFLICT DO NOTHING;

-- +goose Down
-- Delete seeded data in reverse order (respecting foreign key dependencies)
DELETE FROM audit_logs;
DELETE FROM barcodes;
DELETE FROM sales_order_items;
DELETE FROM sales_orders;
DELETE FROM purchase_order_items;
DELETE FROM purchase_orders;
DELETE FROM stock_movements;
DELETE FROM inventory;
DELETE FROM batches;
DELETE FROM locations;
DELETE FROM warehouses;
DELETE FROM products;
