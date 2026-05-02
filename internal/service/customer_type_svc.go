package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/pkg/errors"
)

type CustomerTypeService interface {
	CreateCustomerType(ctx context.Context, input *model.CustomerTypeRequest) (int, error)
	GetCustomerTypeByID(ctx context.Context, customerTypeID int) (*model.CustomerTypeResponse, error)
	UpdateCustomerType(ctx context.Context, customerTypeID int, input *model.CustomerTypeRequest) error
	DeleteCustomerType(ctx context.Context, customerTypeID int) error
	ListCustomerTypes(ctx context.Context, params *model.CustomerTypeParams) ([]*model.CustomerTypeResponse, int, error)
}

type customerTypeService struct {
	repos *repository.Repositories
}

func NewCustomerTypeService(repositories *repository.Repositories) CustomerTypeService {
	return &customerTypeService{repos: repositories}
}

func (s *customerTypeService) CreateCustomerType(ctx context.Context, input *model.CustomerTypeRequest) (int, error) {
	if input == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.Name == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "name is required")
	}
	if strings.TrimSpace(*input.Name) == "" {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "name cannot be empty")
	}

	if input.IsActive == nil {
		trueVal := true
		input.IsActive = &trueVal
	}

	id, err := s.repos.CustomerType.Create(ctx, input)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return 0, apperrors.New(apperrors.CodeAlreadyExists, "customer type with this name already exists")
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create customer type")
	}

	return id, nil
}

func (s *customerTypeService) GetCustomerTypeByID(ctx context.Context, customerTypeID int) (*model.CustomerTypeResponse, error) {
	customerType, err := s.repos.CustomerType.GetByID(ctx, customerTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerTypeNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "customer type not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load customer type")
	}

	return customerType.ToAPIResponse(), nil
}

func (s *customerTypeService) UpdateCustomerType(ctx context.Context, customerTypeID int, input *model.CustomerTypeRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "name cannot be empty")
	}

	err := s.repos.CustomerType.Update(ctx, customerTypeID, input)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerTypeNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "customer type not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update customer type")
	}

	return nil
}

func (s *customerTypeService) DeleteCustomerType(ctx context.Context, customerTypeID int) error {
	err := s.repos.CustomerType.Delete(ctx, customerTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerTypeNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "customer type not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete customer type")
	}

	return nil
}

func (s *customerTypeService) ListCustomerTypes(ctx context.Context, params *model.CustomerTypeParams) ([]*model.CustomerTypeResponse, int, error) {
	if params == nil {
		params = &model.CustomerTypeParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	customerTypes, err := s.repos.CustomerType.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list customer types")
	}

	count, err := s.repos.CustomerType.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count customer types")
	}

	return customerTypes.ToAPIResponse(), count, nil
}