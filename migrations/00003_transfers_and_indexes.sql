-- +goose Up
CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transfer_number VARCHAR(50) UNIQUE NOT NULL,
    source_warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    dest_warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    requested_by UUID,
    approved_by UUID,
    shipped_date TIMESTAMP WITH TIME ZONE,
    received_date TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CHECK (source_warehouse_id <> dest_warehouse_id)
);

CREATE TABLE IF NOT EXISTS transfer_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transfer_id UUID NOT NULL REFERENCES transfers(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    batch_id UUID REFERENCES batches(id),
    quantity_requested INTEGER NOT NULL,
    quantity_shipped INTEGER NOT NULL DEFAULT 0,
    quantity_received INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CHECK (quantity_requested >= 0),
    CHECK (quantity_shipped >= 0),
    CHECK (quantity_received >= 0)
);

CREATE INDEX IF NOT EXISTS idx_transfers_source_warehouse ON transfers(source_warehouse_id);
CREATE INDEX IF NOT EXISTS idx_transfers_dest_warehouse ON transfers(dest_warehouse_id);
CREATE INDEX IF NOT EXISTS idx_transfers_status_created ON transfers(status, created_at);
CREATE INDEX IF NOT EXISTS idx_transfer_items_transfer_id ON transfer_items(transfer_id);

CREATE INDEX IF NOT EXISTS idx_batches_expiry_date ON batches(expiry_date);
CREATE INDEX IF NOT EXISTS idx_sales_orders_status_order_date ON sales_orders(status, order_date);
CREATE INDEX IF NOT EXISTS idx_inventory_product_warehouse ON inventory(product_id, warehouse_id);

-- +goose Down
DROP INDEX IF EXISTS idx_inventory_product_warehouse;
DROP INDEX IF EXISTS idx_sales_orders_status_order_date;
DROP INDEX IF EXISTS idx_batches_expiry_date;

DROP INDEX IF EXISTS idx_transfer_items_transfer_id;
DROP INDEX IF EXISTS idx_transfers_status_created;
DROP INDEX IF EXISTS idx_transfers_dest_warehouse;
DROP INDEX IF EXISTS idx_transfers_source_warehouse;

DROP TABLE IF EXISTS transfer_items CASCADE;
DROP TABLE IF EXISTS transfers CASCADE;
