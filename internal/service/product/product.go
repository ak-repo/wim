package product

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	repo        postgres.ProductRepository
	barcodeRepo postgres.BarcodeRepository
}

func NewService(repo postgres.ProductRepository, barcodeRepo postgres.BarcodeRepository) *Service {
	return &Service{repo: repo, barcodeRepo: barcodeRepo}
}

type CreateInput struct {
	SKU           string
	Name          string
	Description   string
	Category      string
	UnitOfMeasure string
	Weight        *float64
	Barcode       string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*domain.Product, error) {
	existing, err := s.repo.GetBySKU(ctx, input.SKU)
	if err == nil && existing != nil {
		return nil, ErrSKUExists
	}

	product := domain.NewProduct(
		input.SKU,
		input.Name,
		input.Description,
		input.Category,
		input.UnitOfMeasure,
		input.Barcode,
	)

	if input.Weight != nil {
		product.Weight = input.Weight
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	return s.repo.GetBySKU(ctx, sku)
}

func (s *Service) List(ctx context.Context, filter postgres.ProductFilter) ([]*domain.Product, error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input CreateInput) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.SKU != product.SKU {
		existing, err := s.repo.GetBySKU(ctx, input.SKU)
		if err == nil && existing != nil {
			return nil, ErrSKUExists
		}
		product.SKU = input.SKU
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Category = input.Category
	product.UnitOfMeasure = input.UnitOfMeasure
	product.Barcode = input.Barcode
	if input.Weight != nil {
		product.Weight = input.Weight
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
