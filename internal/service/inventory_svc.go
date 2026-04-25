package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/internal/errs"
)

type InventoryService interface {
	AdjustInventory(ctx context.Context, input *model.AdjustInventoryRequest) error
	GetInventoryByID(ctx context.Context, inventoryID int) (*model.InventoryResponse, error)
	ListInventory(ctx context.Context, params *model.InventoryParams) ([]*model.InventoryResponse, int, error)
	ListStockMovements(ctx context.Context, params *model.StockMovementParams) ([]*model.StockMovementResponse, int, error)
}

type inventoryService struct {
	repos *repository.Repositories
}

func NewInventoryService(repositories *repository.Repositories) InventoryService {
	return &inventoryService{repos: repositories}
}

func (s *inventoryService) AdjustInventory(ctx context.Context, input *model.AdjustInventoryRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.ProductID <= 0 || input.WarehouseID <= 0 || input.LocationID <= 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "productId, warehouseId, and locationId are required")
	}
	if input.Quantity == 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "quantity cannot be zero")
	}
	if strings.TrimSpace(input.Reason) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "reason is required")
	}

	if _, err := s.repos.Product.GetByID(ctx, input.ProductID); err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "product not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product")
	}

	location, err := s.repos.Location.GetByID(ctx, input.LocationID)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "location not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load location")
	}
	if location.WarehouseID != input.WarehouseID {
		return apperrors.New(apperrors.CodeInvalidInput, "location does not belong to warehouse")
	}

	notes := strings.TrimSpace(input.Reason)
	if strings.TrimSpace(input.Notes) != "" {
		notes = notes + " | " + strings.TrimSpace(input.Notes)
	}

	err = s.repos.Inventory.Adjust(
		ctx,
		input.ProductID,
		input.WarehouseID,
		input.LocationID,
		input.BatchID,
		input.Quantity,
		constants.MovementAdjustment,
		constants.ReferenceManualAdjustment,
		nil,
		nil,
		notes,
	)
	if err != nil {
		if errors.Is(err, apperrors.ErrInsufficientStock) {
			return apperrors.New(apperrors.CodeInsufficientStock, "insufficient stock")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to adjust inventory")
	}

	return nil
}

func (s *inventoryService) GetInventoryByID(ctx context.Context, inventoryID int) (*model.InventoryResponse, error) {
	inventory, err := s.repos.Inventory.GetByID(ctx, inventoryID)
	if err != nil {
		if errors.Is(err, repository.ErrInventoryNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "inventory not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load inventory")
	}
	return inventory.ToAPIResponse(), nil
}

func (s *inventoryService) ListInventory(ctx context.Context, params *model.InventoryParams) ([]*model.InventoryResponse, int, error) {
	if params == nil {
		params = &model.InventoryParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	items, err := s.repos.Inventory.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list inventory")
	}

	count, err := s.repos.Inventory.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count inventory")
	}

	return items.ToAPIResponse(), count, nil
}

func (s *inventoryService) ListStockMovements(ctx context.Context, params *model.StockMovementParams) ([]*model.StockMovementResponse, int, error) {
	if params == nil {
		params = &model.StockMovementParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	movements, err := s.repos.Inventory.ListMovements(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list stock movements")
	}

	count, err := s.repos.Inventory.CountMovements(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count stock movements")
	}

	return movements.ToAPIResponse(), count, nil
}
