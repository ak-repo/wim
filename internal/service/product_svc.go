package service

import (
	"context"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/pkg/errors"
)

type ProductService interface {
	CreateProduct(ctx context.Context, input *model.CreateProductRequest) (int, error)
	GetProductByID(ctx context.Context, productID int) (*model.ProductResponse, error)
	GetProductBySKU(ctx context.Context, sku string) (*model.ProductResponse, error)
	UpdateProduct(ctx context.Context, productID int, input *model.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, productID int) error
	ListProducts(ctx context.Context, params *model.ProductParams) ([]*model.ProductResponse, int, error)
}

type productService struct {
	repos *repository.Repositories
}

func NewProductService(repositories *repository.Repositories) ProductService {
	return &productService{
		repos: repositories,
	}
}

func (s *productService) CreateProduct(ctx context.Context, input *model.CreateProductRequest) (int, error) {
	// Validate input
	if strings.TrimSpace(input.SKU) == "" || strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.UnitOfMeasure) == "" {
		return 0, apperrors.ErrInvalidInput
	}

	// Check if SKU already exists
	exists, err := s.repos.Product.ExistsBySKU(ctx, input.SKU)
	if err != nil {
		return 0, apperrors.ErrCheckingFaild
	}
	if exists {
		return 0, apperrors.ErrAlreadyExists
	}

	return s.repos.Product.Create(ctx, input)
}

func (s *productService) GetProductByID(ctx context.Context, productID int) (*model.ProductResponse, error) {
	product, err := s.repos.Product.GetByID(ctx, productID)
	if err != nil {
		return nil, apperrors.ErrDataBase
	}
	return product.ToAPIResponse(), nil
}

func (s *productService) GetProductBySKU(ctx context.Context, sku string) (*model.ProductResponse, error) {
	product, err := s.repos.Product.GetBySKU(ctx, sku)
	if err != nil {
		return nil, apperrors.ErrDataBase
	}
	return product.ToAPIResponse(), nil
}

func (s *productService) UpdateProduct(ctx context.Context, productID int, input *model.UpdateProductRequest) error {
	return s.repos.Product.Update(ctx, productID, input)
}

func (s *productService) DeleteProduct(ctx context.Context, productID int) error {
	return s.repos.Product.Delete(ctx, productID)
}

func (s *productService) ListProducts(ctx context.Context, params *model.ProductParams) ([]*model.ProductResponse, int, error) {
	// Validate params
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	products, err := s.repos.Product.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.ErrDataBase
	}

	count, err := s.repos.Product.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.ErrDataBase
	}

	return products.ToAPIResponse(), count, nil
}
