# Production-Grade Kafka Implementation

This document describes the production-grade Kafka event system implemented for the Warehouse Inventory Management (WIM) system.

## Overview

The event system has been redesigned from a basic single-topic publisher to a comprehensive, production-ready event streaming platform with:

- **Multi-topic architecture** for different event types
- **Producer** with async/sync modes, batching, compression, and retries
- **Consumer** with consumer groups, manual offset management, and graceful shutdown
- **Worker pool** with configurable concurrency, DLQ support, and metrics
- **Topic management** with auto-creation and proper partitioning
- **Schema versioning** for event compatibility

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           WIM Kafka System                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐        │
│  │   API Server │──────▶│   Producer   │──────▶│   Kafka      │        │
│  │              │      │   (Async)    │      │   Cluster    │        │
│  └──────────────┘      └──────────────┘      └──────┬───────┘        │
│                                                      │                  │
│                       ┌──────────────────────────────┼────────────────┐│
│                       │                              │                ││
│                       ▼                              ▼                ││
│              ┌────────────────┐           ┌────────────────┐           ││
│              │wim.inventory.  │           │wim.order.      │           ││
│              │events          │           │events          │           ││
│              │wim.alerts      │           │wim.stock.      │           ││
│              │wim.audit.events│           │movements       │           ││
│              └────────────────┘           └────────────────┘           ││
│                                                      │                ││
│                       ┌──────────────────────────────┘                ││
│                       ▼                                             ││
│              ┌────────────────┐                                   ││
│              │   Consumer     │                                   ││
│              │   (Worker)     │                                   ││
│              └───────┬────────┘                                   ││
│                      │                                             ││
│       ┌──────────────┼────────────────┐                          ││
│       ▼              ▼                ▼                          ││
│  ┌────────┐    ┌────────┐    ┌────────┐                          ││
│  │Handler │    │Handler │    │Handler │                          ││
│  │Inv     │    │Order   │    │Alert   │                          ││
│  └────────┘    └────────┘    └────────┘                          ││
│                                                                     ││
└─────────────────────────────────────────────────────────────────────┘
```

## Topics

| Topic | Purpose | Partitions | Retention | Cleanup |
|-------|---------|------------|-----------|---------|
| `wim.inventory.events` | Product/inventory changes | 6 | 7 days | delete |
| `wim.order.events` | Order lifecycle events | 6 | 7 days | delete |
| `wim.stock.movements` | Stock movement records | 6 | 30 days | delete |
| `wim.alerts` | Expiry and system alerts | 3 | 1 day | delete |
| `wim.audit.events` | Audit/compliance events | 6 | 90 days | compact |
| `wim.system.events` | System events | 3 | 3 days | delete |
| `wim.dlq` | Dead letter queue | 3 | 30 days | delete |

## Event Schema

All events follow a standardized envelope format:

```json
{
  "version": "1.0",
  "id": "uuid-v4",
  "type": "event.type.name",
  "source": "wim-api",
  "timestamp": "2024-01-15T10:30:00Z",
  "correlation_id": "optional-trace-id",
  "payload": { ... },
  "metadata": { "key": "value" }
}
```

## Producer Features

### Configuration

```yaml
kafka:
  brokers:
    - localhost:9092
  topic: warehouse-events
  group_id: warehouse-worker
  
  # Producer settings
  producer_async: true      # Enable async publishing
  batch_size: 100           # Messages per batch
  batch_timeout: 100ms      # Max time to wait for batch
  required_acks: 1          # 0=none, 1=leader, -1=all
  compression: snappy       # none, gzip, snappy, lz4
  idempotent: true          # Enable idempotent producer
  
  # Security (optional)
  enable_sasl: false
  sasl_mechanism: plain
  sasl_username: ""
  sasl_password: ""
  enable_tls: false
```

### Usage

```go
// Create producer from config
producer, err := event.InitKafka(cfg.Kafka)
if err != nil {
    log.Fatal(err)
}
defer producer.Close()

// Publish synchronously
err = producer.Publish(ctx, event.Event{
    ID:        uuid.New().String(),
    Type:      event.EventOrderCreated,
    Payload:   payload,
    Timestamp: time.Now(),
})

// Or publish asynchronously
err = producer.PublishAsync(ctx, event)

// With callback
producer.PublishAsyncWithCallback(ctx, event, func(err error) {
    if err != nil {
        log.Printf("Failed to publish: %v", err)
    }
})

// Health check
if err := producer.HealthCheck(ctx); err != nil {
    // Handle connection failure
}

// Get metrics
metrics := producer.GetMetrics()
log.Printf("Sent: %d, Failed: %d", metrics.MessagesSent, metrics.MessagesFailed)
```

## Consumer Features

### Configuration

```go
config := event.ConsumerConfig{
    Brokers:           []string{"localhost:9092"},
    GroupID:           "warehouse-worker",
    Topics:            []string{event.TopicOrderEvents},
    MinBytes:          1024,           // Minimum fetch size
    MaxBytes:          10 * 1024 * 1024, // Maximum fetch size
    MaxWait:           500 * time.Millisecond,
    AutoCommit:        false,         // Manual commit for reliability
    CommitInterval:      5 * time.Second,
    MaxRetries:        3,
    RetryDelay:        time.Second,
}

consumer, err := event.NewConsumer(config)
if err != nil {
    log.Fatal(err)
}

// Register handlers
consumer.RegisterHandler(event.TopicOrderEvents, func(ctx context.Context, msg *event.ConsumedMessage) error {
    log.Printf("Received: %s", msg.Event.Type)
    // Process message
    return nil
})

// Start consuming
if err := consumer.Start(); err != nil {
    log.Fatal(err)
}
```

### Consumer Groups

Multiple consumers can join the same group for parallel processing:

```
Consumer 1 (Group: warehouse-worker) ─┐
                                       ├──▶ Topic: wim.order.events
Consumer 2 (Group: warehouse-worker) ─┤
                                       ├──▶ Partitions distributed
Consumer 3 (Group: warehouse-worker) ─┘
```

## Worker Pool

The worker pool provides concurrent message processing with retry and DLQ support:

```go
poolConfig := worker.PoolConfig{
    ConsumerConfig: event.ConsumerConfig{
        Brokers: []string{"localhost:9092"},
        GroupID: "warehouse-worker",
    },
    ProcessorConfig: queue.ProcessorConfig{
        Workers:      10,   // Concurrent workers
        QueueSize:    100,  // Pending queue size
        MaxRetries:   3,    // Retry attempts
        DLQEnabled:   true, // Enable dead letter queue
    },
    Topics: []string{
        event.TopicOrderEvents,
        event.TopicInventoryEvents,
    },
}

pool, err := worker.NewPool(poolConfig)
if err != nil {
    log.Fatal(err)
}

// Register handlers
pool.RegisterHandler(event.TopicOrderEvents, handleOrderEvent)
pool.RegisterHandler(event.TopicInventoryEvents, handleInventoryEvent)

// Start
if err := pool.Start(); err != nil {
    log.Fatal(err)
}

// Wait for shutdown signal
pool.Wait()
```

## Running the Worker

```bash
# Ensure Kafka is running
docker-compose up -d kafka

# Start the worker
go run cmd/worker/main.go
```

## Docker Compose

The included `docker-compose.yml` provides a complete Kafka setup:

```yaml
kafka:
  image: confluentinc/cp-kafka:7.5.0
  environment:
    KAFKA_NODE_ID: 1
    KAFKA_PROCESS_ROLES: broker,controller
    KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093
    KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
    KAFKA_CONTROLLER_QUORUM_VOTERS: 1@kafka:9093
    KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
```

## Monitoring

### Producer Metrics

- `MessagesSent`: Total messages successfully sent
- `MessagesFailed`: Total messages failed after retries
- `MessagesDropped`: Messages dropped (async queue full)
- `BytesSent`: Total bytes transmitted

### Consumer Metrics

- `MessagesConsumed`: Total messages consumed
- `MessagesFailed`: Messages that failed processing
- `MessagesRetried`: Messages retried due to errors
- `MessagesDLQ`: Messages sent to dead letter queue
- `CommitsSucceeded`: Successful offset commits
- `CommitsFailed`: Failed offset commits

### Health Checks

Both producer and consumer support health check endpoints:

```go
// Producer health
if err := producer.HealthCheck(ctx); err != nil {
    // Connection failed
}

// Consumer health
if err := consumer.HealthCheck(ctx); err != nil {
    // Connection failed
}
```

## Security

### SASL Authentication

```yaml
kafka:
  enable_sasl: true
  sasl_mechanism: plain
  sasl_username: "wim-user"
  sasl_password: "secure-password"
```

### TLS Encryption

```yaml
kafka:
  enable_tls: true
  tls_ca_file: "/path/to/ca.crt"
  tls_cert_file: "/path/to/client.crt"
  tls_key_file: "/path/to/client.key"
```

## Migration from Old System

The old `KafkaPublisher` is still available but deprecated:

```go
// Old way (deprecated)
publisher := event.NewKafkaPublisher(brokers, topic)

// New way (recommended)
producer, err := event.NewProducer(event.ProducerConfig{
    Brokers: brokers,
    DefaultTopic: topic,
    Async: true,
    // ... other options
})
```

## Best Practices

1. **Use async publishing** for high-throughput scenarios
2. **Enable idempotent producer** for exactly-once semantics
3. **Set appropriate batch sizes** to balance latency vs throughput
4. **Use manual commits** for critical events (auto-commit for logs)
5. **Implement proper error handling** in message handlers
6. **Monitor metrics** to detect issues early
7. **Configure DLQ** for failed messages
8. **Use correlation IDs** for request tracing across services

## Troubleshooting

### Connection Issues

```bash
# Test Kafka connectivity
nc -zv localhost 9092

# List topics
kafka-topics.sh --bootstrap-server localhost:9092 --list
```

### Performance Tuning

- **Increase batch_size** for higher throughput (trades latency)
- **Enable compression** (snappy recommended)
- **Adjust fetch.min.bytes** on consumer side
- **Monitor consumer lag** with `kafka-consumer-groups.sh`

## API Integration

The API server now automatically uses Kafka if configured:

```go
// In cmd/api/main.go
eventPublisher, err := event.InitKafka(cfg.Kafka)
if err != nil {
    log.Printf("Failed to initialize Kafka, using mock: %v", err)
    eventPublisher = event.NewMockPublisher()
}
```

Services automatically publish events:

```go
// Service layer
s.events.Publish(ctx, event.Event{
    Type:      event.EventOrderCreated,
    ID:        order.ID.String(),
    Payload:   eventData,
    Timestamp: time.Now(),
})
```
