package order

import (
	"context"
	"errors"
	"time"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
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
}

func NewService(
	purchaseOrderRepo postgres.PurchaseOrderRepository,
	salesOrderRepo postgres.SalesOrderRepository,
	inventoryRepo postgres.InventoryRepository,
	stockMovementRepo postgres.StockMovementRepository,
) *Service {
	return &Service{
		purchaseOrderRepo: purchaseOrderRepo,
		salesOrderRepo:    salesOrderRepo,
		inventoryRepo:     inventoryRepo,
		stockMovementRepo: stockMovementRepo,
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

	return order, nil
}

func (s *Service) GetSalesOrder(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	return s.salesOrderRepo.GetByID(ctx, id)
}

func (s *Service) ListSalesOrders(ctx context.Context, filter postgres.SalesOrderFilter) ([]*domain.SalesOrder, error) {
	return s.salesOrderRepo.List(ctx, filter)
}

func (s *Service) AllocateSalesOrder(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	order, err := s.salesOrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.AllocationStatus != "UNALLOCATED" {
		return nil, ErrInvalidOrderState
	}

	items, err := s.salesOrderRepo.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		available, err := s.inventoryRepo.GetTotalQuantity(ctx, item.ProductID, order.WarehouseID)
		if err != nil {
			return nil, err
		}

		if available < item.QuantityOrdered {
			return nil, errors.New("insufficient stock for product")
		}

		invList, err := s.inventoryRepo.GetByProductWarehouse(ctx, item.ProductID, order.WarehouseID)
		if err != nil {
			return nil, err
		}

		remaining := item.QuantityOrdered
		for _, inv := range invList {
			if remaining <= 0 {
				break
			}
			allocQty := inv.AvailableQty()
			if allocQty > remaining {
				allocQty = remaining
			}

			inv.ReservedQty += allocQty
			inv.UpdatedAt = time.Now()
			if err := s.inventoryRepo.Update(ctx, inv); err != nil {
				return nil, err
			}

			remaining -= allocQty
		}

		item.QuantityAllocated = item.QuantityOrdered
		item.UpdatedAt = time.Now()
		if err := s.salesOrderRepo.UpdateItem(ctx, item); err != nil {
			return nil, err
		}
	}

	order.AllocationStatus = "ALLOCATED"
	order.Status = "PROCESSING"
	order.UpdatedAt = time.Now()

	if err := s.salesOrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) ShipSalesOrder(ctx context.Context, id uuid.UUID) (*domain.SalesOrder, error) {
	order, err := s.salesOrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.Status != "PROCESSING" {
		return nil, ErrInvalidOrderState
	}

	items, err := s.salesOrderRepo.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		invList, err := s.inventoryRepo.GetByProductWarehouse(ctx, item.ProductID, order.WarehouseID)
		if err != nil {
			return nil, err
		}

		remaining := item.QuantityAllocated
		for _, inv := range invList {
			if remaining <= 0 {
				break
			}
			releaseQty := inv.ReservedQty
			if releaseQty > remaining {
				releaseQty = remaining
			}

			inv.Quantity -= releaseQty
			inv.ReservedQty -= releaseQty
			inv.UpdatedAt = time.Now()
			if err := s.inventoryRepo.Update(ctx, inv); err != nil {
				return nil, err
			}

			movement := &domain.StockMovement{
				ID:             uuid.New(),
				MovementType:   "SHIP",
				ProductID:      item.ProductID,
				WarehouseID:    order.WarehouseID,
				LocationIDFrom: &inv.LocationID,
				Quantity:       releaseQty,
				ReferenceType:  "sales_order",
				ReferenceID:    &order.ID,
				Notes:          "Order shipped",
				CreatedAt:      time.Now(),
			}
			if err := s.stockMovementRepo.Create(ctx, movement); err != nil {
				return nil, err
			}

			remaining -= releaseQty
		}

		item.QuantityShipped = item.QuantityAllocated
		item.UpdatedAt = time.Now()
		if err := s.salesOrderRepo.UpdateItem(ctx, item); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	order.ShippedDate = &now
	order.Status = "SHIPPED"
	order.UpdatedAt = time.Now()

	if err := s.salesOrderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
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
