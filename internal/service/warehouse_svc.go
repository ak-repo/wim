package service

import (
	"context"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/google/uuid"
)

type WarehouseService interface {
	CreateWarehouse(ctx context.Context, input *model.CreateWarehouseRequest) (uuid.UUID, error)
	GetWarehouseByID(ctx context.Context, warehouseID uuid.UUID) (*model.WarehouseResponse, error)
	GetWarehouseByCode(ctx context.Context, code string) (*model.WarehouseResponse, error)
	UpdateWarehouse(ctx context.Context, warehouseID uuid.UUID, input *model.UpdateWarehouseRequest) error
	DeleteWarehouse(ctx context.Context, warehouseID uuid.UUID) error
	ListWarehouses(ctx context.Context, params *model.WarehouseParams) ([]*model.WarehouseResponse, int, error)
}

type warehouseService struct {
	repos *repository.Repositories
}

func NewWarehouseService(repositories *repository.Repositories) WarehouseService {
	return &warehouseService{
		repos: repositories,
	}
}

func (s *warehouseService) CreateWarehouse(ctx context.Context, input *model.CreateWarehouseRequest) (uuid.UUID, error) {
	// Validate input
	if strings.TrimSpace(input.Code) == "" || strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Country) == "" {
		return uuid.Nil, apperrors.ErrInvalidInput
	}

	// Check if code already exists
	exists, err := s.repos.Warehouse.ExistsByCode(ctx, input.Code)
	if err != nil {
		return uuid.Nil, apperrors.ErrCheckingFaild
	}
	if exists {
		return uuid.Nil, apperrors.ErrAlreadyExists
	}

	return s.repos.Warehouse.Create(ctx, input)
}

func (s *warehouseService) GetWarehouseByID(ctx context.Context, warehouseID uuid.UUID) (*model.WarehouseResponse, error) {
	warehouse, err := s.repos.Warehouse.GetByID(ctx, warehouseID)
	if err != nil {
		return nil, apperrors.ErrDataBase
	}
	return warehouse.ToAPIResponse(), nil
}

func (s *warehouseService) GetWarehouseByCode(ctx context.Context, code string) (*model.WarehouseResponse, error) {
	warehouse, err := s.repos.Warehouse.GetByCode(ctx, code)
	if err != nil {
		return nil, apperrors.ErrDataBase
	}
	return warehouse.ToAPIResponse(), nil
}

func (s *warehouseService) UpdateWarehouse(ctx context.Context, warehouseID uuid.UUID, input *model.UpdateWarehouseRequest) error {
	return s.repos.Warehouse.Update(ctx, warehouseID, input)
}

func (s *warehouseService) DeleteWarehouse(ctx context.Context, warehouseID uuid.UUID) error {
	return s.repos.Warehouse.Delete(ctx, warehouseID)
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
		return nil, 0, apperrors.ErrDataBase
	}

	count, err := s.repos.Warehouse.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.ErrDataBase
	}

	return warehouses.ToAPIResponse(), count, nil
}
