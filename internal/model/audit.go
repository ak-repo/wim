package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLogResponse struct {
	ID         uuid.UUID  `json:"id"`
	EntityType string     `json:"entityType"`
	EntityID   uuid.UUID  `json:"entityId"`
	Action     string     `json:"action"`
	UserID     *uuid.UUID `json:"userId,omitempty"`
	OldValues  *string    `json:"oldValues,omitempty"`
	NewValues  *string    `json:"newValues,omitempty"`
	IPAddress  *string    `json:"ipAddress,omitempty"`
	UserAgent  *string    `json:"userAgent,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
}
