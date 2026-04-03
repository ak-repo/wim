package service

import (
	"context"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/google/uuid"
)

type LocationService interface {
	CreateLocation(ctx context.Context, input *model.CreateLocationRequest) (uuid.UUID, error)
	GetLocationByID(ctx context.Context, locationID uuid.UUID) (*model.LocationResponse, error)
	GetLocationByCode(ctx context.Context, code string) (*model.LocationResponse, error)
	UpdateLocation(ctx context.Context, locationID uuid.UUID, input *model.UpdateLocationRequest) error
	DeleteLocation(ctx context.Context, locationID uuid.UUID) error
	ListLocations(ctx context.Context, params *model.LocationParams) ([]*model.LocationResponse, int, error)
	ListLocationsByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*model.LocationResponse, error)
}

type locationService struct {
	repos *repository.Repositories
}

func NewLocationService(repositories *repository.Repositories) LocationService {
	return &locationService{
		repos: repositories,
	}
}

func (s *locationService) CreateLocation(ctx context.Context, input *model.CreateLocationRequest) (uuid.UUID, error) {
	// Validate input
	if input.WarehouseID == uuid.Nil {
		return uuid.Nil, apperrors.ErrInvalidInput
	}
	if strings.TrimSpace(input.Zone) == "" || strings.TrimSpace(input.LocationCode) == "" || strings.TrimSpace(input.LocationType) == "" {
		return uuid.Nil, apperrors.ErrInvalidInput
	}

	// Check if location code already exists
	exists, err := s.repos.Location.ExistsByCode(ctx, input.LocationCode)
	if err != nil {
		return uuid.Nil, apperrors.ErrCheckingFaild
	}
	if exists {
		return uuid.Nil, apperrors.ErrAlreadyExists
	}

	return s.repos.Location.Create(ctx, input)
}

func (s *locationService) GetLocationByID(ctx context.Context, locationID uuid.UUID) (*model.LocationResponse, error) {
	location, err := s.repos.Location.GetByID(ctx, locationID)
	if err != nil {
		return nil, apperrors.ErrDataBase
	}
	return location.ToAPIResponse(), nil
}

func (s *locationService) GetLocationByCode(ctx context.Context, code string) (*model.LocationResponse, error) {
	location, err := s.repos.Location.GetByCode(ctx, code)
	if err != nil {
		return nil, apperrors.ErrDataBase
	}
	return location.ToAPIResponse(), nil
}

func (s *locationService) UpdateLocation(ctx context.Context, locationID uuid.UUID, input *model.UpdateLocationRequest) error {
	return s.repos.Location.Update(ctx, locationID, input)
}

func (s *locationService) DeleteLocation(ctx context.Context, locationID uuid.UUID) error {
	return s.repos.Location.Delete(ctx, locationID)
}

func (s *locationService) ListLocations(ctx context.Context, params *model.LocationParams) ([]*model.LocationResponse, int, error) {
	// Validate params
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	locations, err := s.repos.Location.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.ErrDataBase
	}

	count, err := s.repos.Location.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.ErrDataBase
	}

	return locations.ToAPIResponse(), count, nil
}

func (s *locationService) ListLocationsByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*model.LocationResponse, error) {
	locations, err := s.repos.Location.ListByWarehouse(ctx, warehouseID)
	if err != nil {
		return nil, apperrors.ErrDataBase
	}
	return locations.ToAPIResponse(), nil
}
