-- +goose Up
CREATE TABLE IF NOT EXISTS picking_tasks (
    id BIGSERIAL PRIMARY KEY,
    ref_code TEXT UNIQUE NOT NULL,
    sales_order_id BIGINT NOT NULL REFERENCES sales_orders(id),
    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    priority VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    assigned_to BIGINT REFERENCES users(id),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS picking_task_items (
    id BIGSERIAL PRIMARY KEY,
    picking_task_id BIGINT NOT NULL REFERENCES picking_tasks(id) ON DELETE CASCADE,
    sales_order_item_id BIGINT NOT NULL REFERENCES sales_order_items(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    location_id BIGINT REFERENCES locations(id),
    batch_id BIGINT,
    quantity_required INT NOT NULL,
    quantity_picked INT NOT NULL DEFAULT 0,
    picked_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_picking_tasks_sales_order ON picking_tasks(sales_order_id);
CREATE INDEX IF NOT EXISTS idx_picking_tasks_warehouse ON picking_tasks(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_picking_tasks_status ON picking_tasks(status);
CREATE INDEX IF NOT EXISTS idx_picking_tasks_assigned_to ON picking_tasks(assigned_to);
CREATE INDEX IF NOT EXISTS idx_picking_tasks_priority ON picking_tasks(priority);
CREATE INDEX IF NOT EXISTS idx_picking_tasks_created ON picking_tasks(created_at);
CREATE INDEX IF NOT EXISTS idx_picking_task_items_task ON picking_task_items(picking_task_id);
CREATE INDEX IF NOT EXISTS idx_picking_task_items_product ON picking_task_items(product_id);
CREATE INDEX IF NOT EXISTS idx_picking_task_items_location ON picking_task_items(location_id);
CREATE INDEX IF NOT EXISTS idx_picking_task_items_status ON picking_task_items(status);

-- +goose Down
DROP INDEX IF EXISTS idx_picking_task_items_status;
DROP INDEX IF EXISTS idx_picking_task_items_location;
DROP INDEX IF EXISTS idx_picking_task_items_product;
DROP INDEX IF EXISTS idx_picking_task_items_task;
DROP TABLE IF EXISTS picking_task_items;

DROP INDEX IF EXISTS idx_picking_tasks_created;
DROP INDEX IF EXISTS idx_picking_tasks_priority;
DROP INDEX IF EXISTS idx_picking_tasks_assigned_to;
DROP INDEX IF EXISTS idx_picking_tasks_status;
DROP INDEX IF EXISTS idx_picking_tasks_warehouse;
DROP INDEX IF EXISTS idx_picking_tasks_sales_order;
DROP TABLE IF EXISTS picking_tasks;
