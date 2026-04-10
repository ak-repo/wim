package event

import (
	"context"
	"log"
)

// MockPublisher is a no-op event publisher for when Kafka is not available
type MockPublisher struct{}

func NewMockPublisher() EventPublisher {
	return &MockPublisher{}
}

func (m *MockPublisher) Publish(ctx context.Context, event Event) error {
	log.Printf("[MOCK EVENT] Type: %s, ID: %s", event.Type, event.ID)
	return nil
}

func (m *MockPublisher) Subscribe(topic string, handler EventHandler) error {
	return nil
}

func (m *MockPublisher) Close() error {
	return nil
}
