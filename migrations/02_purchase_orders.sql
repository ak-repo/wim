-- +goose Up
CREATE TABLE IF NOT EXISTS purchase_orders (
	id BIGSERIAL PRIMARY KEY,
	ref_code TEXT UNIQUE NOT NULL,
	supplier_id BIGINT NOT NULL,
	warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
	status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
	expected_date TIMESTAMP WITH TIME ZONE,
	received_date TIMESTAMP WITH TIME ZONE,
	notes TEXT,
	created_by BIGINT REFERENCES users(id),
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS purchase_order_items (
	id BIGSERIAL PRIMARY KEY,
	purchase_order_id BIGINT NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
	product_id BIGINT NOT NULL REFERENCES products(id),
	quantity_ordered INT NOT NULL,
	quantity_received INT NOT NULL DEFAULT 0,
	batch_number TEXT,
	unit_price DECIMAL(12, 4),
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_purchase_orders_ref_code ON purchase_orders(ref_code);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_supplier ON purchase_orders(supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_warehouse ON purchase_orders(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON purchase_orders(status);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_order ON purchase_order_items(purchase_order_id);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_product ON purchase_order_items(product_id);

-- +goose Down
DROP INDEX IF EXISTS idx_purchase_order_items_product;
DROP INDEX IF EXISTS idx_purchase_order_items_order;
DROP INDEX IF EXISTS idx_purchase_orders_status;
DROP INDEX IF EXISTS idx_purchase_orders_warehouse;
DROP INDEX IF EXISTS idx_purchase_orders_supplier;
DROP INDEX IF EXISTS idx_purchase_orders_ref_code;
DROP TABLE IF EXISTS purchase_order_items;
DROP TABLE IF EXISTS purchase_orders;
