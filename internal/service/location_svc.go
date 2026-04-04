package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/pkg/errors"
)

type LocationService interface {
	CreateLocation(ctx context.Context, input *model.LocationRequest) (int, error)
	GetLocationByID(ctx context.Context, locationID int) (*model.LocationResponse, error)
	GetLocationByCode(ctx context.Context, code string) (*model.LocationResponse, error)
	UpdateLocation(ctx context.Context, locationID int, input *model.LocationRequest) error
	DeleteLocation(ctx context.Context, locationID int) error
	ListLocations(ctx context.Context, params *model.LocationParams) ([]*model.LocationResponse, int, error)
	ListLocationsByWarehouse(ctx context.Context, warehouseID int) ([]*model.LocationResponse, error)
}

type locationService struct {
	repos *repository.Repositories
}

func NewLocationService(repositories *repository.Repositories) LocationService {
	return &locationService{
		repos: repositories,
	}
}

func (s *locationService) CreateLocation(ctx context.Context, input *model.LocationRequest) (int, error) {
	if input == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	// Validate input
	if input.WarehouseID == 0 {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "warehouseId is required")
	}
	if input.Zone == nil || input.LocationCode == nil || input.LocationType == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "zone, locationCode and locationType are required")
	}
	if strings.TrimSpace(*input.Zone) == "" || strings.TrimSpace(*input.LocationCode) == "" || strings.TrimSpace(*input.LocationType) == "" {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "zone, locationCode and locationType cannot be empty")
	}

	// Check if location code already exists
	exists, err := s.repos.Location.ExistsByCode(ctx, *input.LocationCode)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check location code")
	}
	if exists {
		return 0, apperrors.New(apperrors.CodeAlreadyExists, "location with this code already exists")
	}

	// Refcode
	refCode, err := s.repos.RefCode.GenerateLocationRefCode(ctx)
	if err != nil {
		return 0, err
	}
	input.RefCode = refCode

	id, err := s.repos.Location.Create(ctx, input)
	if err != nil {
		if errors.Is(err, repository.ErrWarehouseNotFound) {
			return 0, apperrors.New(apperrors.CodeNotFound, "warehouse not found")
		}
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return 0, apperrors.New(apperrors.CodeAlreadyExists, "location with this code already exists")
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create location")
	}
	return id, nil
}

func (s *locationService) GetLocationByID(ctx context.Context, locationID int) (*model.LocationResponse, error) {
	location, err := s.repos.Location.GetByID(ctx, locationID)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "location not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load location")
	}
	return location.ToAPIResponse(), nil
}

func (s *locationService) GetLocationByCode(ctx context.Context, code string) (*model.LocationResponse, error) {
	if strings.TrimSpace(code) == "" {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "code is required")
	}

	location, err := s.repos.Location.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "location not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load location")
	}
	return location.ToAPIResponse(), nil
}

func (s *locationService) UpdateLocation(ctx context.Context, locationID int, input *model.LocationRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	// PATCH validation - only validate provided fields
	if input.Zone != nil && strings.TrimSpace(*input.Zone) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "zone cannot be empty")
	}

	if input.LocationCode != nil {
		if strings.TrimSpace(*input.LocationCode) == "" {
			return apperrors.New(apperrors.CodeInvalidInput, "locationCode cannot be empty")
		}
		// Check code uniqueness if being updated
		exists, err := s.repos.Location.ExistsByCode(ctx, *input.LocationCode)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check location code")
		}
		if exists {
			return apperrors.New(apperrors.CodeAlreadyExists, "location with this code already exists")
		}
	}

	if input.LocationType != nil && strings.TrimSpace(*input.LocationType) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "locationType cannot be empty")
	}

	err := s.repos.Location.Update(ctx, locationID, input)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "location not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update location")
	}
	return nil
}

func (s *locationService) DeleteLocation(ctx context.Context, locationID int) error {
	err := s.repos.Location.Delete(ctx, locationID)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "location not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete location")
	}
	return nil
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
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list locations")
	}

	count, err := s.repos.Location.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count locations")
	}

	return locations.ToAPIResponse(), count, nil
}

func (s *locationService) ListLocationsByWarehouse(ctx context.Context, warehouseID int) ([]*model.LocationResponse, error) {
	locations, err := s.repos.Location.ListByWarehouse(ctx, warehouseID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list locations by warehouse")
	}
	return locations.ToAPIResponse(), nil
}
