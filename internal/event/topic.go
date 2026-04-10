package event

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// TopicManager handles Kafka topic creation and management
type TopicManager struct {
	brokers []string
}

// TopicConfig represents topic configuration
type TopicConfig struct {
	Topic             string
	NumPartitions     int
	ReplicationFactor int
	RetentionHours    int
	CleanupPolicy     string // "delete" or "compact"
	MinInsyncReplicas int
}

// DefaultTopicConfigs provides default configurations for known topics
var DefaultTopicConfigs = map[string]TopicConfig{
	TopicInventoryEvents: {
		Topic:             TopicInventoryEvents,
		NumPartitions:     6,
		ReplicationFactor: 3,
		RetentionHours:    168, // 7 days
		CleanupPolicy:     "delete",
		MinInsyncReplicas: 2,
	},
	TopicOrderEvents: {
		Topic:             TopicOrderEvents,
		NumPartitions:     6,
		ReplicationFactor: 3,
		RetentionHours:    168,
		CleanupPolicy:     "delete",
		MinInsyncReplicas: 2,
	},
	TopicStockMovements: {
		Topic:             TopicStockMovements,
		NumPartitions:     6,
		ReplicationFactor: 3,
		RetentionHours:    720, // 30 days for audit purposes
		CleanupPolicy:     "delete",
		MinInsyncReplicas: 2,
	},
	TopicAlerts: {
		Topic:             TopicAlerts,
		NumPartitions:     3,
		ReplicationFactor: 3,
		RetentionHours:    24, // 1 day for alerts
		CleanupPolicy:     "delete",
		MinInsyncReplicas: 2,
	},
	TopicAuditEvents: {
		Topic:             TopicAuditEvents,
		NumPartitions:     6,
		ReplicationFactor: 3,
		RetentionHours:    2160, // 90 days for compliance
		CleanupPolicy:     "compact",
		MinInsyncReplicas: 2,
	},
	TopicSystemEvents: {
		Topic:             TopicSystemEvents,
		NumPartitions:     3,
		ReplicationFactor: 3,
		RetentionHours:    72, // 3 days
		CleanupPolicy:     "delete",
		MinInsyncReplicas: 2,
	},
	TopicDLQ: {
		Topic:             TopicDLQ,
		NumPartitions:     3,
		ReplicationFactor: 3,
		RetentionHours:    720, // 30 days
		CleanupPolicy:     "delete",
		MinInsyncReplicas: 2,
	},
}

// NewTopicManager creates a new topic manager
func NewTopicManager(brokers []string) *TopicManager {
	return &TopicManager{
		brokers: brokers,
	}
}

// CreateTopic creates a Kafka topic with the specified configuration
func (tm *TopicManager) CreateTopic(ctx context.Context, config TopicConfig) error {
	conn, err := kafka.Dial("tcp", tm.brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.ConfigEntry{
		{
			ConfigName:  "cleanup.policy",
			ConfigValue: config.CleanupPolicy,
		},
		{
			ConfigName:  "min.insync.replicas",
			ConfigValue: fmt.Sprintf("%d", config.MinInsyncReplicas),
		},
		{
			ConfigName:  "retention.ms",
			ConfigValue: fmt.Sprintf("%d", config.RetentionHours*60*60*1000),
		},
	}

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             config.Topic,
		NumPartitions:     config.NumPartitions,
		ReplicationFactor: config.ReplicationFactor,
		ConfigEntries:     topicConfigs,
	})

	if err != nil {
		// Check if topic already exists
		exists, _ := tm.TopicExists(ctx, config.Topic)
		if exists {
			return nil // Not an error
		}
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}

// CreateTopics creates multiple topics
func (tm *TopicManager) CreateTopics(ctx context.Context, configs []TopicConfig) error {
	for _, config := range configs {
		if err := tm.CreateTopic(ctx, config); err != nil {
			return fmt.Errorf("failed to create topic %s: %w", config.Topic, err)
		}
	}
	return nil
}

// CreateDefaultTopics creates all default WIM topics
func (tm *TopicManager) CreateDefaultTopics(ctx context.Context) error {
	for _, config := range DefaultTopicConfigs {
		if err := tm.CreateTopic(ctx, config); err != nil {
			return fmt.Errorf("failed to create topic %s: %w", config.Topic, err)
		}
	}
	return nil
}

// TopicExists checks if a topic exists
func (tm *TopicManager) TopicExists(ctx context.Context, topic string) (bool, error) {
	conn, err := kafka.Dial("tcp", tm.brokers[0])
	if err != nil {
		return false, fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		if err == kafka.UnknownTopicOrPartition {
			return false, nil
		}
		return false, fmt.Errorf("failed to read partitions: %w", err)
	}

	return len(partitions) > 0, nil
}

// DeleteTopic deletes a topic
func (tm *TopicManager) DeleteTopic(ctx context.Context, topic string) error {
	conn, err := kafka.Dial("tcp", tm.brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to controller: %w", err)
	}
	defer controllerConn.Close()

	return controllerConn.DeleteTopics(topic)
}

// ListTopics lists all topics
func (tm *TopicManager) ListTopics(ctx context.Context) ([]string, error) {
	conn, err := kafka.Dial("tcp", tm.brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions: %w", err)
	}

	topicMap := make(map[string]bool)
	for _, p := range partitions {
		topicMap[p.Topic] = true
	}

	topics := make([]string, 0, len(topicMap))
	for topic := range topicMap {
		topics = append(topics, topic)
	}

	return topics, nil
}

// GetTopicInfo returns information about a topic
func (tm *TopicManager) GetTopicInfo(ctx context.Context, topic string) (*TopicInfo, error) {
	conn, err := kafka.Dial("tcp", tm.brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions: %w", err)
	}

	return &TopicInfo{
		Topic:         topic,
		NumPartitions: len(partitions),
	}, nil
}

// TopicInfo holds information about a topic
type TopicInfo struct {
	Topic         string
	NumPartitions int
}

// EnsureAllTopics ensures all required topics exist with proper configuration
func (tm *TopicManager) EnsureAllTopics(ctx context.Context) error {
	// Add retry mechanism
	maxRetries := 5
	retryDelay := 2 * time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			time.Sleep(retryDelay)
		}

		err := tm.CreateDefaultTopics(ctx)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("failed to ensure topics after %d retries: %w", maxRetries, lastErr)
}
