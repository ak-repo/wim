package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// ConsumerConfig holds configuration for the Kafka consumer
type ConsumerConfig struct {
	Brokers           []string
	GroupID           string
	Topics            []string
	MinBytes          int
	MaxBytes          int
	MaxWait           time.Duration
	HeartbeatInterval time.Duration
	SessionTimeout    time.Duration
	RebalanceTimeout  time.Duration
	RetentionTime     time.Duration
	TLS               *TLSConfig
	SASL              *SASLConfig
	AutoCommit        bool
	CommitInterval    time.Duration
	MaxRetries        int
	RetryDelay        time.Duration
}

// MessageHandler processes consumed messages
type MessageHandler func(ctx context.Context, msg *ConsumedMessage) error

// ConsumedMessage represents a consumed Kafka message
type ConsumedMessage struct {
	Topic      string
	Partition  int
	Offset     int64
	Key        []byte
	Value      []byte
	Headers    map[string]string
	Timestamp  time.Time
	Event      EventSchema
	RawMessage kafka.Message
}

// HandlerRegistration represents a handler registration
type HandlerRegistration struct {
	Topic      string
	Handler    MessageHandler
	RetryCount int
	RetryDelay time.Duration
}

// Consumer is a production-grade Kafka consumer with consumer group support
type Consumer struct {
	config   ConsumerConfig
	handlers map[string][]HandlerRegistration
	mu       sync.RWMutex
	runners  []*consumerRunner
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	closed   bool
	metrics  *ConsumerMetrics
	ready    chan struct{}
}

// ConsumerMetrics holds consumer metrics
type ConsumerMetrics struct {
	mu               sync.RWMutex
	MessagesConsumed int64
	MessagesFailed   int64
	MessagesRetried  int64
	MessagesDLQ      int64
	CommitsSucceeded int64
	CommitsFailed    int64
}

type consumerRunner struct {
	reader   *kafka.Reader
	handlers []HandlerRegistration
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewConsumer creates a new production-grade Kafka consumer
func NewConsumer(config ConsumerConfig) (*Consumer, error) {
	if len(config.Brokers) == 0 {
		return nil, errors.New("at least one broker required")
	}
	if config.GroupID == "" {
		return nil, errors.New("group ID is required")
	}

	// Set defaults
	if config.MinBytes <= 0 {
		config.MinBytes = 1024
	}
	if config.MaxBytes <= 0 {
		config.MaxBytes = 10 * 1024 * 1024 // 10MB
	}
	if config.MaxWait <= 0 {
		config.MaxWait = 500 * time.Millisecond
	}
	if config.HeartbeatInterval <= 0 {
		config.HeartbeatInterval = 3 * time.Second
	}
	if config.SessionTimeout <= 0 {
		config.SessionTimeout = 30 * time.Second
	}
	if config.CommitInterval <= 0 {
		config.CommitInterval = 5 * time.Second
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = 1 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &Consumer{
		config:   config,
		handlers: make(map[string][]HandlerRegistration),
		ctx:      ctx,
		cancel:   cancel,
		metrics:  &ConsumerMetrics{},
		ready:    make(chan struct{}),
	}

	return c, nil
}

// RegisterHandler registers a handler for a specific topic
func (c *Consumer) RegisterHandler(topic string, handler MessageHandler) {
	c.RegisterHandlerWithRetry(topic, handler, 0, 0)
}

// RegisterHandlerWithRetry registers a handler with retry configuration
func (c *Consumer) RegisterHandlerWithRetry(topic string, handler MessageHandler, retryCount int, retryDelay time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	reg := HandlerRegistration{
		Topic:      topic,
		Handler:    handler,
		RetryCount: retryCount,
		RetryDelay: retryDelay,
	}

	c.handlers[topic] = append(c.handlers[topic], reg)
}

// createDialer creates a kafka.Dialer with TLS and SASL configuration
func (c *Consumer) createDialer() (*kafka.Dialer, error) {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	if c.config.SASL != nil {
		switch c.config.SASL.Mechanism {
		case "plain":
			dialer.SASLMechanism = plain.Mechanism{
				Username: c.config.SASL.Username,
				Password: c.config.SASL.Password,
			}
		default:
			return nil, fmt.Errorf("unsupported SASL mechanism: %s", c.config.SASL.Mechanism)
		}
	}

	return dialer, nil
}

// Start begins consuming messages from all registered topics
func (c *Consumer) Start() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return errors.New("consumer is closed")
	}

	dialer, err := c.createDialer()
	if err != nil {
		return err
	}

	// Create a reader for each unique topic
	for topic := range c.handlers {
		readerConfig := kafka.ReaderConfig{
			Brokers:           c.config.Brokers,
			GroupID:           c.config.GroupID,
			Topic:             topic,
			MinBytes:          c.config.MinBytes,
			MaxBytes:          c.config.MaxBytes,
			MaxWait:           c.config.MaxWait,
			HeartbeatInterval: c.config.HeartbeatInterval,
			SessionTimeout:    c.config.SessionTimeout,
			RebalanceTimeout:  c.config.RebalanceTimeout,
			RetentionTime:     c.config.RetentionTime,
			StartOffset:       kafka.LastOffset,
			ReadBackoffMin:    100 * time.Millisecond,
			ReadBackoffMax:    1 * time.Second,
			Logger:            nil, // Can add logging here
			ErrorLogger:       nil, // Can add error logging here
			Dialer:            dialer,
		}

		if !c.config.AutoCommit {
			readerConfig.CommitInterval = 0 // Disable auto commit
		} else {
			readerConfig.CommitInterval = c.config.CommitInterval
		}

		reader := kafka.NewReader(readerConfig)

		runnerCtx, runnerCancel := context.WithCancel(c.ctx)
		runner := &consumerRunner{
			reader:   reader,
			handlers: c.handlers[topic],
			ctx:      runnerCtx,
			cancel:   runnerCancel,
		}

		c.runners = append(c.runners, runner)

		c.wg.Add(1)
		go c.runConsumer(runner, topic)
	}

	close(c.ready)
	return nil
}

// WaitForReady blocks until the consumer is ready
func (c *Consumer) WaitForReady(timeout time.Duration) bool {
	select {
	case <-c.ready:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (c *Consumer) runConsumer(runner *consumerRunner, topic string) {
	defer c.wg.Done()

	for {
		select {
		case <-runner.ctx.Done():
			return
		default:
		}

		msg, err := runner.reader.ReadMessage(runner.ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			// Log error and continue
			fmt.Printf("Error reading message from topic %s: %v\n", topic, err)
			continue
		}

		c.metrics.mu.Lock()
		c.metrics.MessagesConsumed++
		c.metrics.mu.Unlock()

		// Parse the consumed message
		consumedMsg, err := c.parseMessage(msg)
		if err != nil {
			// Log error but commit the message to avoid infinite loop
			fmt.Printf("Error parsing message: %v\n", err)
			if !c.config.AutoCommit {
				_ = runner.reader.CommitMessages(runner.ctx, msg)
			}
			continue
		}

		// Process with handlers
		for _, reg := range runner.handlers {
			err := c.processWithRetry(runner.ctx, reg, consumedMsg)
			if err != nil {
				c.metrics.mu.Lock()
				c.metrics.MessagesFailed++
				c.metrics.mu.Unlock()
				// Could send to DLQ here
			}
		}

		// Manual commit if auto commit is disabled
		if !c.config.AutoCommit {
			err := runner.reader.CommitMessages(runner.ctx, msg)
			if err != nil {
				c.metrics.mu.Lock()
				c.metrics.CommitsFailed++
				c.metrics.mu.Unlock()
			} else {
				c.metrics.mu.Lock()
				c.metrics.CommitsSucceeded++
				c.metrics.mu.Unlock()
			}
		}
	}
}

func (c *Consumer) parseMessage(msg kafka.Message) (*ConsumedMessage, error) {
	// Parse headers
	headers := make(map[string]string)
	for _, h := range msg.Headers {
		headers[h.Key] = string(h.Value)
	}

	// Parse event schema from value
	var schema EventSchema
	if err := json.Unmarshal(msg.Value, &schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	return &ConsumedMessage{
		Topic:      msg.Topic,
		Partition:  msg.Partition,
		Offset:     msg.Offset,
		Key:        msg.Key,
		Value:      msg.Value,
		Headers:    headers,
		Timestamp:  msg.Time,
		Event:      schema,
		RawMessage: msg,
	}, nil
}

func (c *Consumer) processWithRetry(ctx context.Context, reg HandlerRegistration, msg *ConsumedMessage) error {
	var lastErr error

	for attempt := 0; attempt <= reg.RetryCount; attempt++ {
		if attempt > 0 {
			c.metrics.mu.Lock()
			c.metrics.MessagesRetried++
			c.metrics.mu.Unlock()
			time.Sleep(reg.RetryDelay * time.Duration(attempt))
		}

		err := reg.Handler(ctx, msg)
		if err == nil {
			return nil
		}

		lastErr = err
	}

	return fmt.Errorf("failed after %d retries: %w", reg.RetryCount, lastErr)
}

// GetMetrics returns current consumer metrics
func (c *Consumer) GetMetrics() ConsumerMetrics {
	c.metrics.mu.RLock()
	defer c.metrics.mu.RUnlock()
	return ConsumerMetrics{
		MessagesConsumed: c.metrics.MessagesConsumed,
		MessagesFailed:   c.metrics.MessagesFailed,
		MessagesRetried:  c.metrics.MessagesRetried,
		MessagesDLQ:      c.metrics.MessagesDLQ,
		CommitsSucceeded: c.metrics.CommitsSucceeded,
		CommitsFailed:    c.metrics.CommitsFailed,
	}
}

// HealthCheck checks if the consumer can connect to Kafka
func (c *Consumer) HealthCheck(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", c.config.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()
	return nil
}

// Close gracefully shuts down the consumer
func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	c.cancel()

	// Cancel all runners
	for _, runner := range c.runners {
		runner.cancel()
	}

	// Wait for all runners to finish
	c.wg.Wait()

	// Close all readers
	for _, runner := range c.runners {
		if err := runner.reader.Close(); err != nil {
			// Log error
			fmt.Printf("Error closing reader: %v\n", err)
		}
	}

	return nil
}

// IsHealthy returns true if the consumer is running
func (c *Consumer) IsHealthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return !c.closed
}
