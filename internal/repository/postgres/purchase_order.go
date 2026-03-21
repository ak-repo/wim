package postgres

import (
	"context"
	"fmt"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type PurchaseOrderRepository interface {
	Create(ctx context.Context, po *domain.PurchaseOrder) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseOrder, error)
	GetByNumber(ctx context.Context, number string) (*domain.PurchaseOrder, error)
	List(ctx context.Context, filter PurchaseOrderFilter) ([]*domain.PurchaseOrder, error)
	Update(ctx context.Context, po *domain.PurchaseOrder) error

	CreateItem(ctx context.Context, item *domain.PurchaseOrderItem) error
	GetItems(ctx context.Context, poID uuid.UUID) ([]*domain.PurchaseOrderItem, error)
	UpdateItem(ctx context.Context, item *domain.PurchaseOrderItem) error
}

type PurchaseOrderFilter struct {
	SupplierID  *uuid.UUID
	WarehouseID *uuid.UUID
	Status      string
	FromDate    *string
	ToDate      *string
	Limit       int
	Offset      int
}

type purchaseOrderRepo struct {
	db *DB
}

func NewPurchaseOrderRepository(db *DB) PurchaseOrderRepository {
	return &purchaseOrderRepo{db: db}
}

func (r *purchaseOrderRepo) Create(ctx context.Context, po *domain.PurchaseOrder) error {
	query := `
		INSERT INTO purchase_orders (id, po_number, supplier_id, warehouse_id, order_date, expected_date, status, total_amount, notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Pool.Exec(ctx, query,
		po.ID, po.PONumber, po.SupplierID, po.WarehouseID, po.OrderDate, po.ExpectedDate,
		po.Status, po.TotalAmount, po.Notes, po.CreatedBy, po.CreatedAt, po.UpdatedAt,
	)
	return err
}

func (r *purchaseOrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseOrder, error) {
	query := `
		SELECT id, po_number, supplier_id, warehouse_id, order_date, expected_date, received_date, status, total_amount, notes, created_by, created_at, updated_at
		FROM purchase_orders WHERE id = $1`

	var po domain.PurchaseOrder
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&po.ID, &po.PONumber, &po.SupplierID, &po.WarehouseID, &po.OrderDate, &po.ExpectedDate,
		&po.ReceivedDate, &po.Status, &po.TotalAmount, &po.Notes, &po.CreatedBy, &po.CreatedAt, &po.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *purchaseOrderRepo) GetByNumber(ctx context.Context, number string) (*domain.PurchaseOrder, error) {
	query := `
		SELECT id, po_number, supplier_id, warehouse_id, order_date, expected_date, received_date, status, total_amount, notes, created_by, created_at, updated_at
		FROM purchase_orders WHERE po_number = $1`

	var po domain.PurchaseOrder
	err := r.db.Pool.QueryRow(ctx, query, number).Scan(
		&po.ID, &po.PONumber, &po.SupplierID, &po.WarehouseID, &po.OrderDate, &po.ExpectedDate,
		&po.ReceivedDate, &po.Status, &po.TotalAmount, &po.Notes, &po.CreatedBy, &po.CreatedAt, &po.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *purchaseOrderRepo) List(ctx context.Context, filter PurchaseOrderFilter) ([]*domain.PurchaseOrder, error) {
	query := `
		SELECT id, po_number, supplier_id, warehouse_id, order_date, expected_date, received_date, status, total_amount, notes, created_by, created_at, updated_at
		FROM purchase_orders WHERE 1=1`

	args := []interface{}{}
	argNum := 1

	if filter.SupplierID != nil {
		query += fmt.Sprintf(" AND supplier_id = $%d", argNum)
		args = append(args, *filter.SupplierID)
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

	var orders []*domain.PurchaseOrder
	for rows.Next() {
		var po domain.PurchaseOrder
		err := rows.Scan(
			&po.ID, &po.PONumber, &po.SupplierID, &po.WarehouseID, &po.OrderDate, &po.ExpectedDate,
			&po.ReceivedDate, &po.Status, &po.TotalAmount, &po.Notes, &po.CreatedBy, &po.CreatedAt, &po.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &po)
	}
	return orders, nil
}

func (r *purchaseOrderRepo) Update(ctx context.Context, po *domain.PurchaseOrder) error {
	query := `
		UPDATE purchase_orders SET 
			po_number = $2, supplier_id = $3, warehouse_id = $4, order_date = $5,
			expected_date = $6, received_date = $7, status = $8, total_amount = $9,
			notes = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		po.ID, po.PONumber, po.SupplierID, po.WarehouseID, po.OrderDate, po.ExpectedDate,
		po.ReceivedDate, po.Status, po.TotalAmount, po.Notes, po.UpdatedAt,
	)
	return err
}

func (r *purchaseOrderRepo) CreateItem(ctx context.Context, item *domain.PurchaseOrderItem) error {
	query := `
		INSERT INTO purchase_order_items (id, purchase_order_id, product_id, batch_number, quantity_ordered, quantity_received, unit_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.PurchaseOrderID, item.ProductID, item.BatchNumber, item.QuantityOrdered,
		item.QuantityReceived, item.UnitPrice, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *purchaseOrderRepo) GetItems(ctx context.Context, poID uuid.UUID) ([]*domain.PurchaseOrderItem, error) {
	query := `
		SELECT id, purchase_order_id, product_id, batch_number, quantity_ordered, quantity_received, unit_price, created_at, updated_at
		FROM purchase_order_items WHERE purchase_order_id = $1`

	rows, err := r.db.Pool.Query(ctx, query, poID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.PurchaseOrderItem
	for rows.Next() {
		var item domain.PurchaseOrderItem
		err := rows.Scan(
			&item.ID, &item.PurchaseOrderID, &item.ProductID, &item.BatchNumber, &item.QuantityOrdered,
			&item.QuantityReceived, &item.UnitPrice, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (r *purchaseOrderRepo) UpdateItem(ctx context.Context, item *domain.PurchaseOrderItem) error {
	query := `
		UPDATE purchase_order_items SET 
			product_id = $2, batch_number = $3, quantity_ordered = $4,
			quantity_received = $5, unit_price = $6, updated_at = $7
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.ProductID, item.BatchNumber, item.QuantityOrdered,
		item.QuantityReceived, item.UnitPrice, item.UpdatedAt,
	)
	return err
}
