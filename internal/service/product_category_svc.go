package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/pkg/errors"
)

type ProductCategoryService interface {
	CreateProductCategory(ctx context.Context, input *model.ProductCategoryRequest) (int, error)
	GetProductCategoryByID(ctx context.Context, productCategoryID int) (*model.ProductCategoryResponse, error)
	UpdateProductCategory(ctx context.Context, productCategoryID int, input *model.ProductCategoryRequest) error
	DeleteProductCategory(ctx context.Context, productCategoryID int) error
	ListProductCategories(ctx context.Context, params *model.ProductCategoryParams) ([]*model.ProductCategoryResponse, int, error)
}

type productCategoryService struct {
	repos *repository.Repositories
}

func NewProductCategoryService(repositories *repository.Repositories) ProductCategoryService {
	return &productCategoryService{repos: repositories}
}

func (s *productCategoryService) CreateProductCategory(ctx context.Context, input *model.ProductCategoryRequest) (int, error) {
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

	id, err := s.repos.ProductCategory.Create(ctx, input)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return 0, apperrors.New(apperrors.CodeAlreadyExists, "product category with this name already exists")
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create product category")
	}

	return id, nil
}

func (s *productCategoryService) GetProductCategoryByID(ctx context.Context, productCategoryID int) (*model.ProductCategoryResponse, error) {
	productCategory, err := s.repos.ProductCategory.GetByID(ctx, productCategoryID)
	if err != nil {
		if errors.Is(err, repository.ErrProductCategoryNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "product category not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product category")
	}

	return productCategory.ToAPIResponse(), nil
}

func (s *productCategoryService) UpdateProductCategory(ctx context.Context, productCategoryID int, input *model.ProductCategoryRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "name cannot be empty")
	}

	err := s.repos.ProductCategory.Update(ctx, productCategoryID, input)
	if err != nil {
		if errors.Is(err, repository.ErrProductCategoryNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "product category not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update product category")
	}

	return nil
}

func (s *productCategoryService) DeleteProductCategory(ctx context.Context, productCategoryID int) error {
	err := s.repos.ProductCategory.Delete(ctx, productCategoryID)
	if err != nil {
		if errors.Is(err, repository.ErrProductCategoryNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "product category not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete product category")
	}

	return nil
}

func (s *productCategoryService) ListProductCategories(ctx context.Context, params *model.ProductCategoryParams) ([]*model.ProductCategoryResponse, int, error) {
	if params == nil {
		params = &model.ProductCategoryParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	productCategories, err := s.repos.ProductCategory.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list product categories")
	}

	count, err := s.repos.ProductCategory.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count product categories")
	}

	return productCategories.ToAPIResponse(), count, nil
}