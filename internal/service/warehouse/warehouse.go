package warehouse

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	warehouseRepo postgres.WarehouseRepository
	locationRepo  postgres.LocationRepository
}

func NewService(warehouseRepo postgres.WarehouseRepository, locationRepo postgres.LocationRepository) *Service {
	return &Service{
		warehouseRepo: warehouseRepo,
		locationRepo:  locationRepo,
	}
}

func (s *Service) CreateWarehouse(ctx context.Context, input CreateWarehouseInput) (*domain.Warehouse, error) {
	warehouse := domain.NewWarehouse(
		input.Code,
		input.Name,
		input.Country,
	)
	warehouse.AddressLine1 = input.AddressLine1
	warehouse.AddressLine2 = input.AddressLine2
	warehouse.City = input.City
	warehouse.State = input.State
	warehouse.PostalCode = input.PostalCode

	if err := s.warehouseRepo.Create(ctx, warehouse); err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (s *Service) GetWarehouse(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error) {
	return s.warehouseRepo.GetByID(ctx, id)
}

func (s *Service) ListWarehouses(ctx context.Context, filter postgres.WarehouseFilter) ([]*domain.Warehouse, error) {
	return s.warehouseRepo.List(ctx, filter)
}

func (s *Service) UpdateWarehouse(ctx context.Context, id uuid.UUID, input CreateWarehouseInput) (*domain.Warehouse, error) {
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	warehouse.Name = input.Name
	warehouse.AddressLine1 = input.AddressLine1
	warehouse.AddressLine2 = input.AddressLine2
	warehouse.City = input.City
	warehouse.State = input.State
	warehouse.PostalCode = input.PostalCode

	if err := s.warehouseRepo.Update(ctx, warehouse); err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (s *Service) DeleteWarehouse(ctx context.Context, id uuid.UUID) error {
	return s.warehouseRepo.Delete(ctx, id)
}

func (s *Service) CreateLocation(ctx context.Context, warehouseID uuid.UUID, input CreateLocationInput) (*domain.Location, error) {
	location := domain.NewLocation(
		warehouseID,
		input.Zone,
		input.Aisle,
		input.Rack,
		input.Bin,
		input.LocationCode,
		input.LocationType,
	)
	location.IsPickFace = input.IsPickFace
	location.MaxWeight = input.MaxWeight

	if err := s.locationRepo.Create(ctx, location); err != nil {
		return nil, err
	}

	return location, nil
}

func (s *Service) GetLocations(ctx context.Context, warehouseID uuid.UUID) ([]*domain.Location, error) {
	return s.locationRepo.GetByWarehouse(ctx, warehouseID)
}

type CreateWarehouseInput struct {
	Code         string
	Name         string
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string
	Country      string
}

type CreateLocationInput struct {
	Zone         string
	Aisle        string
	Rack         string
	Bin          string
	LocationCode string
	LocationType string
	IsPickFace   bool
	MaxWeight    *float64
}
