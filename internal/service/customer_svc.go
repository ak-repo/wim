package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/pkg/errors"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, input *model.CustomerRequest) (int, error)
	GetCustomerByID(ctx context.Context, customerID int) (*model.CustomerResponse, error)
	GetCustomerByEmail(ctx context.Context, email string) (*model.CustomerResponse, error)
	UpdateCustomer(ctx context.Context, customerID int, input *model.CustomerRequest) error
	DeleteCustomer(ctx context.Context, customerID int) error
	ListCustomers(ctx context.Context, params *model.CustomerParams) ([]*model.CustomerResponse, int, error)
}

type customerService struct {
	repos *repository.Repositories
}

func NewCustomerService(repositories *repository.Repositories) CustomerService {
	return &customerService{repos: repositories}
}

func (s *customerService) CreateCustomer(ctx context.Context, input *model.CustomerRequest) (int, error) {
	if input == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.Name == nil || input.Email == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "name and email are required")
	}
	if strings.TrimSpace(*input.Name) == "" || strings.TrimSpace(*input.Email) == "" {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "name and email cannot be empty")
	}

	exists, err := s.repos.Customer.ExistsByEmail(ctx, *input.Email)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check customer email")
	}
	if exists {
		return 0, apperrors.New(apperrors.CodeAlreadyExists, "customer with this email already exists")
	}

	refCode, err := s.repos.RefCode.GenerateCustomerRefCode(ctx)
	if err != nil {
		return 0, err
	}
	input.RefCode = refCode

	id, err := s.repos.Customer.Create(ctx, input)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return 0, apperrors.New(apperrors.CodeAlreadyExists, "customer with this email already exists")
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create customer")
	}

	return id, nil
}

func (s *customerService) GetCustomerByID(ctx context.Context, customerID int) (*model.CustomerResponse, error) {
	customer, err := s.repos.Customer.GetByID(ctx, customerID)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "customer not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load customer")
	}

	return customer.ToAPIResponse(), nil
}

func (s *customerService) GetCustomerByEmail(ctx context.Context, email string) (*model.CustomerResponse, error) {
	if strings.TrimSpace(email) == "" {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "email is required")
	}

	customer, err := s.repos.Customer.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "customer not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load customer")
	}

	return customer.ToAPIResponse(), nil
}

func (s *customerService) UpdateCustomer(ctx context.Context, customerID int, input *model.CustomerRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "name cannot be empty")
	}
	if input.Email != nil {
		if strings.TrimSpace(*input.Email) == "" {
			return apperrors.New(apperrors.CodeInvalidInput, "email cannot be empty")
		}
		exists, err := s.repos.Customer.ExistsByEmail(ctx, *input.Email)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check customer email")
		}
		if exists {
			return apperrors.New(apperrors.CodeAlreadyExists, "customer with this email already exists")
		}
	}

	err := s.repos.Customer.Update(ctx, customerID, input)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "customer not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update customer")
	}

	return nil
}

func (s *customerService) DeleteCustomer(ctx context.Context, customerID int) error {
	err := s.repos.Customer.Delete(ctx, customerID)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "customer not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete customer")
	}

	return nil
}

func (s *customerService) ListCustomers(ctx context.Context, params *model.CustomerParams) ([]*model.CustomerResponse, int, error) {
	if params == nil {
		params = &model.CustomerParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	customers, err := s.repos.Customer.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list customers")
	}

	count, err := s.repos.Customer.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count customers")
	}

	return customers.ToAPIResponse(), count, nil
}
