package event

import (
	"context"
	"time"
)

type EventType string

const (
	EventProductCreated    EventType = "product.created"
	EventProductUpdated    EventType = "product.updated"
	EventProductDeleted    EventType = "product.deleted"
	EventInventoryAdjusted EventType = "inventory.adjusted"
	EventOrderCreated      EventType = "order.created"
	EventOrderAllocated    EventType = "order.allocated"
	EventOrderShipped      EventType = "order.shipped"
	EventTransferCreated   EventType = "transfer.created"
	EventTransferCompleted EventType = "transfer.completed"
	EventBatchCreated      EventType = "batch.created"
	EventExpiryAlert       EventType = "inventory.expiry_alert"
)

type Event struct {
	ID        string
	Type      EventType
	Payload   []byte
	Timestamp time.Time
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(topic string, handler EventHandler) error
	Close() error
}

type EventHandler func(ctx context.Context, event Event) error
