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
	GetProductCategoryByID(ctx context.Context, categoryID int) (*model.ProductCategoryResponse, error)
	UpdateProductCategory(ctx context.Context, categoryID int, input *model.ProductCategoryRequest) error
	DeleteProductCategory(ctx context.Context, categoryID int) error
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

	if input.Name == nil || strings.TrimSpace(*input.Name) == "" {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "name is required")
	}

	// Check name uniqueness
	exists, err := s.repos.ProductCategory.ExistsByName(ctx, *input.Name)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check category name")
	}
	if exists {
		return 0, apperrors.New(apperrors.CodeAlreadyExists, "product category with this name already exists")
	}

	refCode, err := s.repos.RefCode.GenerateProductCategoryRefCode(ctx)
	if err != nil {
		return 0, err
	}
	input.RefCode = refCode

	id, err := s.repos.ProductCategory.Create(ctx, input)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return 0, apperrors.New(apperrors.CodeAlreadyExists, "product category with this name already exists")
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create product category")
	}

	return id, nil
}

func (s *productCategoryService) GetProductCategoryByID(ctx context.Context, categoryID int) (*model.ProductCategoryResponse, error) {
	category, err := s.repos.ProductCategory.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, repository.ErrProductCategoryNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "product category not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product category")
	}
	return category.ToAPIResponse(), nil
}

func (s *productCategoryService) UpdateProductCategory(ctx context.Context, categoryID int, input *model.ProductCategoryRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			return apperrors.New(apperrors.CodeInvalidInput, "name cannot be empty")
		}
		exists, err := s.repos.ProductCategory.ExistsByName(ctx, *input.Name)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check category name")
		}
		if exists {
			return apperrors.New(apperrors.CodeAlreadyExists, "product category with this name already exists")
		}
	}

	if err := s.repos.ProductCategory.Update(ctx, categoryID, input); err != nil {
		if errors.Is(err, repository.ErrProductCategoryNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "product category not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update product category")
	}

	return nil
}

func (s *productCategoryService) DeleteProductCategory(ctx context.Context, categoryID int) error {
	if err := s.repos.ProductCategory.Delete(ctx, categoryID); err != nil {
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

	categories, err := s.repos.ProductCategory.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list product categories")
	}

	count, err := s.repos.ProductCategory.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count product categories")
	}

	return categories.ToAPIResponse(), count, nil
}
