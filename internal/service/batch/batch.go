package batch

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	repo postgres.BatchRepository
}

func NewService(repo postgres.BatchRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Batch, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.Batch, error) {
	return s.repo.GetByProduct(ctx, productID)
}

func (s *Service) GetExpiringSoon(ctx context.Context, days int) ([]*domain.Batch, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetExpiringSoon(ctx, days)
}
