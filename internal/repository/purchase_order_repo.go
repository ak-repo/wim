package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/internal/db"
	apperrors "github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/model"
	"github.com/jackc/pgx/v5"
)

var ErrPurchaseOrderNotFound = errors.New("purchase order not found")

type PurchaseOrderRepository interface {
	Create(ctx context.Context, order *model.PurchaseOrderDTO, items []*model.PurchaseOrderItemDTO) (*model.PurchaseOrderDTO, error)
	GetByID(ctx context.Context, orderID int) (*model.PurchaseOrderDTO, error)
	GetByRefCode(ctx context.Context, refCode string) (*model.PurchaseOrderDTO, error)
	List(ctx context.Context, params *model.PurchaseOrderParams) (model.PurchaseOrderDTOs, error)
	Count(ctx context.Context, params *model.PurchaseOrderParams) (int, error)
	GetItemsByOrderID(ctx context.Context, orderID int) (model.PurchaseOrderItemDTOs, error)
	Receive(ctx context.Context, orderID int, receivedDate *time.Time, notes string, items []model.ReceivePurchaseOrderItemRequest, performedBy *int) error
	PutAway(ctx context.Context, orderID int, notes string, items []model.PutAwayPurchaseOrderItemRequest, performedBy *int) error
}

type purchaseOrderRepository struct {
	db *db.DB
}

func NewPurchaseOrderRepository(database *db.DB) PurchaseOrderRepository {
	return &purchaseOrderRepository{db: database}
}

func (r *purchaseOrderRepository) Create(ctx context.Context, order *model.PurchaseOrderDTO, items []*model.PurchaseOrderItemDTO) (*model.PurchaseOrderDTO, error) {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start purchase order transaction")
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO purchase_orders (
			ref_code, supplier_id, warehouse_id, status, expected_date, notes, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var expectedDate any
	if order.ExpectedDate != nil {
		expectedDate = *order.ExpectedDate
	}
	var notes any
	if order.Notes != nil {
		notes = *order.Notes
	}
	var createdBy any
	if order.CreatedBy != nil {
		createdBy = *order.CreatedBy
	}

	err = tx.QueryRow(ctx, query, order.RefCode, order.SupplierID, order.WarehouseID, order.Status, expectedDate, notes, createdBy).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create purchase order")
	}

	if len(items) > 0 {
		itemQuery := `
			INSERT INTO purchase_order_items (
				purchase_order_id, product_id, quantity_ordered, quantity_received, batch_number, unit_price, created_at, updated_at
			) VALUES ($1, $2, $3, 0, $4, $5, NOW(), NOW())
			RETURNING id, created_at, updated_at
		`
		for _, item := range items {
			item.PurchaseOrderID = order.ID
			var batchNumber any
			if item.BatchNumber != nil && strings.TrimSpace(*item.BatchNumber) != "" {
				batchNumber = *item.BatchNumber
			}
			var unitPrice any
			if item.UnitPrice != nil {
				unitPrice = *item.UnitPrice
			}
			err = tx.QueryRow(ctx, itemQuery, item.PurchaseOrderID, item.ProductID, item.QuantityOrdered, batchNumber, unitPrice).
				Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
			if err != nil {
				return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create purchase order item")
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit purchase order transaction")
	}

	return order, nil
}

func (r *purchaseOrderRepository) GetByID(ctx context.Context, orderID int) (*model.PurchaseOrderDTO, error) {
	return scanPurchaseOrder(ctx, r.db, `
		SELECT id, ref_code, supplier_id, warehouse_id, status, expected_date, received_date, notes, created_by, created_at, updated_at
		FROM purchase_orders
		WHERE id = $1
	`, orderID)
}

func (r *purchaseOrderRepository) GetByRefCode(ctx context.Context, refCode string) (*model.PurchaseOrderDTO, error) {
	return scanPurchaseOrder(ctx, r.db, `
		SELECT id, ref_code, supplier_id, warehouse_id, status, expected_date, received_date, notes, created_by, created_at, updated_at
		FROM purchase_orders
		WHERE ref_code = $1
	`, refCode)
}

func (r *purchaseOrderRepository) List(ctx context.Context, params *model.PurchaseOrderParams) (model.PurchaseOrderDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, ref_code, supplier_id, warehouse_id, status, expected_date, received_date, notes, created_by, created_at, updated_at
		FROM purchase_orders
	`

	if params.SupplierID != nil {
		conditions = append(conditions, fmt.Sprintf("supplier_id = $%d", len(args)+1))
		args = append(args, *params.SupplierID)
	}
	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *params.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, params.Limit, offset)

	return scanPurchaseOrders(ctx, r.db, query, args...)
}

func (r *purchaseOrderRepository) Count(ctx context.Context, params *model.PurchaseOrderParams) (int, error) {
	var (
		args       []any
		conditions []string
	)

	query := `SELECT COUNT(*) FROM purchase_orders`

	if params.SupplierID != nil {
		conditions = append(conditions, fmt.Sprintf("supplier_id = $%d", len(args)+1))
		args = append(args, *params.SupplierID)
	}
	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *params.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count purchase orders")
	}
	return count, nil
}

func (r *purchaseOrderRepository) GetItemsByOrderID(ctx context.Context, orderID int) (model.PurchaseOrderItemDTOs, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, purchase_order_id, product_id, quantity_ordered, quantity_received, batch_number, unit_price, created_at, updated_at
		FROM purchase_order_items
		WHERE purchase_order_id = $1
		ORDER BY id
	`, orderID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list purchase order items")
	}
	defer rows.Close()

	var items model.PurchaseOrderItemDTOs
	for rows.Next() {
		var row model.PurchaseOrderItemDTO
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&row.ID, &row.PurchaseOrderID, &row.ProductID, &row.QuantityOrdered, &row.QuantityReceived, &row.BatchNumber, &row.UnitPrice, &createdAt, &updatedAt); err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan purchase order item")
		}
		row.ApplyNullScalars(createdAt, updatedAt)
		items = append(items, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate purchase order items")
	}

	return items, nil
}

func (r *purchaseOrderRepository) Receive(ctx context.Context, orderID int, receivedDate *time.Time, notes string, items []model.ReceivePurchaseOrderItemRequest, performedBy *int) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start receiving transaction")
	}
	defer tx.Rollback(ctx)

	order, err := scanPurchaseOrderWithTx(ctx, tx, `
		SELECT id, ref_code, supplier_id, warehouse_id, status, expected_date, received_date, notes, created_by, created_at, updated_at
		FROM purchase_orders
		WHERE id = $1
		FOR UPDATE
	`, orderID)
	if err != nil {
		return err
	}
	if order.Status == constants.StatusReceived || order.Status == constants.StatusCancelled {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot receive this purchase order")
	}

	for _, itemReq := range items {
		if itemReq.QuantityReceived <= 0 {
			return apperrors.New(apperrors.CodeInvalidInput, "quantityReceived must be greater than 0")
		}
		if itemReq.LocationID == nil || *itemReq.LocationID <= 0 {
			return apperrors.New(apperrors.CodeInvalidInput, "locationId is required for receiving")
		}

		poItem, err := scanPurchaseOrderItemWithTx(ctx, tx, `
			SELECT id, purchase_order_id, product_id, quantity_ordered, quantity_received, batch_number, unit_price, created_at, updated_at
			FROM purchase_order_items
			WHERE id = $1 AND purchase_order_id = $2
			FOR UPDATE
		`, itemReq.PurchaseOrderItemID, orderID)
		if err != nil {
			return err
		}

		newReceived := poItem.QuantityReceived + itemReq.QuantityReceived
		if newReceived > poItem.QuantityOrdered {
			return apperrors.New(apperrors.CodeInvalidOperation, "received quantity exceeds ordered quantity")
		}

		location, err := scanLocationWithTx(ctx, tx, `
			SELECT id, ref_code, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at, deleted_at
			FROM locations
			WHERE id = $1 AND deleted_at IS NULL
			FOR UPDATE
		`, *itemReq.LocationID)
		if err != nil {
			return err
		}
		if location.WarehouseID != order.WarehouseID {
			return apperrors.New(apperrors.CodeInvalidInput, "receiving location does not belong to warehouse")
		}

		if err := upsertInventoryForMovement(ctx, tx, order.WarehouseID, poItem.ProductID, *itemReq.LocationID, itemReq.BatchID, itemReq.QuantityReceived, constants.MovementReceipt, constants.ReferencePurchaseOrder, &orderID, performedBy, notes); err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, `
			UPDATE purchase_order_items
			SET quantity_received = $2, updated_at = NOW()
			WHERE id = $1
		`, poItem.ID, newReceived); err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update purchase order item receipt")
		}
	}

	allItems, err := scanPurchaseOrderItemsWithTx(ctx, tx, `
		SELECT id, purchase_order_id, product_id, quantity_ordered, quantity_received, batch_number, unit_price, created_at, updated_at
		FROM purchase_order_items
		WHERE purchase_order_id = $1
	`, orderID)
	if err != nil {
		return err
	}

	var totalOrdered, totalReceived int
	for _, item := range allItems {
		totalOrdered += item.QuantityOrdered
		totalReceived += item.QuantityReceived
	}

	status := constants.StatusPartiallyReceived
	if totalReceived >= totalOrdered && totalOrdered > 0 {
		status = constants.StatusReceived
	}

	var receivedDateValue any
	if status == constants.StatusReceived {
		receivedDateValue = receivedDate
	}

	if _, err := tx.Exec(ctx, `
		UPDATE purchase_orders
		SET status = $2,
			received_date = COALESCE($3, received_date),
			updated_at = NOW()
		WHERE id = $1
	`, orderID, status, receivedDateValue); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update purchase order status")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit receiving transaction")
	}

	return nil
}

func (r *purchaseOrderRepository) PutAway(ctx context.Context, orderID int, notes string, items []model.PutAwayPurchaseOrderItemRequest, performedBy *int) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start put-away transaction")
	}
	defer tx.Rollback(ctx)

	order, err := scanPurchaseOrderWithTx(ctx, tx, `
		SELECT id, ref_code, supplier_id, warehouse_id, status, expected_date, received_date, notes, created_by, created_at, updated_at
		FROM purchase_orders
		WHERE id = $1
		FOR UPDATE
	`, orderID)
	if err != nil {
		return err
	}

	for _, itemReq := range items {
		if itemReq.Quantity <= 0 {
			return apperrors.New(apperrors.CodeInvalidInput, "quantity must be greater than 0")
		}

		poItem, err := scanPurchaseOrderItemWithTx(ctx, tx, `
			SELECT id, purchase_order_id, product_id, quantity_ordered, quantity_received, batch_number, unit_price, created_at, updated_at
			FROM purchase_order_items
			WHERE id = $1 AND purchase_order_id = $2
			FOR UPDATE
		`, itemReq.PurchaseOrderItemID, orderID)
		if err != nil {
			return err
		}
		if itemReq.Quantity > poItem.QuantityReceived {
			return apperrors.New(apperrors.CodeInvalidOperation, "put-away quantity exceeds received quantity")
		}

		fromLocation, err := scanLocationWithTx(ctx, tx, `
			SELECT id, ref_code, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at, deleted_at
			FROM locations
			WHERE id = $1 AND deleted_at IS NULL
			FOR UPDATE
		`, itemReq.FromLocationID)
		if err != nil {
			return err
		}
		toLocation, err := scanLocationWithTx(ctx, tx, `
			SELECT id, ref_code, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at, deleted_at
			FROM locations
			WHERE id = $1 AND deleted_at IS NULL
			FOR UPDATE
		`, itemReq.ToLocationID)
		if err != nil {
			return err
		}
		if fromLocation.WarehouseID != order.WarehouseID || toLocation.WarehouseID != order.WarehouseID {
			return apperrors.New(apperrors.CodeInvalidInput, "put-away locations must belong to the purchase order warehouse")
		}

		if err := moveInventoryBetweenLocations(ctx, tx, order.WarehouseID, poItem.ProductID, itemReq.FromLocationID, itemReq.ToLocationID, itemReq.BatchID, itemReq.Quantity, constants.MovementPutaway, constants.ReferencePurchaseOrder, &orderID, performedBy, notes); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit put-away transaction")
	}

	return nil
}

func scanPurchaseOrder(ctx context.Context, database *db.DB, query string, args ...any) (*model.PurchaseOrderDTO, error) {
	var row model.PurchaseOrderDTO
	var expectedDate, receivedDate, createdAt, updatedAt sql.NullTime
	if err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.RefCode,
		&row.SupplierID,
		&row.WarehouseID,
		&row.Status,
		&expectedDate,
		&receivedDate,
		&row.Notes,
		&row.CreatedBy,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPurchaseOrderNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load purchase order")
	}
	row.ApplyNullScalars(expectedDate, receivedDate, createdAt, updatedAt)
	return &row, nil
}

func scanPurchaseOrders(ctx context.Context, database *db.DB, query string, args ...any) (model.PurchaseOrderDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list purchase orders")
	}
	defer rows.Close()

	var orders model.PurchaseOrderDTOs
	for rows.Next() {
		var row model.PurchaseOrderDTO
		var expectedDate, receivedDate, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&row.ID, &row.RefCode, &row.SupplierID, &row.WarehouseID, &row.Status, &expectedDate, &receivedDate, &row.Notes, &row.CreatedBy, &createdAt, &updatedAt); err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan purchase order row")
		}
		row.ApplyNullScalars(expectedDate, receivedDate, createdAt, updatedAt)
		orders = append(orders, &row)
	}
	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate purchase orders")
	}
	return orders, nil
}

func scanPurchaseOrderWithTx(ctx context.Context, tx pgx.Tx, query string, args ...any) (*model.PurchaseOrderDTO, error) {
	var row model.PurchaseOrderDTO
	var expectedDate, receivedDate, createdAt, updatedAt sql.NullTime
	if err := tx.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.RefCode,
		&row.SupplierID,
		&row.WarehouseID,
		&row.Status,
		&expectedDate,
		&receivedDate,
		&row.Notes,
		&row.CreatedBy,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPurchaseOrderNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load purchase order")
	}
	row.ApplyNullScalars(expectedDate, receivedDate, createdAt, updatedAt)
	return &row, nil
}

func scanPurchaseOrderItemWithTx(ctx context.Context, tx pgx.Tx, query string, args ...any) (*model.PurchaseOrderItemDTO, error) {
	var row model.PurchaseOrderItemDTO
	var createdAt, updatedAt sql.NullTime
	if err := tx.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.PurchaseOrderID,
		&row.ProductID,
		&row.QuantityOrdered,
		&row.QuantityReceived,
		&row.BatchNumber,
		&row.UnitPrice,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPurchaseOrderNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load purchase order item")
	}
	row.ApplyNullScalars(createdAt, updatedAt)
	return &row, nil
}

func scanLocationWithTx(ctx context.Context, tx pgx.Tx, query string, args ...any) (*model.LocationDTO, error) {
	var row model.LocationDTO
	var isPickFace, isActive sql.NullBool
	var createdAt, updatedAt, deletedAt sql.NullTime

	if err := tx.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.RefCode,
		&row.WarehouseID,
		&row.Zone,
		&row.Aisle,
		&row.Rack,
		&row.Bin,
		&row.LocationCode,
		&row.LocationType,
		&isPickFace,
		&row.MaxWeight,
		&isActive,
		&createdAt,
		&updatedAt,
		&deletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load location")
	}

	row.ApplyNullScalars(isActive, isPickFace, createdAt, updatedAt, deletedAt)
	return &row, nil
}

func scanPurchaseOrderItemsWithTx(ctx context.Context, tx pgx.Tx, query string, args ...any) (model.PurchaseOrderItemDTOs, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query purchase order items")
	}
	defer rows.Close()

	var items model.PurchaseOrderItemDTOs
	for rows.Next() {
		var row model.PurchaseOrderItemDTO
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&row.ID, &row.PurchaseOrderID, &row.ProductID, &row.QuantityOrdered, &row.QuantityReceived, &row.BatchNumber, &row.UnitPrice, &createdAt, &updatedAt); err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan purchase order item row")
		}
		row.ApplyNullScalars(createdAt, updatedAt)
		items = append(items, &row)
	}
	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate purchase order items")
	}
	return items, nil
}

func moveInventoryBetweenLocations(ctx context.Context, tx pgx.Tx, warehouseID, productID, fromLocationID, toLocationID int, batchID *int, delta int, movementType, referenceType string, referenceID *int, performedBy *int, notes string) error {
	if delta <= 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "quantity must be greater than 0")
	}

	selectQuery := `
		SELECT id, quantity, reserved_qty
		FROM inventories
		WHERE product_id = $1
		  AND warehouse_id = $2
		  AND location_id = $3
		  AND ((batch_id IS NULL AND $4 IS NULL) OR batch_id = $4)
		FOR UPDATE
	`
	var sourceID, sourceQty, sourceReserved int
	if err := tx.QueryRow(ctx, selectQuery, productID, warehouseID, fromLocationID, batchID).Scan(&sourceID, &sourceQty, &sourceReserved); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.New(apperrors.CodeNotFound, "source inventory not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load source inventory")
	}
	if sourceQty < delta || sourceQty-delta < sourceReserved {
		return apperrors.ErrInsufficientStock
	}

	if _, err := tx.Exec(ctx, `
		UPDATE inventories
		SET quantity = quantity - $2, version = version + 1, updated_at = NOW()
		WHERE id = $1
	`, sourceID, delta); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to decrement source inventory")
	}

	targetID, err := ensureInventoryAtLocation(ctx, tx, productID, warehouseID, toLocationID, batchID, delta)
	if err != nil {
		return err
	}
	_ = targetID

	var locationIDFrom *int
	var locationIDTo *int
	locationIDFrom = &fromLocationID
	locationIDTo = &toLocationID

	var notesValue any
	if strings.TrimSpace(notes) != "" {
		notesValue = notes
	}
	var referenceTypeValue any
	if strings.TrimSpace(referenceType) != "" {
		referenceTypeValue = referenceType
	}

	quantity := int(math.Abs(float64(delta)))
	if _, err := tx.Exec(ctx, `
		INSERT INTO stock_movements (
			movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id,
			quantity, reference_type, reference_id, performed_by, notes, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW())
	`, movementType, productID, warehouseID, locationIDFrom, locationIDTo, batchID, quantity, referenceTypeValue, referenceID, performedBy, notesValue); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create movement")
	}

	return nil
}

func upsertInventoryForMovement(ctx context.Context, tx pgx.Tx, warehouseID, productID, locationID int, batchID *int, delta int, movementType, referenceType string, referenceID *int, performedBy *int, notes string) error {
	if delta <= 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "quantity must be greater than 0")
	}

	selectQuery := `
		SELECT id, quantity, reserved_qty
		FROM inventories
		WHERE product_id = $1
		  AND warehouse_id = $2
		  AND location_id = $3
		  AND ((batch_id IS NULL AND $4 IS NULL) OR batch_id = $4)
		FOR UPDATE
	`
	var inventoryID, quantity, reservedQty int
	err := tx.QueryRow(ctx, selectQuery, productID, warehouseID, locationID, batchID).Scan(&inventoryID, &quantity, &reservedQty)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if _, err := tx.Exec(ctx, `
				INSERT INTO inventories (product_id, warehouse_id, location_id, batch_id, quantity, reserved_qty, version, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, 0, 1, NOW(), NOW())
			`, productID, warehouseID, locationID, batchID, delta); err != nil {
				return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create inventory")
			}
		} else {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load inventory")
		}
	} else {
		newQty := quantity + delta
		if newQty < 0 || newQty < reservedQty {
			return apperrors.ErrInsufficientStock
		}
		if _, err := tx.Exec(ctx, `
			UPDATE inventories
			SET quantity = $2, version = version + 1, updated_at = NOW()
			WHERE id = $1
		`, inventoryID, newQty); err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update inventory")
		}
	}

	var locationIDTo *int = &locationID
	quantityAbs := int(math.Abs(float64(delta)))
	var notesValue any
	if strings.TrimSpace(notes) != "" {
		notesValue = notes
	}
	var referenceTypeValue any
	if strings.TrimSpace(referenceType) != "" {
		referenceTypeValue = referenceType
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO stock_movements (
			movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id,
			quantity, reference_type, reference_id, performed_by, notes, created_at
		) VALUES ($1,$2,$3,NULL,$4,$5,$6,$7,$8,$9,$10,NOW())
	`, movementType, productID, warehouseID, locationIDTo, batchID, quantityAbs, referenceTypeValue, referenceID, performedBy, notesValue); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to record movement")
	}

	return nil
}

func ensureInventoryAtLocation(ctx context.Context, tx pgx.Tx, productID, warehouseID, locationID int, batchID *int, delta int) (int, error) {
	selectQuery := `
		SELECT id, quantity, reserved_qty
		FROM inventories
		WHERE product_id = $1
		  AND warehouse_id = $2
		  AND location_id = $3
		  AND ((batch_id IS NULL AND $4 IS NULL) OR batch_id = $4)
		FOR UPDATE
	`
	var inventoryID, quantity, reservedQty int
	if err := tx.QueryRow(ctx, selectQuery, productID, warehouseID, locationID, batchID).Scan(&inventoryID, &quantity, &reservedQty); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if _, err := tx.Exec(ctx, `
				INSERT INTO inventories (product_id, warehouse_id, location_id, batch_id, quantity, reserved_qty, version, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, 0, 1, NOW(), NOW())
			`, productID, warehouseID, locationID, batchID, delta); err != nil {
				return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create target inventory")
			}
			return 0, nil
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load target inventory")
	}
	if _, err := tx.Exec(ctx, `
		UPDATE inventories
		SET quantity = quantity + $2, version = version + 1, updated_at = NOW()
		WHERE id = $1
	`, inventoryID, delta); err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update target inventory")
	}
	return inventoryID, nil
}
