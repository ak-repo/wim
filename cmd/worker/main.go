package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ak-repo/wim/config"
	"github.com/ak-repo/wim/internal/event"
	"github.com/ak-repo/wim/internal/queue"
	"github.com/ak-repo/wim/internal/worker"
)

func main() {
	log.Println("Starting WIM Worker")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", "error", err)
	}

	// Check if Kafka is configured
	if len(cfg.Kafka.Brokers) == 0 {
		log.Fatal("No Kafka brokers configured")
	}

	// Ensure topics exist
	tm := event.NewTopicManager(cfg.Kafka.Brokers)
	if err := tm.EnsureAllTopics(context.Background()); err != nil {
		log.Printf("Warning: Failed to ensure topics: %v", err)
	}

	// Create pool configuration
	poolConfig := worker.PoolConfig{
		ConsumerConfig: event.ConsumerConfig{
			Brokers: cfg.Kafka.Brokers,
			GroupID: cfg.Kafka.GroupID,
			Topics: []string{
				event.TopicInventoryEvents,
				event.TopicOrderEvents,
				event.TopicStockMovements,
				event.TopicAlerts,
				event.TopicAuditEvents,
				event.TopicSystemEvents,
			},
			MinBytes:       cfg.Kafka.MinBytes,
			MaxBytes:       cfg.Kafka.MaxBytes,
			MaxWait:        cfg.Kafka.MaxWait,
			AutoCommit:     cfg.Kafka.AutoCommit,
			CommitInterval: cfg.Kafka.CommitInterval,
			MaxRetries:     cfg.Worker.RetryCount,
			RetryDelay:     cfg.Worker.RetryDelay,
		},
		ProcessorConfig: queue.ProcessorConfig{
			Workers:    cfg.Worker.PoolSize,
			QueueSize:  cfg.Worker.QueueSize,
			BatchSize:  cfg.Worker.BatchSize,
			MaxRetries: cfg.Worker.RetryCount,
			RetryDelay: cfg.Worker.RetryDelay,
			DLQEnabled: true,
			DLQTopic:   event.TopicDLQ,
		},
		Topics: []string{
			event.TopicInventoryEvents,
			event.TopicOrderEvents,
			event.TopicStockMovements,
			event.TopicAlerts,
			event.TopicAuditEvents,
			event.TopicSystemEvents,
		},
		GroupID: cfg.Kafka.GroupID,
	}

	// Create worker pool
	pool, err := worker.NewPool(poolConfig)
	if err != nil {
		log.Fatalf("Failed to create worker pool: %v", err)
	}

	// Register handlers
	registerHandlers(pool)

	// Start the pool
	if err := pool.Start(); err != nil {
		log.Fatalf("Failed to start worker pool: %v", err)
	}

	log.Println("WIM Worker started successfully")
	log.Printf("Listening on topics: %v", poolConfig.Topics)
	log.Printf("Worker pool size: %d", pool.PoolSize())

	// Wait for shutdown signal
	pool.Wait()

	log.Println("WIM Worker shutdown complete")
}

func registerHandlers(pool *worker.Pool) {
	// Inventory event handler
	pool.RegisterHandler(event.TopicInventoryEvents, func(ctx context.Context, msg *event.ConsumedMessage) error {
		log.Printf("[Inventory] Processing event: %s (ID: %s)", msg.Event.Type, msg.Event.ID)

		switch msg.Event.Type {
		case event.EventProductCreated:
			return handleProductCreated(ctx, msg)
		case event.EventProductUpdated:
			return handleProductUpdated(ctx, msg)
		case event.EventProductDeleted:
			return handleProductDeleted(ctx, msg)
		case event.EventInventoryAdjusted:
			return handleInventoryAdjusted(ctx, msg)
		case event.EventBatchCreated:
			return handleBatchCreated(ctx, msg)
		default:
			log.Printf("[Inventory] Unknown event type: %s", msg.Event.Type)
			return nil
		}
	})

	// Order event handler
	pool.RegisterHandler(event.TopicOrderEvents, func(ctx context.Context, msg *event.ConsumedMessage) error {
		log.Printf("[Order] Processing event: %s (ID: %s)", msg.Event.Type, msg.Event.ID)

		switch msg.Event.Type {
		case event.EventOrderCreated:
			return handleOrderCreated(ctx, msg)
		case event.EventOrderAllocated:
			return handleOrderAllocated(ctx, msg)
		case event.EventOrderShipped:
			return handleOrderShipped(ctx, msg)
		case event.EventTransferCreated:
			return handleTransferCreated(ctx, msg)
		case event.EventTransferCompleted:
			return handleTransferCompleted(ctx, msg)
		default:
			log.Printf("[Order] Unknown event type: %s", msg.Event.Type)
			return nil
		}
	})

	// Stock movement handler
	pool.RegisterHandler(event.TopicStockMovements, func(ctx context.Context, msg *event.ConsumedMessage) error {
		log.Printf("[Stock Movement] Processing movement: %s", msg.Event.ID)
		// Process stock movement for auditing/reporting
		return handleStockMovement(ctx, msg)
	})

	// Alert handler
	pool.RegisterHandler(event.TopicAlerts, func(ctx context.Context, msg *event.ConsumedMessage) error {
		log.Printf("[Alert] Processing alert: %s (Type: %s)", msg.Event.ID, msg.Event.Type)

		switch msg.Event.Type {
		case event.EventExpiryAlert:
			return handleExpiryAlert(ctx, msg)
		default:
			return handleGenericAlert(ctx, msg)
		}
	})

	// Audit events handler
	pool.RegisterHandler(event.TopicAuditEvents, func(ctx context.Context, msg *event.ConsumedMessage) error {
		log.Printf("[Audit] Processing audit event: %s", msg.Event.ID)
		// Store audit events for compliance
		return handleAuditEvent(ctx, msg)
	})

	// System events handler
	pool.RegisterHandler(event.TopicSystemEvents, func(ctx context.Context, msg *event.ConsumedMessage) error {
		log.Printf("[System] Processing system event: %s", msg.Event.ID)
		return handleSystemEvent(ctx, msg)
	})
}

// Handler implementations

func handleProductCreated(ctx context.Context, msg *event.ConsumedMessage) error {
	var data map[string]interface{}
	if err := json.Unmarshal(msg.Event.Payload, &data); err != nil {
		return err
	}
	log.Printf("[Product Created] ID: %s, SKU: %v", msg.Event.ID, data["sku"])
	// TODO: Send notifications, update search index, etc.
	return nil
}

func handleProductUpdated(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[Product Updated] ID: %s", msg.Event.ID)
	// TODO: Invalidate cache, update search index, etc.
	return nil
}

func handleProductDeleted(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[Product Deleted] ID: %s", msg.Event.ID)
	// TODO: Cleanup related data, invalidate cache, etc.
	return nil
}

func handleInventoryAdjusted(ctx context.Context, msg *event.ConsumedMessage) error {
	var data map[string]interface{}
	if err := json.Unmarshal(msg.Event.Payload, &data); err != nil {
		return err
	}
	log.Printf("[Inventory Adjusted] Product: %v, Warehouse: %v, Delta: %v",
		data["product_id"], data["warehouse_id"], data["delta"])
	// TODO: Recalculate stock levels, trigger reorder alerts, etc.
	return nil
}

func handleBatchCreated(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[Batch Created] ID: %s", msg.Event.ID)
	// TODO: Schedule expiry checks, etc.
	return nil
}

func handleOrderCreated(ctx context.Context, msg *event.ConsumedMessage) error {
	var data map[string]interface{}
	if err := json.Unmarshal(msg.Event.Payload, &data); err != nil {
		return err
	}
	log.Printf("[Order Created] Order: %s, Customer: %v", msg.Event.ID, data["customer_id"])
	// TODO: Send notifications, update dashboards, etc.
	return nil
}

func handleOrderAllocated(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[Order Allocated] Order: %s", msg.Event.ID)
	// TODO: Update inventory reservations, generate pick tasks, etc.
	return nil
}

func handleOrderShipped(ctx context.Context, msg *event.ConsumedMessage) error {
	var data map[string]interface{}
	if err := json.Unmarshal(msg.Event.Payload, &data); err != nil {
		return err
	}
	log.Printf("[Order Shipped] Order: %s, Tracking: %v", msg.Event.ID, data["tracking_number"])
	// TODO: Send shipping notifications, update order status, etc.
	return nil
}

func handleTransferCreated(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[Transfer Created] Transfer: %s", msg.Event.ID)
	return nil
}

func handleTransferCompleted(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[Transfer Completed] Transfer: %s", msg.Event.ID)
	return nil
}

func handleStockMovement(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[Stock Movement] Processing: %s", msg.Event.ID)
	// TODO: Update analytics, generate reports, etc.
	return nil
}

func handleExpiryAlert(ctx context.Context, msg *event.ConsumedMessage) error {
	var data map[string]interface{}
	if err := json.Unmarshal(msg.Event.Payload, &data); err != nil {
		return err
	}
	log.Printf("[Expiry Alert] Product: %v, Days: %v", data["product_id"], data["days_remaining"])
	// TODO: Send notifications to warehouse staff, update priority, etc.
	return nil
}

func handleGenericAlert(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[Generic Alert] Type: %s, ID: %s", msg.Event.Type, msg.Event.ID)
	return nil
}

func handleAuditEvent(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[Audit Event] Type: %s, ID: %s, Timestamp: %s",
		msg.Event.Type, msg.Event.ID, msg.Event.Timestamp.Format(time.RFC3339))
	// TODO: Store to audit database, etc.
	return nil
}

func handleSystemEvent(ctx context.Context, msg *event.ConsumedMessage) error {
	log.Printf("[System Event] Type: %s, ID: %s", msg.Event.Type, msg.Event.ID)
	return nil
}
