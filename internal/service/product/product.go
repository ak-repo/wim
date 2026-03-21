package product

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/event"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	repo        postgres.ProductRepository
	barcodeRepo postgres.BarcodeRepository
	auditRepo   postgres.AuditLogRepository
	publisher   event.EventPublisher
}

func NewService(repo postgres.ProductRepository, barcodeRepo postgres.BarcodeRepository, auditRepo postgres.AuditLogRepository, publisher event.EventPublisher) *Service {
	return &Service{repo: repo, barcodeRepo: barcodeRepo, auditRepo: auditRepo, publisher: publisher}
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

	s.writeAudit(ctx, "PRODUCT", product.ID, "CREATE", nil, product)
	s.publishEvent(ctx, event.EventProductCreated, product)

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

	s.writeAudit(ctx, "PRODUCT", product.ID, "UPDATE", nil, product)
	s.publishEvent(ctx, event.EventProductUpdated, product)

	return product, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	existing, _ := s.repo.GetByID(ctx, id)
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.writeAudit(ctx, "PRODUCT", id, "DELETE", existing, nil)
	s.publishEvent(ctx, event.EventProductDeleted, map[string]string{"id": id.String()})
	return nil
}

func (s *Service) writeAudit(ctx context.Context, entityType string, entityID uuid.UUID, action string, oldValue any, newValue any) {
	if s.auditRepo == nil {
		return
	}

	var oldJSON *string
	if oldValue != nil {
		if b, err := json.Marshal(oldValue); err == nil {
			v := string(b)
			oldJSON = &v
		}
	}

	var newJSON *string
	if newValue != nil {
		if b, err := json.Marshal(newValue); err == nil {
			v := string(b)
			newJSON = &v
		}
	}

	_ = s.auditRepo.Create(ctx, &domain.AuditLog{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		OldValues:  oldJSON,
		NewValues:  newJSON,
		CreatedAt:  time.Now(),
	})
}

func (s *Service) publishEvent(ctx context.Context, eventType event.EventType, payload any) {
	if s.publisher == nil {
		return
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return
	}

	_ = s.publisher.Publish(ctx, event.Event{
		ID:        uuid.NewString(),
		Type:      eventType,
		Payload:   b,
		Timestamp: time.Now(),
	})
}
