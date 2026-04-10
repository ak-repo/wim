package event

import (
	"encoding/json"
	"time"
)

// Topic names for different event types
const (
	TopicInventoryEvents = "wim.inventory.events"
	TopicOrderEvents     = "wim.order.events"
	TopicStockMovements  = "wim.stock.movements"
	TopicAlerts          = "wim.alerts"
	TopicAuditEvents     = "wim.audit.events"
	TopicSystemEvents    = "wim.system.events"
	TopicDLQ             = "wim.dlq"
	TopicRetryPrefix     = "wim.retry."
)

// EventSchema represents the envelope for all events
type EventSchema struct {
	Version       string            `json:"version"`
	ID            string            `json:"id"`
	Type          EventType         `json:"type"`
	Source        string            `json:"source"`
	Timestamp     time.Time         `json:"timestamp"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	Payload       json.RawMessage   `json:"payload"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// EventVersion for schema versioning
const EventVersion = "1.0"

// EventSource identifies the service
const EventSource = "wim-api"

// TopicMapping maps event types to their respective topics
var TopicMapping = map[EventType]string{
	EventProductCreated:    TopicInventoryEvents,
	EventProductUpdated:    TopicInventoryEvents,
	EventProductDeleted:    TopicInventoryEvents,
	EventInventoryAdjusted: TopicInventoryEvents,
	EventOrderCreated:      TopicOrderEvents,
	EventOrderAllocated:    TopicOrderEvents,
	EventOrderShipped:      TopicOrderEvents,
	EventTransferCreated:   TopicOrderEvents,
	EventTransferCompleted: TopicOrderEvents,
	EventBatchCreated:      TopicInventoryEvents,
	EventExpiryAlert:       TopicAlerts,
}

// GetTopicForEvent returns the appropriate topic for an event type
func GetTopicForEvent(eventType EventType) string {
	if topic, ok := TopicMapping[eventType]; ok {
		return topic
	}
	return TopicSystemEvents
}

// RetryTopic returns the retry topic name for a given topic
func RetryTopic(topic string, attempt int) string {
	return TopicRetryPrefix + topic + "." + string(rune('0'+attempt))
}
