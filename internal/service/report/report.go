package report

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	inventoryRepo     postgres.InventoryRepository
	stockMovementRepo postgres.StockMovementRepository
	batchRepo         postgres.BatchRepository
}

func NewService(
	inventoryRepo postgres.InventoryRepository,
	stockMovementRepo postgres.StockMovementRepository,
	batchRepo postgres.BatchRepository,
) *Service {
	return &Service{
		inventoryRepo:     inventoryRepo,
		stockMovementRepo: stockMovementRepo,
		batchRepo:         batchRepo,
	}
}

func (s *Service) InventoryByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]*domain.Inventory, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.inventoryRepo.GetByWarehouse(ctx, warehouseID, limit, offset)
}

func (s *Service) MovementsByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit int) ([]*domain.StockMovement, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.stockMovementRepo.GetByWarehouse(ctx, warehouseID, limit)
}

func (s *Service) ExpiryReport(ctx context.Context, days int) ([]*domain.Batch, error) {
	if days <= 0 {
		days = 30
	}
	return s.batchRepo.GetExpiringSoon(ctx, days)
}
