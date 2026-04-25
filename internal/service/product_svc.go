package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/internal/errs"
)

type ProductService interface {
	CreateProduct(ctx context.Context, input *model.ProductRequest) (int, error)
	GetProductByID(ctx context.Context, productID int) (*model.ProductResponse, error)
	GetProductBySKU(ctx context.Context, sku string) (*model.ProductResponse, error)
	UpdateProduct(ctx context.Context, productID int, input *model.ProductRequest) error
	DeleteProduct(ctx context.Context, productID int) error
	ListProducts(ctx context.Context, params *model.ProductParams) ([]*model.ProductResponse, int, error)
	GetProductCount(ctx context.Context, param *model.ProductParams) (int, error)
}

type productService struct {
	repos *repository.Repositories
}

func NewProductService(repositories *repository.Repositories) ProductService {
	return &productService{repos: repositories}
}

// CREATE
func (s *productService) CreateProduct(ctx context.Context, input *model.ProductRequest) (int, error) {
	if input == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	// Required field validation
	if input.SKU == nil || input.Name == nil || input.UnitOfMeasure == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if strings.TrimSpace(*input.SKU) == "" ||
		strings.TrimSpace(*input.Name) == "" ||
		strings.TrimSpace(*input.UnitOfMeasure) == "" {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	// Check SKU uniqueness
	exists, err := s.repos.Product.ExistsBySKU(ctx, *input.SKU)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check sku uniqueness")
	}
	if exists {
		return 0, apperrors.New(apperrors.CodeAlreadyExists, "product with this sku already exists")
	}

	// Generate RefCode
	refCode, err := s.repos.RefCode.GenerateProductRefCode(ctx)
	if err != nil {
		return 0, err
	}
	input.RefCode = refCode

	id, err := s.repos.Product.Create(ctx, input)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return 0, apperrors.New(apperrors.CodeAlreadyExists, "product with this sku already exists")
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create product")
	}
	return id, nil
}

// GET BY ID
func (s *productService) GetProductByID(ctx context.Context, productID int) (*model.ProductResponse, error) {
	product, err := s.repos.Product.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "product not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product")
	}

	return product.ToAPIResponse(), nil
}

// GET BY SKU
func (s *productService) GetProductBySKU(ctx context.Context, sku string) (*model.ProductResponse, error) {
	if strings.TrimSpace(sku) == "" {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	product, err := s.repos.Product.GetBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "product not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product")
	}

	return product.ToAPIResponse(), nil
}

// UPDATE (PATCH VALIDATION)
func (s *productService) UpdateProduct(ctx context.Context, productID int, input *model.ProductRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	// Validate only provided fields
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if input.SKU != nil && strings.TrimSpace(*input.SKU) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if input.UnitOfMeasure != nil && strings.TrimSpace(*input.UnitOfMeasure) == "" {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	err := s.repos.Product.Update(ctx, productID, input)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "product not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update product")
	}

	return nil
}

// DELETE
func (s *productService) DeleteProduct(ctx context.Context, productID int) error {
	err := s.repos.Product.Delete(ctx, productID)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "product not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete product")
	}

	return nil
}

// LIST
func (s *productService) ListProducts(ctx context.Context, params *model.ProductParams) ([]*model.ProductResponse, int, error) {
	if params == nil {
		params = &model.ProductParams{}
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	products, err := s.repos.Product.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list products")
	}

	count, err := s.repos.Product.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count products")
	}

	return products.ToAPIResponse(), count, nil
}

func (s *productService) GetProductCount(ctx context.Context, param *model.ProductParams) (int, error) {
	count, err := s.repos.Product.Count(ctx, param)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count product")
	}
	return count, nil
}
