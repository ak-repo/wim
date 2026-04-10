package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/model"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/jackc/pgx/v5"
)

var (
	ErrSalesOrderNotFound             = errors.New("sales order not found")
	ErrSalesOrderItemNotFound         = errors.New("sales order item not found")
	ErrInsufficientStockForAllocation = errors.New("insufficient stock for allocation")
)

type SalesOrderRepository interface {
	// Sales Order CRUD
	Create(ctx context.Context, order *model.SalesOrderDTO, items []*model.SalesOrderItemDTO) (*model.SalesOrderDTO, error)
	GetByID(ctx context.Context, orderID int) (*model.SalesOrderDTO, error)
	GetByRefCode(ctx context.Context, refCode string) (*model.SalesOrderDTO, error)
	List(ctx context.Context, params *model.SalesOrderParams) (model.SalesOrderDTOs, error)
	Count(ctx context.Context, params *model.SalesOrderParams) (int, error)
	Update(ctx context.Context, order *model.SalesOrderDTO) error
	UpdateStatus(ctx context.Context, orderID int, status string) error
	UpdateAllocationStatus(ctx context.Context, orderID int, allocationStatus string) error

	// Sales Order Items
	GetItemsByOrderID(ctx context.Context, orderID int) (model.SalesOrderItemDTOs, error)
	GetItemByID(ctx context.Context, itemID int) (*model.SalesOrderItemDTO, error)
	UpdateItem(ctx context.Context, item *model.SalesOrderItemDTO) error
	UpdateItemAllocation(ctx context.Context, itemID int, quantityReserved int, allocationStatus string, locationID *int, batchID *int) error
	UpdateItemShippedQty(ctx context.Context, itemID int, quantityShipped int) error

	// Allocation Operations (Transaction-based)
	AllocateStock(ctx context.Context, orderID int, performedBy *int) error
	DeallocateStock(ctx context.Context, orderID int) error

	// Shipping Operations
	ShipOrder(ctx context.Context, orderID int, shippedDate *time.Time, performedBy *int) error
	ShipOrderItem(ctx context.Context, itemID int, quantity int) error

	// Locking for concurrency control
	LockOrderForUpdate(ctx context.Context, orderID int) (*model.SalesOrderDTO, error)
	LockItemForUpdate(ctx context.Context, itemID int) (*model.SalesOrderItemDTO, error)
}

type salesOrderRepository struct {
	db *db.DB
}

func NewSalesOrderRepository(database *db.DB) SalesOrderRepository {
	return &salesOrderRepository{db: database}
}

func (r *salesOrderRepository) Create(ctx context.Context, order *model.SalesOrderDTO, items []*model.SalesOrderItemDTO) (*model.SalesOrderDTO, error) {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start transaction")
	}
	defer tx.Rollback(ctx)

	// Insert sales order
	orderQuery := `
		INSERT INTO sales_orders (ref_code, customer_id, warehouse_id, status, allocation_status, order_date, required_date, shipping_method, shipping_address, billing_address, notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var requiredDate interface{}
	if order.RequiredDate != nil {
		requiredDate = *order.RequiredDate
	}
	var shippingMethod interface{}
	if order.ShippingMethod != nil {
		shippingMethod = *order.ShippingMethod
	}
	var shippingAddress interface{}
	if order.ShippingAddress != nil {
		shippingAddress = *order.ShippingAddress
	}
	var billingAddress interface{}
	if order.BillingAddress != nil {
		billingAddress = *order.BillingAddress
	}
	var notes interface{}
	if order.Notes != nil {
		notes = *order.Notes
	}
	var createdBy interface{}
	if order.CreatedBy != nil {
		createdBy = *order.CreatedBy
	}

	err = tx.QueryRow(ctx, orderQuery,
		order.RefCode, order.CustomerID, order.WarehouseID, order.Status, order.AllocationStatus,
		requiredDate, shippingMethod, shippingAddress, billingAddress, notes, createdBy,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create sales order")
	}

	// Insert items
	if len(items) > 0 {
		itemQuery := `
			INSERT INTO sales_order_items (sales_order_id, product_id, quantity_ordered, quantity_shipped, quantity_reserved, unit_price, allocation_status, created_at, updated_at)
			VALUES ($1, $2, $3, 0, 0, $4, $5, NOW(), NOW())
			RETURNING id, created_at, updated_at
		`
		for _, item := range items {
			item.SalesOrderID = order.ID

			var unitPrice interface{}
			if item.UnitPrice != nil {
				unitPrice = *item.UnitPrice
			}

			err = tx.QueryRow(ctx, itemQuery,
				item.SalesOrderID, item.ProductID, item.QuantityOrdered, unitPrice, item.AllocationStatus,
			).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
			if err != nil {
				return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create sales order item")
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit transaction")
	}

	return order, nil
}

func (r *salesOrderRepository) GetByID(ctx context.Context, orderID int) (*model.SalesOrderDTO, error) {
	query := `
		SELECT id, ref_code, customer_id, warehouse_id, status, allocation_status, order_date, required_date, shipped_date, shipping_method, shipping_address, billing_address, notes, created_by, created_at, updated_at
		FROM sales_orders
		WHERE id = $1
	`
	return scanSalesOrder(ctx, r.db, query, orderID)
}

func (r *salesOrderRepository) GetByRefCode(ctx context.Context, refCode string) (*model.SalesOrderDTO, error) {
	query := `
		SELECT id, ref_code, customer_id, warehouse_id, status, allocation_status, order_date, required_date, shipped_date, shipping_method, shipping_address, billing_address, notes, created_by, created_at, updated_at
		FROM sales_orders
		WHERE ref_code = $1
	`
	return scanSalesOrder(ctx, r.db, query, refCode)
}

func (r *salesOrderRepository) List(ctx context.Context, params *model.SalesOrderParams) (model.SalesOrderDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, ref_code, customer_id, warehouse_id, status, allocation_status, order_date, required_date, shipped_date, shipping_method, shipping_address, billing_address, notes, created_by, created_at, updated_at
		FROM sales_orders
	`

	if params.CustomerID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", len(args)+1))
		args = append(args, *params.CustomerID)
	}
	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *params.Status)
	}
	if params.AllocationStatus != nil {
		conditions = append(conditions, fmt.Sprintf("allocation_status = $%d", len(args)+1))
		args = append(args, *params.AllocationStatus)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, params.Limit, offset)

	return scanSalesOrders(ctx, r.db, query, args...)
}

func (r *salesOrderRepository) Count(ctx context.Context, params *model.SalesOrderParams) (int, error) {
	var (
		args       []any
		conditions []string
	)

	query := `SELECT COUNT(*) FROM sales_orders`

	if params.CustomerID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", len(args)+1))
		args = append(args, *params.CustomerID)
	}
	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *params.Status)
	}
	if params.AllocationStatus != nil {
		conditions = append(conditions, fmt.Sprintf("allocation_status = $%d", len(args)+1))
		args = append(args, *params.AllocationStatus)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count sales orders")
	}
	return count, nil
}

func (r *salesOrderRepository) Update(ctx context.Context, order *model.SalesOrderDTO) error {
	query := `
		UPDATE sales_orders
		SET customer_id = $2, warehouse_id = $3, required_date = $4, shipping_method = $5, shipping_address = $6, billing_address = $7, notes = $8, updated_at = NOW()
		WHERE id = $1
	`

	var requiredDate interface{}
	if order.RequiredDate != nil {
		requiredDate = *order.RequiredDate
	}
	var shippingMethod interface{}
	if order.ShippingMethod != nil {
		shippingMethod = *order.ShippingMethod
	}
	var shippingAddress interface{}
	if order.ShippingAddress != nil {
		shippingAddress = *order.ShippingAddress
	}
	var billingAddress interface{}
	if order.BillingAddress != nil {
		billingAddress = *order.BillingAddress
	}
	var notes interface{}
	if order.Notes != nil {
		notes = *order.Notes
	}

	_, err := r.db.Pool.Exec(ctx, query, order.ID, order.CustomerID, order.WarehouseID, requiredDate, shippingMethod, shippingAddress, billingAddress, notes)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update sales order")
	}
	return nil
}

func (r *salesOrderRepository) UpdateStatus(ctx context.Context, orderID int, status string) error {
	query := `UPDATE sales_orders SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, orderID, status)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update sales order status")
	}
	return nil
}

func (r *salesOrderRepository) UpdateAllocationStatus(ctx context.Context, orderID int, allocationStatus string) error {
	query := `UPDATE sales_orders SET allocation_status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, orderID, allocationStatus)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update sales order allocation status")
	}
	return nil
}

func (r *salesOrderRepository) GetItemsByOrderID(ctx context.Context, orderID int) (model.SalesOrderItemDTOs, error) {
	query := `
		SELECT id, sales_order_id, product_id, quantity_ordered, quantity_shipped, quantity_reserved, unit_price, allocation_status, batch_id, allocated_location_id, created_at, updated_at
		FROM sales_order_items
		WHERE sales_order_id = $1
		ORDER BY id
	`
	return scanSalesOrderItems(ctx, r.db, query, orderID)
}

func (r *salesOrderRepository) GetItemByID(ctx context.Context, itemID int) (*model.SalesOrderItemDTO, error) {
	query := `
		SELECT id, sales_order_id, product_id, quantity_ordered, quantity_shipped, quantity_reserved, unit_price, allocation_status, batch_id, allocated_location_id, created_at, updated_at
		FROM sales_order_items
		WHERE id = $1
	`
	return scanSalesOrderItem(ctx, r.db, query, itemID)
}

func (r *salesOrderRepository) UpdateItem(ctx context.Context, item *model.SalesOrderItemDTO) error {
	query := `
		UPDATE sales_order_items
		SET product_id = $2, quantity_ordered = $3, unit_price = $4, updated_at = NOW()
		WHERE id = $1
	`

	var unitPrice interface{}
	if item.UnitPrice != nil {
		unitPrice = *item.UnitPrice
	}

	_, err := r.db.Pool.Exec(ctx, query, item.ID, item.ProductID, item.QuantityOrdered, unitPrice)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update sales order item")
	}
	return nil
}

func (r *salesOrderRepository) UpdateItemAllocation(ctx context.Context, itemID int, quantityReserved int, allocationStatus string, locationID *int, batchID *int) error {
	query := `
		UPDATE sales_order_items
		SET quantity_reserved = $2, allocation_status = $3, allocated_location_id = $4, batch_id = $5, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, itemID, quantityReserved, allocationStatus, locationID, batchID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update item allocation")
	}
	return nil
}

func (r *salesOrderRepository) UpdateItemShippedQty(ctx context.Context, itemID int, quantityShipped int) error {
	query := `
		UPDATE sales_order_items
		SET quantity_shipped = $2, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, itemID, quantityShipped)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update item shipped quantity")
	}
	return nil
}

func (r *salesOrderRepository) LockOrderForUpdate(ctx context.Context, orderID int) (*model.SalesOrderDTO, error) {
	query := `
		SELECT id, ref_code, customer_id, warehouse_id, status, allocation_status, order_date, required_date, shipped_date, shipping_method, shipping_address, billing_address, notes, created_by, created_at, updated_at
		FROM sales_orders
		WHERE id = $1
		FOR UPDATE
	`
	return scanSalesOrder(ctx, r.db, query, orderID)
}

func (r *salesOrderRepository) LockItemForUpdate(ctx context.Context, itemID int) (*model.SalesOrderItemDTO, error) {
	query := `
		SELECT id, sales_order_id, product_id, quantity_ordered, quantity_shipped, quantity_reserved, unit_price, allocation_status, batch_id, allocated_location_id, created_at, updated_at
		FROM sales_order_items
		WHERE id = $1
		FOR UPDATE
	`
	return scanSalesOrderItem(ctx, r.db, query, itemID)
}

func (r *salesOrderRepository) AllocateStock(ctx context.Context, orderID int, performedBy *int) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start allocation transaction")
	}
	defer tx.Rollback(ctx)

	// Lock order
	orderQuery := `
		SELECT id, ref_code, customer_id, warehouse_id, status, allocation_status, order_date, required_date, shipped_date, shipping_method, shipping_address, billing_address, notes, created_by, created_at, updated_at
		FROM sales_orders
		WHERE id = $1
		FOR UPDATE
	`
	order, err := scanSalesOrderWithTx(ctx, tx, orderQuery, orderID)
	if err != nil {
		return err
	}

	// Check if order can be allocated
	if order.Status == "SHIPPED" || order.Status == "CANCELLED" {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot allocate to shipped or cancelled order")
	}

	// Get order items
	itemsQuery := `
		SELECT id, sales_order_id, product_id, quantity_ordered, quantity_shipped, quantity_reserved, unit_price, allocation_status, batch_id, allocated_location_id, created_at, updated_at
		FROM sales_order_items
		WHERE sales_order_id = $1
		FOR UPDATE
	`
	items, err := scanSalesOrderItemsWithTx(ctx, tx, itemsQuery, orderID)
	if err != nil {
		return err
	}

	// For each item, try to allocate from inventory
	for _, item := range items {
		if item.AllocationStatus == "FULLY_ALLOCATED" || item.QuantityOrdered == item.QuantityReserved {
			continue
		}

		qtyToAllocate := item.QuantityOrdered - item.QuantityReserved

		// Find available inventory in the warehouse
		invQuery := `
			SELECT id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_qty, version
			FROM inventories
			WHERE product_id = $1 AND warehouse_id = $2
			ORDER BY quantity - reserved_qty DESC
			FOR UPDATE
		`
		rows, err := tx.Query(ctx, invQuery, item.ProductID, order.WarehouseID)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query inventory for allocation")
		}

		var allocatedQty int
		var allocLocationID, allocBatchID *int

		for rows.Next() && allocatedQty < qtyToAllocate {
			var inv struct {
				ID          int
				ProductID   int
				WarehouseID int
				LocationID  int
				BatchID     sql.NullInt64
				Quantity    int
				ReservedQty int
				Version     int
			}
			err := rows.Scan(&inv.ID, &inv.ProductID, &inv.WarehouseID, &inv.LocationID, &inv.BatchID, &inv.Quantity, &inv.ReservedQty, &inv.Version)
			if err != nil {
				rows.Close()
				return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan inventory")
			}

			available := inv.Quantity - inv.ReservedQty
			if available <= 0 {
				continue
			}

			allocQty := qtyToAllocate - allocatedQty
			if allocQty > available {
				allocQty = available
			}

			// Update inventory reserved quantity
			updateInvQuery := `
				UPDATE inventories
				SET reserved_qty = reserved_qty + $2, version = version + 1, updated_at = NOW()
				WHERE id = $1 AND version = $3
			`
			cmdTag, err := tx.Exec(ctx, updateInvQuery, inv.ID, allocQty, inv.Version)
			if err != nil {
				rows.Close()
				return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update inventory")
			}
			if cmdTag.RowsAffected() == 0 {
				rows.Close()
				return apperrors.ErrConcurrentUpdate
			}

			// Record stock movement
			var batchID interface{}
			if inv.BatchID.Valid {
				batchID = int(inv.BatchID.Int64)
			}

			movementQuery := `
				INSERT INTO stock_movements (movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id, quantity, reference_type, reference_id, performed_by, notes, created_at)
				VALUES ('RESERVATION', $1, $2, $3, NULL, $4, $5, 'SALES_ORDER', $6, $7, 'Allocation for sales order', NOW())
			`
			_, err = tx.Exec(ctx, movementQuery, item.ProductID, order.WarehouseID, inv.LocationID, batchID, allocQty, orderID, performedBy)
			if err != nil {
				rows.Close()
				return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to record allocation movement")
			}

			allocatedQty += allocQty
			allocLocationID = &inv.LocationID
			if inv.BatchID.Valid {
				bID := int(inv.BatchID.Int64)
				allocBatchID = &bID
			}
		}
		rows.Close()

		// Update item allocation status
		newReservedQty := item.QuantityReserved + allocatedQty
		var newAllocStatus string
		if newReservedQty >= item.QuantityOrdered {
			newAllocStatus = "FULLY_ALLOCATED"
		} else if newReservedQty > 0 {
			newAllocStatus = "PARTIALLY_ALLOCATED"
		} else {
			newAllocStatus = "UNALLOCATED"
		}

		updateItemQuery := `
			UPDATE sales_order_items
			SET quantity_reserved = $2, allocation_status = $3, allocated_location_id = $4, batch_id = $5, updated_at = NOW()
			WHERE id = $1
		`
		_, err = tx.Exec(ctx, updateItemQuery, item.ID, newReservedQty, newAllocStatus, allocLocationID, allocBatchID)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update item allocation status")
		}
	}

	// Update order allocation status
	allItems, err := scanSalesOrderItemsWithTx(ctx, tx, `
		SELECT id, sales_order_id, product_id, quantity_ordered, quantity_shipped, quantity_reserved, unit_price, allocation_status, batch_id, allocated_location_id, created_at, updated_at
		FROM sales_order_items
		WHERE sales_order_id = $1
	`, orderID)
	if err != nil {
		return err
	}

	var totalOrdered, totalReserved int
	for _, item := range allItems {
		totalOrdered += item.QuantityOrdered
		totalReserved += item.QuantityReserved
	}

	var orderAllocStatus string
	if totalReserved == 0 {
		orderAllocStatus = "UNALLOCATED"
	} else if totalReserved >= totalOrdered {
		orderAllocStatus = "FULLY_ALLOCATED"
	} else {
		orderAllocStatus = "PARTIALLY_ALLOCATED"
	}

	_, err = tx.Exec(ctx, `UPDATE sales_orders SET allocation_status = $2, status = 'PROCESSING', updated_at = NOW() WHERE id = $1`, orderID, orderAllocStatus)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update order allocation status")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit allocation transaction")
	}

	return nil
}

func (r *salesOrderRepository) DeallocateStock(ctx context.Context, orderID int) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start deallocation transaction")
	}
	defer tx.Rollback(ctx)

	// Get items with allocations
	items, err := scanSalesOrderItemsWithTx(ctx, tx, `
		SELECT id, sales_order_id, product_id, quantity_ordered, quantity_shipped, quantity_reserved, unit_price, allocation_status, batch_id, allocated_location_id, created_at, updated_at
		FROM sales_order_items
		WHERE sales_order_id = $1 AND quantity_reserved > 0
		FOR UPDATE
	`, orderID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.QuantityReserved <= 0 {
			continue
		}

		// Find the inventory and release reservation
		if item.AllocatedLocationID != nil {
			invQuery := `
				SELECT id, reserved_qty, version
				FROM inventories
				WHERE product_id = $1 AND warehouse_id = $2 AND location_id = $3 AND (($4 IS NULL AND batch_id IS NULL) OR batch_id = $4)
				FOR UPDATE
			`
			var invID, reservedQty, version int
			err := tx.QueryRow(ctx, invQuery, item.ProductID, item.SalesOrderID, *item.AllocatedLocationID, item.BatchID).Scan(&invID, &reservedQty, &version)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to find inventory for deallocation")
			}

			if err == nil {
				newReserved := reservedQty - item.QuantityReserved
				if newReserved < 0 {
					newReserved = 0
				}

				updateInvQuery := `
					UPDATE inventories
					SET reserved_qty = $2, version = version + 1, updated_at = NOW()
					WHERE id = $1
				`
				_, err = tx.Exec(ctx, updateInvQuery, invID, newReserved)
				if err != nil {
					return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to release inventory reservation")
				}
			}
		}

		// Update item
		_, err = tx.Exec(ctx, `
			UPDATE sales_order_items
			SET quantity_reserved = 0, allocation_status = 'UNALLOCATED', allocated_location_id = NULL, updated_at = NOW()
			WHERE id = $1
		`, item.ID)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to deallocate item")
		}
	}

	// Update order
	_, err = tx.Exec(ctx, `
		UPDATE sales_orders
		SET allocation_status = 'UNALLOCATED', updated_at = NOW()
		WHERE id = $1
	`, orderID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update order allocation status")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit deallocation transaction")
	}

	return nil
}

func (r *salesOrderRepository) ShipOrder(ctx context.Context, orderID int, shippedDate *time.Time, performedBy *int) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start shipping transaction")
	}
	defer tx.Rollback(ctx)

	// Lock order
	order, err := scanSalesOrderWithTx(ctx, tx, `
		SELECT id, ref_code, customer_id, warehouse_id, status, allocation_status, order_date, required_date, shipped_date, shipping_method, shipping_address, billing_address, notes, created_by, created_at, updated_at
		FROM sales_orders
		WHERE id = $1
		FOR UPDATE
	`, orderID)
	if err != nil {
		return err
	}

	if order.Status == "SHIPPED" {
		return apperrors.New(apperrors.CodeInvalidOperation, "order already shipped")
	}
	if order.Status == "CANCELLED" {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot ship cancelled order")
	}

	// Get items
	items, err := scanSalesOrderItemsWithTx(ctx, tx, `
		SELECT id, sales_order_id, product_id, quantity_ordered, quantity_shipped, quantity_reserved, unit_price, allocation_status, batch_id, allocated_location_id, created_at, updated_at
		FROM sales_order_items
		WHERE sales_order_id = $1
		FOR UPDATE
	`, orderID)
	if err != nil {
		return err
	}

	for _, item := range items {
		qtyToShip := item.QuantityReserved
		if qtyToShip <= 0 {
			continue
		}

		if item.AllocatedLocationID != nil {
			// Reduce actual quantity and reserved quantity
			invQuery := `
				SELECT id, quantity, reserved_qty, version
				FROM inventories
				WHERE product_id = $1 AND warehouse_id = $2 AND location_id = $3 AND (($4 IS NULL AND batch_id IS NULL) OR batch_id = $4)
				FOR UPDATE
			`
			var invID, quantity, reservedQty, version int
			err := tx.QueryRow(ctx, invQuery, item.ProductID, order.WarehouseID, *item.AllocatedLocationID, item.BatchID).Scan(&invID, &quantity, &reservedQty, &version)
			if err != nil {
				return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to find inventory for shipping")
			}

			newQty := quantity - qtyToShip
			newReserved := reservedQty - qtyToShip
			if newQty < 0 {
				return apperrors.ErrInsufficientStock
			}
			if newReserved < 0 {
				newReserved = 0
			}

			updateInvQuery := `
				UPDATE inventories
				SET quantity = $2, reserved_qty = $3, version = version + 1, updated_at = NOW()
				WHERE id = $1 AND version = $4
			`
			cmdTag, err := tx.Exec(ctx, updateInvQuery, invID, newQty, newReserved, version)
			if err != nil {
				return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update inventory")
			}
			if cmdTag.RowsAffected() == 0 {
				return apperrors.ErrConcurrentUpdate
			}

			// Record SHIP movement
			var batchID interface{}
			if item.BatchID != nil {
				batchID = *item.BatchID
			}

			movementQuery := `
				INSERT INTO stock_movements (movement_type, product_id, warehouse_id, location_id_from, batch_id, quantity, reference_type, reference_id, performed_by, notes, created_at)
				VALUES ('SHIP', $1, $2, $3, $4, $5, 'SALES_ORDER', $6, $7, 'Sales order shipment', NOW())
			`
			_, err = tx.Exec(ctx, movementQuery, item.ProductID, order.WarehouseID, *item.AllocatedLocationID, batchID, qtyToShip, orderID, performedBy)
			if err != nil {
				return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to record ship movement")
			}
		}

		// Update item shipped quantity
		newShippedQty := item.QuantityShipped + qtyToShip
		_, err = tx.Exec(ctx, `
			UPDATE sales_order_items
			SET quantity_shipped = $2, quantity_reserved = quantity_reserved - $3, updated_at = NOW()
			WHERE id = $1
		`, item.ID, newShippedQty, qtyToShip)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update item shipped quantity")
		}
	}

	// Update order status
	shipDate := time.Now()
	if shippedDate != nil {
		shipDate = *shippedDate
	}

	_, err = tx.Exec(ctx, `
		UPDATE sales_orders
		SET status = 'SHIPPED', shipped_date = $2, allocation_status = 'SHIPPED', updated_at = NOW()
		WHERE id = $1
	`, orderID, shipDate)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update order status")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit shipping transaction")
	}

	return nil
}

func (r *salesOrderRepository) ShipOrderItem(ctx context.Context, itemID int, quantity int) error {
	// This is handled within ShipOrder transaction
	return nil
}

// Helper functions for scanning

func scanSalesOrder(ctx context.Context, database *db.DB, query string, args ...any) (*model.SalesOrderDTO, error) {
	var row model.SalesOrderDTO
	var requiredDate sql.NullTime
	var shippedDate sql.NullTime
	var shippingMethod sql.NullString
	var shippingAddress sql.NullString
	var billingAddress sql.NullString
	var notes sql.NullString
	var createdBy sql.NullInt64
	var createdAt, updatedAt sql.NullTime

	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID, &row.RefCode, &row.CustomerID, &row.WarehouseID, &row.Status, &row.AllocationStatus,
		&row.OrderDate, &requiredDate, &shippedDate, &shippingMethod, &shippingAddress, &billingAddress, &notes, &createdBy, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSalesOrderNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	if requiredDate.Valid {
		row.RequiredDate = &requiredDate.Time
	}
	if shippedDate.Valid {
		row.ShippedDate = &shippedDate.Time
	}
	if shippingMethod.Valid {
		row.ShippingMethod = &shippingMethod.String
	}
	if shippingAddress.Valid {
		row.ShippingAddress = &shippingAddress.String
	}
	if billingAddress.Valid {
		row.BillingAddress = &billingAddress.String
	}
	if notes.Valid {
		row.Notes = &notes.String
	}
	if createdBy.Valid {
		v := int(createdBy.Int64)
		row.CreatedBy = &v
	}
	row.ApplyNullScalars(createdAt, updatedAt)

	return &row, nil
}

func scanSalesOrderWithTx(ctx context.Context, tx pgx.Tx, query string, args ...any) (*model.SalesOrderDTO, error) {
	var row model.SalesOrderDTO
	var requiredDate sql.NullTime
	var shippedDate sql.NullTime
	var shippingMethod sql.NullString
	var shippingAddress sql.NullString
	var billingAddress sql.NullString
	var notes sql.NullString
	var createdBy sql.NullInt64
	var createdAt, updatedAt sql.NullTime

	err := tx.QueryRow(ctx, query, args...).Scan(
		&row.ID, &row.RefCode, &row.CustomerID, &row.WarehouseID, &row.Status, &row.AllocationStatus,
		&row.OrderDate, &requiredDate, &shippedDate, &shippingMethod, &shippingAddress, &billingAddress, &notes, &createdBy, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSalesOrderNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	if requiredDate.Valid {
		row.RequiredDate = &requiredDate.Time
	}
	if shippedDate.Valid {
		row.ShippedDate = &shippedDate.Time
	}
	if shippingMethod.Valid {
		row.ShippingMethod = &shippingMethod.String
	}
	if shippingAddress.Valid {
		row.ShippingAddress = &shippingAddress.String
	}
	if billingAddress.Valid {
		row.BillingAddress = &billingAddress.String
	}
	if notes.Valid {
		row.Notes = &notes.String
	}
	if createdBy.Valid {
		v := int(createdBy.Int64)
		row.CreatedBy = &v
	}
	row.ApplyNullScalars(createdAt, updatedAt)

	return &row, nil
}

func scanSalesOrders(ctx context.Context, database *db.DB, query string, args ...any) (model.SalesOrderDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query sales orders")
	}
	defer rows.Close()

	var orders model.SalesOrderDTOs
	for rows.Next() {
		var row model.SalesOrderDTO
		var requiredDate sql.NullTime
		var shippedDate sql.NullTime
		var shippingMethod sql.NullString
		var shippingAddress sql.NullString
		var billingAddress sql.NullString
		var notes sql.NullString
		var createdBy sql.NullInt64
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&row.ID, &row.RefCode, &row.CustomerID, &row.WarehouseID, &row.Status, &row.AllocationStatus,
			&row.OrderDate, &requiredDate, &shippedDate, &shippingMethod, &shippingAddress, &billingAddress, &notes, &createdBy, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan sales order row")
		}

		if requiredDate.Valid {
			row.RequiredDate = &requiredDate.Time
		}
		if shippedDate.Valid {
			row.ShippedDate = &shippedDate.Time
		}
		if shippingMethod.Valid {
			row.ShippingMethod = &shippingMethod.String
		}
		if shippingAddress.Valid {
			row.ShippingAddress = &shippingAddress.String
		}
		if billingAddress.Valid {
			row.BillingAddress = &billingAddress.String
		}
		if notes.Valid {
			row.Notes = &notes.String
		}
		if createdBy.Valid {
			v := int(createdBy.Int64)
			row.CreatedBy = &v
		}
		row.ApplyNullScalars(createdAt, updatedAt)
		orders = append(orders, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate sales orders")
	}

	return orders, nil
}

func scanSalesOrderItem(ctx context.Context, database *db.DB, query string, args ...any) (*model.SalesOrderItemDTO, error) {
	var row model.SalesOrderItemDTO
	var unitPrice sql.NullFloat64
	var batchID sql.NullInt64
	var allocatedLocationID sql.NullInt64
	var createdAt, updatedAt sql.NullTime

	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID, &row.SalesOrderID, &row.ProductID, &row.QuantityOrdered, &row.QuantityShipped, &row.QuantityReserved, &unitPrice, &row.AllocationStatus, &batchID, &allocatedLocationID, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSalesOrderItemNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order item")
	}

	if unitPrice.Valid {
		v := unitPrice.Float64
		row.UnitPrice = &v
	}
	if batchID.Valid {
		v := int(batchID.Int64)
		row.BatchID = &v
	}
	if allocatedLocationID.Valid {
		v := int(allocatedLocationID.Int64)
		row.AllocatedLocationID = &v
	}
	row.ApplyNullScalars(createdAt, updatedAt)

	return &row, nil
}

func scanSalesOrderItems(ctx context.Context, database *db.DB, query string, args ...any) (model.SalesOrderItemDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query sales order items")
	}
	defer rows.Close()

	var items model.SalesOrderItemDTOs
	for rows.Next() {
		var row model.SalesOrderItemDTO
		var unitPrice sql.NullFloat64
		var batchID sql.NullInt64
		var allocatedLocationID sql.NullInt64
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&row.ID, &row.SalesOrderID, &row.ProductID, &row.QuantityOrdered, &row.QuantityShipped, &row.QuantityReserved, &unitPrice, &row.AllocationStatus, &batchID, &allocatedLocationID, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan sales order item row")
		}

		if unitPrice.Valid {
			v := unitPrice.Float64
			row.UnitPrice = &v
		}
		if batchID.Valid {
			v := int(batchID.Int64)
			row.BatchID = &v
		}
		if allocatedLocationID.Valid {
			v := int(allocatedLocationID.Int64)
			row.AllocatedLocationID = &v
		}
		row.ApplyNullScalars(createdAt, updatedAt)
		items = append(items, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate sales order items")
	}

	return items, nil
}

func scanSalesOrderItemsWithTx(ctx context.Context, tx pgx.Tx, query string, args ...any) (model.SalesOrderItemDTOs, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query sales order items")
	}
	defer rows.Close()

	var items model.SalesOrderItemDTOs
	for rows.Next() {
		var row model.SalesOrderItemDTO
		var unitPrice sql.NullFloat64
		var batchID sql.NullInt64
		var allocatedLocationID sql.NullInt64
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&row.ID, &row.SalesOrderID, &row.ProductID, &row.QuantityOrdered, &row.QuantityShipped, &row.QuantityReserved, &unitPrice, &row.AllocationStatus, &batchID, &allocatedLocationID, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan sales order item row")
		}

		if unitPrice.Valid {
			v := unitPrice.Float64
			row.UnitPrice = &v
		}
		if batchID.Valid {
			v := int(batchID.Int64)
			row.BatchID = &v
		}
		if allocatedLocationID.Valid {
			v := int(allocatedLocationID.Int64)
			row.AllocatedLocationID = &v
		}
		row.ApplyNullScalars(createdAt, updatedAt)
		items = append(items, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate sales order items")
	}

	return items, nil
}
