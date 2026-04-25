package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/internal/errs"
)

type WarehouseService interface {
	CreateWarehouse(ctx context.Context, input *model.WarehouseRequest) (int, error)
	GetWarehouseByID(ctx context.Context, warehouseID int) (*model.WarehouseResponse, error)
	GetWarehouseByCode(ctx context.Context, code string) (*model.WarehouseResponse, error)
	UpdateWarehouse(ctx context.Context, warehouseID int, input *model.WarehouseRequest) error
	DeleteWarehouse(ctx context.Context, warehouseID int) error
	ListWarehouses(ctx context.Context, params *model.WarehouseParams) ([]*model.WarehouseResponse, int, error)
	GetWarehouseCount(ctx context.Context, params *model.WarehouseParams) (int, error)
}

type warehouseService struct {
	repos *repository.Repositories
}

func NewWarehouseService(repositories *repository.Repositories) WarehouseService {
	return &warehouseService{
		repos: repositories,
	}
}

func (s *warehouseService) CreateWarehouse(ctx context.Context, input *model.WarehouseRequest) (int, error) {
	if input == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	// Required field validation
	if input.Code == nil || input.Name == nil || input.Country == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "code, name and country are required")
	}
	if strings.TrimSpace(*input.Code) == "" || strings.TrimSpace(*input.Name) == "" || strings.TrimSpace(*input.Country) == "" {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "code, name and country cannot be empty")
	}

	// Check if code already exists
	exists, err := s.repos.Warehouse.ExistsByCode(ctx, *input.Code)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check warehouse code")
	}
	if exists {
		return 0, apperrors.New(apperrors.CodeAlreadyExists, "warehouse with this code already exists")
	}

	// Refcode
	refCode, err := s.repos.RefCode.GenerateWarehouseRefCode(ctx)
	if err != nil {
		return 0, err
	}
	input.RefCode = refCode

	id, err := s.repos.Warehouse.Create(ctx, input)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return 0, apperrors.New(apperrors.CodeAlreadyExists, "warehouse with this code already exists")
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create warehouse")
	}
	return id, nil
}

func (s *warehouseService) GetWarehouseByID(ctx context.Context, warehouseID int) (*model.WarehouseResponse, error) {
	warehouse, err := s.repos.Warehouse.GetByID(ctx, warehouseID)
	if err != nil {
		if errors.Is(err, repository.ErrWarehouseNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "warehouse not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load warehouse")
	}
	return warehouse.ToAPIResponse(), nil
}

func (s *warehouseService) GetWarehouseByCode(ctx context.Context, code string) (*model.WarehouseResponse, error) {
	if strings.TrimSpace(code) == "" {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "code is required")
	}

	warehouse, err := s.repos.Warehouse.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrWarehouseNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "warehouse not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load warehouse")
	}
	return warehouse.ToAPIResponse(), nil
}

func (s *warehouseService) UpdateWarehouse(ctx context.Context, warehouseID int, input *model.WarehouseRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	// PATCH validation - only validate provided fields
	if input.Code != nil {
		if strings.TrimSpace(*input.Code) == "" {
			return apperrors.New(apperrors.CodeInvalidInput, "code cannot be empty")
		}
		// Check code uniqueness if being updated
		exists, err := s.repos.Warehouse.ExistsByCode(ctx, *input.Code)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check warehouse code")
		}
		if exists {
			return apperrors.New(apperrors.CodeAlreadyExists, "warehouse with this code already exists")
		}
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "name cannot be empty")
	}

	if input.Country != nil && strings.TrimSpace(*input.Country) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "country cannot be empty")
	}

	err := s.repos.Warehouse.Update(ctx, warehouseID, input)
	if err != nil {
		if errors.Is(err, repository.ErrWarehouseNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "warehouse not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update warehouse")
	}
	return nil
}

func (s *warehouseService) DeleteWarehouse(ctx context.Context, warehouseID int) error {
	err := s.repos.Warehouse.Delete(ctx, warehouseID)
	if err != nil {
		if errors.Is(err, repository.ErrWarehouseNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "warehouse not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete warehouse")
	}
	return nil
}

func (s *warehouseService) ListWarehouses(ctx context.Context, params *model.WarehouseParams) ([]*model.WarehouseResponse, int, error) {
	// Validate params
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	warehouses, err := s.repos.Warehouse.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list warehouses")
	}

	count, err := s.repos.Warehouse.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count warehouses")
	}

	return warehouses.ToAPIResponse(), count, nil
}

func (s *warehouseService) GetWarehouseCount(ctx context.Context, params *model.WarehouseParams) (int, error) {
	count, err := s.repos.Warehouse.Count(ctx, params)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count warehouses")
	}
	return count, nil
}
