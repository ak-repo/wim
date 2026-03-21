package audit

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	repo postgres.AuditLogRepository
}

func NewService(repo postgres.AuditLogRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.AuditLog, error) {
	return s.repo.GetByEntity(ctx, entityType, entityID)
}

func (s *Service) GetByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetByUser(ctx, userID, limit)
}
