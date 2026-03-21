package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/event"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidOrderState = errors.New("invalid order state")
)

type Service struct {
	purchaseOrderRepo postgres.PurchaseOrderRepository
	salesOrderRepo    postgres.SalesOrderRepository
	inventoryRepo     postgres.InventoryRepository
	stockMovementRepo postgres.StockMovementRepository
	auditRepo         postgres.AuditLogRepository
	publisher         event.EventPublisher
	db                *postgres.DB
}

func NewService(
	purchaseOrderRepo postgres.PurchaseOrderRepository,
	salesOrderRepo postgres.SalesOrderRepository,
	inventoryRepo postgres.InventoryRepository,
	stockMovementRepo postgres.StockMovementRepository,
	auditRepo postgres.AuditLogRepository,
	publisher event.EventPublisher,
	db *postgres.DB,
) *Service {
	return &Service{
		purchaseOrderRepo: purchaseOrderRepo,
		salesOrderRepo:    salesOrderRepo,
		inventoryRepo:     inventoryRepo,
		stockMovementRepo: stockMovementRepo,
		auditRepo:         auditRepo,
		publisher:         publisher,
		db:                db,
	}
}

func (s *Service) CreatePurchaseOrder(ctx context.Context, input CreatePurchaseOrderInput) (*domain.PurchaseOrder, error) {
	po := &domain.PurchaseOrder{
		ID:          uuid.New(),
		PONumber:    generatePONumber(),
		SupplierID:  input.SupplierID,
		WarehouseID: input.WarehouseID,
		OrderDate:   time.Now(),
		Status:      "PENDING",
		TotalAmount: input.TotalAmount,
		Notes:       input.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.purchaseOrderRepo.Create(ctx, po); err != nil {
		return nil, err
	}

	s.writeAudit(ctx, "PURCHASE_ORDER", po.ID, "CREATE", nil, po)
	s.publishEvent(ctx, event.EventOrderCreated, po)

	return po, nil
}

func (s *Service) GetPurchaseOrder(ctx context.Context, id uuid.UUID) (*domain.PurchaseOrder, error) {
	return s.purchaseOrderRepo.GetByID(ctx, id)
}

func (s *Service) ListPurchaseOrders(ctx context.Context, filter postgres.PurchaseOrderFilter) ([]*domain.PurchaseOrder, error) {
	return s.purchaseOrderRepo.List(ctx, filter)
}

func (s *Service) ReceivePurchaseOrder(ctx context.Context, id uuid.UUID, input ReceiveInput) (*domain.PurchaseOrder, error) {
	po, err := s.purchaseOrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if po.Status != "PENDING" && po.Status != "PARTIAL" {
		return nil, ErrInvalidOrderState
	}

	items, err := s.purchaseOrderRepo.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}

	for i, item := range items {
		if item.QuantityOrdered > item.QuantityReceived {
			receiveQty := input.Quantity
			if receiveQty > item.QuantityOrdered-item.QuantityReceived {
				receiveQty = item.QuantityOrdered - item.QuantityReceived
			}

			item.QuantityReceived += receiveQty
			item.UpdatedAt = time.Now()
			if err := s.purchaseOrderRepo.UpdateItem(ctx, item); err != nil {
				return nil, err
			}

			locationID, _ := uuid.Parse(input.Location)
			inv := &domain.Inventory{
				ID:          uuid.New(),
				ProductID:   item.ProductID,
				WarehouseID: po.WarehouseID,
				LocationID:  locationID,
				Quantity:    receiveQty,
				ReservedQty: 0,
				Version:     1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			if err := s.inventoryRepo.Create(ctx, inv); err != nil {
				return nil, err
			}

			movement := &domain.StockMovement{
				ID:            uuid.New(),
				MovementType:  "RECEIPT",
				ProductID:     item.ProductID,
				WarehouseID:   po.WarehouseID,
				LocationIDTo:  &locationID,
				Quantity:      receiveQty,
				ReferenceType: "purchase_order",
				ReferenceID:   &po.ID,
				Notes:         "Received from PO",
				CreatedAt:     time.Now(),
			}
			if err := s.stockMovementRepo.Create(ctx, movement); err != nil {
				return nil, err
			}

			items[i] = item
		}
	}

	allReceived := true
	for _, item := range items {
		if item.QuantityReceived < item.QuantityOrdered {
			allReceived = false
			break
		}
	}

	now := time.Now()
	if allReceived {
		po.Status = "RECEIVED"
	} else {
		po.Status = "PARTIAL"
	}
	po.ReceivedDate = &now
	po.UpdatedAt = now

	if err := s.purchaseOrderRepo.Update(ctx, po); err != nil {
		return nil, err
	}

	s.writeAudit(ctx, "PURCHASE_ORDER", po.ID, "UPDATE", nil, po)
	s.publishEvent(ctx, event.EventInventoryAdjusted, po)

	return po, nil
}

func (s *Service) CreateSalesOrder(ctx context.Context, input CreateSalesOrderInput) (*domain.SalesOrder, error) {
	order := &domain.SalesOrder{
		ID:               uuid.New(),
		OrderNumber:      generateSONumber(),
		CustomerID:       input.CustomerID,
		WarehouseID:      input.WarehouseID,
		OrderDate:        time.Now(),
		Status:           "PENDING",
		AllocationStatus: "UNALLOCATED",
		ShippingMethod:   input.ShippingMethod,
		ShippingAddress:  input.ShippingAddress,
		BillingAddress:   input.BillingAddress,
		Subtotal:         input.Subtotal,
		TotalAmount:      input.TotalAmount,
		Notes:            input.Notes,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.salesOrderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	s.writeAudit(ctx, "SALES_ORDER", order.ID, "CREATE", nil, order)
	s.publishEvent(ctx, event.EventOrderCreated, order)

	return order, nil
}

func (s *Service) GetSalesOrder(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	return s.salesOrderRepo.GetByID(ctx, id)
}

func (s *Service) ListSalesOrders(ctx context.Context, filter postgres.SalesOrderFilter) ([]*domain.SalesOrder, error) {
	return s.salesOrderRepo.List(ctx, filter)
}

func (s *Service) AllocateSalesOrder(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	tx, err := s.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var order domain.SalesOrder
	err = tx.QueryRow(ctx, `
		SELECT id, order_number, customer_id, warehouse_id, order_date, required_date, shipped_date, status, allocation_status, shipping_method, shipping_address, billing_address, subtotal, tax_amount, shipping_amount, total_amount, notes, created_by, created_at, updated_at
		FROM sales_orders WHERE id = $1
		FOR UPDATE`, id).Scan(
		&order.ID, &order.OrderNumber, &order.CustomerID, &order.WarehouseID, &order.OrderDate,
		&order.RequiredDate, &order.ShippedDate, &order.Status, &order.AllocationStatus,
		&order.ShippingMethod, &order.ShippingAddress, &order.BillingAddress, &order.Subtotal,
		&order.TaxAmount, &order.ShippingAmount, &order.TotalAmount, &order.Notes, &order.CreatedBy,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.AllocationStatus != "UNALLOCATED" {
		return nil, ErrInvalidOrderState
	}

	rows, err := tx.Query(ctx, `
		SELECT id, product_id, quantity_ordered
		FROM sales_order_items
		WHERE sales_order_id = $1
		FOR UPDATE`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type itemRow struct {
		ID       uuid.UUID
		Product  uuid.UUID
		Quantity int
	}
	var items []itemRow
	for rows.Next() {
		var it itemRow
		if err := rows.Scan(&it.ID, &it.Product, &it.Quantity); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	for _, item := range items {
		invRows, err := tx.Query(ctx, `
			SELECT id, quantity, reserved_quantity
			FROM inventory
			WHERE product_id = $1 AND warehouse_id = $2
			ORDER BY created_at ASC
			FOR UPDATE`, item.Product, order.WarehouseID)
		if err != nil {
			return nil, err
		}

		type invRow struct {
			ID       uuid.UUID
			Quantity int
			Reserved int
		}
		var invs []invRow
		totalAvailable := 0
		for invRows.Next() {
			var inv invRow
			if err := invRows.Scan(&inv.ID, &inv.Quantity, &inv.Reserved); err != nil {
				invRows.Close()
				return nil, err
			}
			available := inv.Quantity - inv.Reserved
			if available > 0 {
				totalAvailable += available
			}
			invs = append(invs, inv)
		}
		invRows.Close()

		if totalAvailable < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", item.Product)
		}

		remaining := item.Quantity
		for _, inv := range invs {
			if remaining <= 0 {
				break
			}
			available := inv.Quantity - inv.Reserved
			if available <= 0 {
				continue
			}
			alloc := available
			if alloc > remaining {
				alloc = remaining
			}

			if _, err := tx.Exec(ctx, `
				UPDATE inventory
				SET reserved_quantity = reserved_quantity + $1, version = version + 1, updated_at = $2
				WHERE id = $3`, alloc, time.Now(), inv.ID); err != nil {
				return nil, err
			}
			remaining -= alloc
		}

		if _, err := tx.Exec(ctx, `
			UPDATE sales_order_items
			SET quantity_allocated = quantity_ordered, updated_at = $2
			WHERE id = $1`, item.ID, time.Now()); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	if _, err := tx.Exec(ctx, `
		UPDATE sales_orders
		SET allocation_status = 'ALLOCATED', status = 'PROCESSING', updated_at = $2
		WHERE id = $1`, id, now); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	orderUpdated, err := s.salesOrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.writeAudit(ctx, "SALES_ORDER", orderUpdated.ID, "UPDATE", nil, orderUpdated)
	s.publishEvent(ctx, event.EventOrderAllocated, orderUpdated)

	return orderUpdated, nil
}

func (s *Service) ShipSalesOrder(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	tx, err := s.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var order domain.SalesOrder
	err = tx.QueryRow(ctx, `
		SELECT id, order_number, customer_id, warehouse_id, order_date, required_date, shipped_date, status, allocation_status, shipping_method, shipping_address, billing_address, subtotal, tax_amount, shipping_amount, total_amount, notes, created_by, created_at, updated_at
		FROM sales_orders WHERE id = $1
		FOR UPDATE`, id).Scan(
		&order.ID, &order.OrderNumber, &order.CustomerID, &order.WarehouseID, &order.OrderDate,
		&order.RequiredDate, &order.ShippedDate, &order.Status, &order.AllocationStatus,
		&order.ShippingMethod, &order.ShippingAddress, &order.BillingAddress, &order.Subtotal,
		&order.TaxAmount, &order.ShippingAmount, &order.TotalAmount, &order.Notes, &order.CreatedBy,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.Status != "PROCESSING" {
		return nil, ErrInvalidOrderState
	}

	rows, err := tx.Query(ctx, `
		SELECT id, product_id, quantity_allocated
		FROM sales_order_items
		WHERE sales_order_id = $1
		FOR UPDATE`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type itemRow struct {
		ID        uuid.UUID
		ProductID uuid.UUID
		Allocated int
	}
	var items []itemRow
	for rows.Next() {
		var it itemRow
		if err := rows.Scan(&it.ID, &it.ProductID, &it.Allocated); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	for _, item := range items {
		invRows, err := tx.Query(ctx, `
			SELECT id, location_id, reserved_quantity
			FROM inventory
			WHERE product_id = $1 AND warehouse_id = $2
			ORDER BY created_at ASC
			FOR UPDATE`, item.ProductID, order.WarehouseID)
		if err != nil {
			return nil, err
		}

		type invRow struct {
			ID        uuid.UUID
			Location  uuid.UUID
			ReservedQ int
		}
		var invs []invRow
		for invRows.Next() {
			var inv invRow
			if err := invRows.Scan(&inv.ID, &inv.Location, &inv.ReservedQ); err != nil {
				invRows.Close()
				return nil, err
			}
			invs = append(invs, inv)
		}
		invRows.Close()

		remaining := item.Allocated
		for _, inv := range invs {
			if remaining <= 0 {
				break
			}
			if inv.ReservedQ <= 0 {
				continue
			}

			releaseQty := inv.ReservedQ
			if releaseQty > remaining {
				releaseQty = remaining
			}

			if _, err := tx.Exec(ctx, `
				UPDATE inventory
				SET quantity = quantity - $1,
					reserved_quantity = reserved_quantity - $1,
					version = version + 1,
					updated_at = $2
				WHERE id = $3`, releaseQty, time.Now(), inv.ID); err != nil {
				return nil, err
			}

			if _, err := tx.Exec(ctx, `
				INSERT INTO stock_movements (id, movement_type, product_id, warehouse_id, location_id_from, quantity, reference_type, reference_id, notes, created_at)
				VALUES ($1, 'SHIP', $2, $3, $4, $5, 'sales_order', $6, 'Order shipped', $7)`,
				uuid.New(), item.ProductID, order.WarehouseID, inv.Location, releaseQty, order.ID, time.Now()); err != nil {
				return nil, err
			}

			remaining -= releaseQty
		}

		if remaining > 0 {
			return nil, fmt.Errorf("insufficient reserved stock for product %s", item.ProductID)
		}

		if _, err := tx.Exec(ctx, `
			UPDATE sales_order_items
			SET quantity_shipped = quantity_allocated, updated_at = $2
			WHERE id = $1`, item.ID, time.Now()); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	if _, err := tx.Exec(ctx, `
		UPDATE sales_orders
		SET shipped_date = $2, status = 'SHIPPED', updated_at = $2
		WHERE id = $1`, id, now); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	orderUpdated, err := s.salesOrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.writeAudit(ctx, "SALES_ORDER", orderUpdated.ID, "UPDATE", nil, orderUpdated)
	s.publishEvent(ctx, event.EventOrderShipped, orderUpdated)

	return orderUpdated, nil
}

func (s *Service) writeAudit(ctx context.Context, entityType string, entityID uuid.UUID, action string, oldValue any, newValue any) {
	if s.auditRepo == nil {
		return
	}

	var oldJSON *string
	if oldValue != nil {
		if b, err := json.Marshal(oldValue); err == nil {
			v := string(b)
			oldJSON = &v
		}
	}

	var newJSON *string
	if newValue != nil {
		if b, err := json.Marshal(newValue); err == nil {
			v := string(b)
			newJSON = &v
		}
	}

	_ = s.auditRepo.Create(ctx, &domain.AuditLog{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		OldValues:  oldJSON,
		NewValues:  newJSON,
		CreatedAt:  time.Now(),
	})
}

func (s *Service) publishEvent(ctx context.Context, eventType event.EventType, payload any) {
	if s.publisher == nil {
		return
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return
	}

	_ = s.publisher.Publish(ctx, event.Event{
		ID:        uuid.NewString(),
		Type:      eventType,
		Payload:   b,
		Timestamp: time.Now(),
	})
}

func generatePONumber() string {
	return "PO-" + time.Now().Format("20060102150405")
}

func generateSONumber() string {
	return "SO-" + time.Now().Format("20060102150405")
}

type CreatePurchaseOrderInput struct {
	SupplierID   uuid.UUID
	WarehouseID  uuid.UUID
	ExpectedDate *time.Time
	TotalAmount  *float64
	Notes        string
}

type ReceiveInput struct {
	Quantity int
	Location string
}

type CreateSalesOrderInput struct {
	CustomerID      uuid.UUID
	WarehouseID     uuid.UUID
	RequiredDate    *time.Time
	ShippingMethod  *string
	ShippingAddress *string
	BillingAddress  *string
	Subtotal        *float64
	TotalAmount     *float64
	Notes           *string
}
