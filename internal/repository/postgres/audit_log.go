package postgres

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	GetByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.AuditLog, error)
	GetByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.AuditLog, error)
}

type auditLogRepo struct {
	db *DB
}

func NewAuditLogRepository(db *DB) AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Create(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, entity_type, entity_id, action, user_id, old_values, new_values, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Pool.Exec(ctx, query,
		log.ID, log.EntityType, log.EntityID, log.Action, log.UserID,
		log.OldValues, log.NewValues, log.IPAddress, log.UserAgent, log.CreatedAt,
	)
	return err
}

func (r *auditLogRepo) GetByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, entity_type, entity_id, action, user_id, old_values, new_values, ip_address, user_agent, created_at
		FROM audit_logs WHERE entity_type = $1 AND entity_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		err := rows.Scan(
			&l.ID, &l.EntityType, &l.EntityID, &l.Action, &l.UserID,
			&l.OldValues, &l.NewValues, &l.IPAddress, &l.UserAgent, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, nil
}

func (r *auditLogRepo) GetByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, entity_type, entity_id, action, user_id, old_values, new_values, ip_address, user_agent, created_at
		FROM audit_logs WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		err := rows.Scan(
			&l.ID, &l.EntityType, &l.EntityID, &l.Action, &l.UserID,
			&l.OldValues, &l.NewValues, &l.IPAddress, &l.UserAgent, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, nil
}
