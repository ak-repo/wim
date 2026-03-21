package postgres

import (
	"context"
	"fmt"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type SalesOrderRepository interface {
	Create(ctx context.Context, order *domain.SalesOrder) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error)
	GetByNumber(ctx context.Context, number string) (*domain.SalesOrder, error)
	List(ctx context.Context, filter SalesOrderFilter) ([]*domain.SalesOrder, error)
	Update(ctx context.Context, order *domain.SalesOrder) error

	CreateItem(ctx context.Context, item *domain.SalesOrderItem) error
	GetItems(ctx context.Context, orderID uuid.UUID) ([]*domain.SalesOrderItem, error)
	UpdateItem(ctx context.Context, item *domain.SalesOrderItem) error
}

type SalesOrderFilter struct {
	CustomerID       *uuid.UUID
	WarehouseID      *uuid.UUID
	Status           string
	AllocationStatus string
	FromDate         *string
	ToDate           *string
	Limit            int
	Offset           int
}

type salesOrderRepo struct {
	db *DB
}

func NewSalesOrderRepository(db *DB) SalesOrderRepository {
	return &salesOrderRepo{db: db}
}

func (r *salesOrderRepo) Create(ctx context.Context, order *domain.SalesOrder) error {
	query := `
		INSERT INTO sales_orders (id, order_number, customer_id, warehouse_id, order_date, required_date, shipped_date, status, allocation_status, shipping_method, shipping_address, billing_address, subtotal, tax_amount, shipping_amount, total_amount, notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	_, err := r.db.Pool.Exec(ctx, query,
		order.ID, order.OrderNumber, order.CustomerID, order.WarehouseID, order.OrderDate,
		order.RequiredDate, order.ShippedDate, order.Status, order.AllocationStatus,
		order.ShippingMethod, order.ShippingAddress, order.BillingAddress, order.Subtotal,
		order.TaxAmount, order.ShippingAmount, order.TotalAmount, order.Notes, order.CreatedBy,
		order.CreatedAt, order.UpdatedAt,
	)
	return err
}

func (r *salesOrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	query := `
		SELECT id, order_number, customer_id, warehouse_id, order_date, required_date, shipped_date, status, allocation_status, shipping_method, shipping_address, billing_address, subtotal, tax_amount, shipping_amount, total_amount, notes, created_by, created_at, updated_at
		FROM sales_orders WHERE id = $1`

	var o domain.SalesOrder
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.OrderNumber, &o.CustomerID, &o.WarehouseID, &o.OrderDate,
		&o.RequiredDate, &o.ShippedDate, &o.Status, &o.AllocationStatus,
		&o.ShippingMethod, &o.ShippingAddress, &o.BillingAddress, &o.Subtotal,
		&o.TaxAmount, &o.ShippingAmount, &o.TotalAmount, &o.Notes, &o.CreatedBy,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *salesOrderRepo) GetByNumber(ctx context.Context, number string) (*domain.SalesOrder, error) {
	query := `
		SELECT id, order_number, customer_id, warehouse_id, order_date, required_date, shipped_date, status, allocation_status, shipping_method, shipping_address, billing_address, subtotal, tax_amount, shipping_amount, total_amount, notes, created_by, created_at, updated_at
		FROM sales_orders WHERE order_number = $1`

	var o domain.SalesOrder
	err := r.db.Pool.QueryRow(ctx, query, number).Scan(
		&o.ID, &o.OrderNumber, &o.CustomerID, &o.WarehouseID, &o.OrderDate,
		&o.RequiredDate, &o.ShippedDate, &o.Status, &o.AllocationStatus,
		&o.ShippingMethod, &o.ShippingAddress, &o.BillingAddress, &o.Subtotal,
		&o.TaxAmount, &o.ShippingAmount, &o.TotalAmount, &o.Notes, &o.CreatedBy,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *salesOrderRepo) List(ctx context.Context, filter SalesOrderFilter) ([]*domain.SalesOrder, error) {
	query := `
		SELECT id, order_number, customer_id, warehouse_id, order_date, required_date, shipped_date, status, allocation_status, shipping_method, shipping_address, billing_address, subtotal, tax_amount, shipping_amount, total_amount, notes, created_by, created_at, updated_at
		FROM sales_orders WHERE 1=1`

	args := []interface{}{}
	argNum := 1

	if filter.CustomerID != nil {
		query += fmt.Sprintf(" AND customer_id = $%d", argNum)
		args = append(args, *filter.CustomerID)
		argNum++
	}

	if filter.WarehouseID != nil {
		query += fmt.Sprintf(" AND warehouse_id = $%d", argNum)
		args = append(args, *filter.WarehouseID)
		argNum++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, filter.Status)
		argNum++
	}

	if filter.AllocationStatus != "" {
		query += fmt.Sprintf(" AND allocation_status = $%d", argNum)
		args = append(args, filter.AllocationStatus)
		argNum++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argNum)
		args = append(args, filter.Limit)
		argNum++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argNum)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.SalesOrder
	for rows.Next() {
		var o domain.SalesOrder
		err := rows.Scan(
			&o.ID, &o.OrderNumber, &o.CustomerID, &o.WarehouseID, &o.OrderDate,
			&o.RequiredDate, &o.ShippedDate, &o.Status, &o.AllocationStatus,
			&o.ShippingMethod, &o.ShippingAddress, &o.BillingAddress, &o.Subtotal,
			&o.TaxAmount, &o.ShippingAmount, &o.TotalAmount, &o.Notes, &o.CreatedBy,
			&o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

func (r *salesOrderRepo) Update(ctx context.Context, order *domain.SalesOrder) error {
	query := `
		UPDATE sales_orders SET 
			order_number = $2, customer_id = $3, warehouse_id = $4, order_date = $5,
			required_date = $6, shipped_date = $7, status = $8, allocation_status = $9,
			shipping_method = $10, shipping_address = $11, billing_address = $12,
			subtotal = $13, tax_amount = $14, shipping_amount = $15, total_amount = $16,
			notes = $17, updated_at = $18
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		order.ID, order.OrderNumber, order.CustomerID, order.WarehouseID, order.OrderDate,
		order.RequiredDate, order.ShippedDate, order.Status, order.AllocationStatus,
		order.ShippingMethod, order.ShippingAddress, order.BillingAddress, order.Subtotal,
		order.TaxAmount, order.ShippingAmount, order.TotalAmount, order.Notes, order.UpdatedAt,
	)
	return err
}

func (r *salesOrderRepo) CreateItem(ctx context.Context, item *domain.SalesOrderItem) error {
	query := `
		INSERT INTO sales_order_items (id, sales_order_id, product_id, batch_id, location_id, quantity_ordered, quantity_allocated, quantity_picked, quantity_shipped, unit_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.SalesOrderID, item.ProductID, item.BatchID, item.LocationID,
		item.QuantityOrdered, item.QuantityAllocated, item.QuantityPicked, item.QuantityShipped,
		item.UnitPrice, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *salesOrderRepo) GetItems(ctx context.Context, orderID uuid.UUID) ([]*domain.SalesOrderItem, error) {
	query := `
		SELECT id, sales_order_id, product_id, batch_id, location_id, quantity_ordered, quantity_allocated, quantity_picked, quantity_shipped, unit_price, created_at, updated_at
		FROM sales_order_items WHERE sales_order_id = $1`

	rows, err := r.db.Pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.SalesOrderItem
	for rows.Next() {
		var item domain.SalesOrderItem
		err := rows.Scan(
			&item.ID, &item.SalesOrderID, &item.ProductID, &item.BatchID, &item.LocationID,
			&item.QuantityOrdered, &item.QuantityAllocated, &item.QuantityPicked, &item.QuantityShipped,
			&item.UnitPrice, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (r *salesOrderRepo) UpdateItem(ctx context.Context, item *domain.SalesOrderItem) error {
	query := `
		UPDATE sales_order_items SET 
			product_id = $2, batch_id = $3, location_id = $4, quantity_ordered = $5,
			quantity_allocated = $6, quantity_picked = $7, quantity_shipped = $8,
			unit_price = $9, updated_at = $10
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.ProductID, item.BatchID, item.LocationID, item.QuantityOrdered,
		item.QuantityAllocated, item.QuantityPicked, item.QuantityShipped, item.UnitPrice,
		item.UpdatedAt,
	)
	return err
}
