package inventory

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/event"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrNotFound          = errors.New("inventory not found")
)

type Service struct {
	inventoryRepo     postgres.InventoryRepository
	stockMovementRepo postgres.StockMovementRepository
	batchRepo         postgres.BatchRepository
	auditRepo         postgres.AuditLogRepository
	publisher         event.EventPublisher
}

func NewService(
	inventoryRepo postgres.InventoryRepository,
	stockMovementRepo postgres.StockMovementRepository,
	batchRepo postgres.BatchRepository,
	auditRepo postgres.AuditLogRepository,
	publisher event.EventPublisher,
) *Service {
	return &Service{
		inventoryRepo:     inventoryRepo,
		stockMovementRepo: stockMovementRepo,
		batchRepo:         batchRepo,
		auditRepo:         auditRepo,
		publisher:         publisher,
	}
}

func (s *Service) AdjustInventory(ctx context.Context, input AdjustInput) (*domain.Inventory, error) {
	productID, _ := uuid.Parse(input.ProductID)
	warehouseID, _ := uuid.Parse(input.WarehouseID)
	locationID, _ := uuid.Parse(input.LocationID)

	existing, err := s.inventoryRepo.GetByProductWarehouse(ctx, productID, warehouseID)
	if err == nil && len(existing) > 0 {
		for _, inv := range existing {
			if inv.LocationID == locationID {
				inv.Quantity += input.Quantity
				inv.UpdatedAt = time.Now()
				if err := s.inventoryRepo.UpdateWithVersion(ctx, inv); err != nil {
					return nil, err
				}

				movement := &domain.StockMovement{
					ID:            uuid.New(),
					MovementType:  input.MovementType,
					ProductID:     productID,
					WarehouseID:   warehouseID,
					LocationIDTo:  &locationID,
					Quantity:      input.Quantity,
					ReferenceType: "adjustment",
					ReferenceID:   &inv.ID,
					Notes:         input.Notes,
					CreatedAt:     time.Now(),
				}

				if err := s.stockMovementRepo.Create(ctx, movement); err != nil {
					return nil, err
				}

				s.writeAudit(ctx, "INVENTORY", inv.ID, "UPDATE", nil, inv)
				s.publishEvent(ctx, event.EventInventoryAdjusted, inv)

				return inv, nil
			}
		}
	}

	inv := &domain.Inventory{
		ID:          uuid.New(),
		ProductID:   productID,
		WarehouseID: warehouseID,
		LocationID:  locationID,
		Quantity:    input.Quantity,
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
		MovementType:  input.MovementType,
		ProductID:     productID,
		WarehouseID:   warehouseID,
		LocationIDTo:  &locationID,
		Quantity:      input.Quantity,
		ReferenceType: "adjustment",
		ReferenceID:   &inv.ID,
		Notes:         input.Notes,
		CreatedAt:     time.Now(),
	}

	if err := s.stockMovementRepo.Create(ctx, movement); err != nil {
		return nil, err
	}

	s.writeAudit(ctx, "INVENTORY", inv.ID, "CREATE", nil, inv)
	s.publishEvent(ctx, event.EventInventoryAdjusted, inv)

	return inv, nil
}

func (s *Service) writeAudit(ctx context.Context, entityType string, entityID uuid.UUID, action string, oldValue any, newValue any) {
	if s.auditRepo == nil {
		return
	}

	var oldJSON *string
	if oldValue != nil {
		if b, err := json.Marshal(oldValue); err == nil {
			s := string(b)
			oldJSON = &s
		}
	}

	var newJSON *string
	if newValue != nil {
		if b, err := json.Marshal(newValue); err == nil {
			s := string(b)
			newJSON = &s
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

func (s *Service) GetByWarehouse(ctx context.Context, warehouseID string, limit, offset int) ([]*domain.Inventory, error) {
	id, _ := uuid.Parse(warehouseID)
	return s.inventoryRepo.GetByWarehouse(ctx, id, limit, offset)
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]*domain.Inventory, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.inventoryRepo.List(ctx, limit, offset)
}

func (s *Service) GetByProduct(ctx context.Context, productID string) ([]*domain.Inventory, error) {
	id, _ := uuid.Parse(productID)
	return s.inventoryRepo.GetByProduct(ctx, id)
}

func (s *Service) GetTotalQuantity(ctx context.Context, productID, warehouseID string) (int, error) {
	productUUID, _ := uuid.Parse(productID)
	warehouseUUID, _ := uuid.Parse(warehouseID)
	return s.inventoryRepo.GetTotalQuantity(ctx, productUUID, warehouseUUID)
}

func (s *Service) GetMovements(ctx context.Context, warehouseID string, limit int) ([]*domain.StockMovement, error) {
	id, _ := uuid.Parse(warehouseID)
	return s.stockMovementRepo.GetByWarehouse(ctx, id, limit)
}

type AdjustInput struct {
	ProductID    string
	WarehouseID  string
	LocationID   string
	BatchID      *string
	Quantity     int
	MovementType string
	Notes        string
}
