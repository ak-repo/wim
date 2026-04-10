package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ak-repo/wim/internal/event"
	"github.com/ak-repo/wim/internal/queue"
)

// Pool represents a worker pool for processing Kafka messages
type Pool struct {
	config     PoolConfig
	consumer   *event.Consumer
	processor  *queue.Processor
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	handlers   map[string]Handler
	mu         sync.RWMutex
	shutdownCh chan struct{}
	ready      bool
}

// PoolConfig holds configuration for the worker pool
type PoolConfig struct {
	ConsumerConfig  event.ConsumerConfig
	ProcessorConfig queue.ProcessorConfig
	Topics          []string
	GroupID         string
}

// Handler is a function that processes a message
type Handler func(ctx context.Context, msg *event.ConsumedMessage) error

// NewPool creates a new worker pool
func NewPool(config PoolConfig) (*Pool, error) {
	if config.GroupID == "" {
		return nil, fmt.Errorf("group ID is required")
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create logger for processor
	logger := &queue.DefaultLogger{}

	// Create message processor
	processor, err := queue.NewProcessor(config.ProcessorConfig, logger)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create processor: %w", err)
	}

	pool := &Pool{
		config:     config,
		processor:  processor,
		ctx:        ctx,
		cancel:     cancel,
		handlers:   make(map[string]Handler),
		shutdownCh: make(chan struct{}),
	}

	return pool, nil
}

// RegisterHandler registers a handler for a specific topic
func (p *Pool) RegisterHandler(topic string, handler Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[topic] = handler
}

// Start starts the worker pool
func (p *Pool) Start() error {
	if p.ready {
		return fmt.Errorf("pool already started")
	}

	// Start the message processor
	p.processor.Start()

	// Create and start the consumer
	consumer, err := event.NewConsumer(p.config.ConsumerConfig)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	p.consumer = consumer

	// Register handlers for each topic
	for _, topic := range p.config.Topics {
		if handler, ok := p.handlers[topic]; ok {
			consumer.RegisterHandler(topic, p.wrapHandler(handler))
		}
	}

	// Start consuming
	if err := consumer.Start(); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	p.ready = true

	log.Printf("Worker pool started with %d workers", p.config.ProcessorConfig.Workers)

	return nil
}

func (p *Pool) wrapHandler(handler Handler) event.MessageHandler {
	return func(ctx context.Context, msg *event.ConsumedMessage) error {
		return handler(ctx, msg)
	}
}

// Wait blocks until the pool is shutdown
func (p *Pool) Wait() {
	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Println("Received shutdown signal")
	case <-p.ctx.Done():
	}

	p.Shutdown()
}

// Shutdown gracefully stops the worker pool
func (p *Pool) Shutdown() error {
	if !p.ready {
		return nil
	}

	p.ready = false
	p.cancel()

	// Shutdown consumer
	if p.consumer != nil {
		if err := p.consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}

	// Shutdown processor
	if p.processor != nil {
		if err := p.processor.Stop(); err != nil {
			log.Printf("Error stopping processor: %v", err)
		}
	}

	close(p.shutdownCh)
	log.Println("Worker pool shutdown complete")

	return nil
}

// HealthCheck checks if the pool is healthy
func (p *Pool) HealthCheck() error {
	if !p.ready {
		return fmt.Errorf("pool not ready")
	}

	if p.consumer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := p.consumer.HealthCheck(ctx); err != nil {
			return fmt.Errorf("consumer health check failed: %w", err)
		}
	}

	return nil
}

// IsReady returns true if the pool is ready
func (p *Pool) IsReady() bool {
	return p.ready
}

// GetMetrics returns metrics from the processor
func (p *Pool) GetMetrics() queue.ProcessorMetrics {
	if p.processor != nil {
		return p.processor.GetMetrics()
	}
	return queue.ProcessorMetrics{}
}

// GetConsumerMetrics returns consumer metrics
func (p *Pool) GetConsumerMetrics() event.ConsumerMetrics {
	if p.consumer != nil {
		return p.consumer.GetMetrics()
	}
	return event.ConsumerMetrics{}
}

// PoolSize returns the configured number of workers
func (p *Pool) PoolSize() int {
	return p.config.ProcessorConfig.Workers
}

// QueueSize returns the current queue size
func (p *Pool) QueueSize() (pending int, capacity int) {
	if p.processor != nil {
		return p.processor.GetQueueStats()
	}
	return 0, 0
}
