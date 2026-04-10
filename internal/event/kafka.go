package event

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
)

// KafkaPublisher is deprecated. Use Producer instead for production-grade features.
// This basic implementation is kept for backwards compatibility.
type KafkaPublisher struct {
	writer *kafka.Writer
}

// NewKafkaPublisher creates a basic Kafka publisher.
// Deprecated: Use NewProducerFromConfig or NewProducer for production use.
func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID),
		Value: payload,
	})
}

func (p *KafkaPublisher) Subscribe(topic string, handler EventHandler) error {
	return errors.New("subscribe is not supported by publisher")
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
