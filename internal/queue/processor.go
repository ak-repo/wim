package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	eventpkg "github.com/ak-repo/wim/internal/event"
)

// ProcessorConfig holds configuration for the message processor
type ProcessorConfig struct {
	Workers       int
	QueueSize     int
	BatchSize     int
	BatchTimeout  time.Duration
	MaxRetries    int
	RetryDelay    time.Duration
	DLQEnabled    bool
	DLQTopic      string
	DeadLetterTTL time.Duration
}

// Processor handles message processing with worker pool
type Processor struct {
	config   ProcessorConfig
	handlers map[string]HandlerFunc
	jobs     chan *Job
	results  chan *Result
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	closed   bool
	mu       sync.RWMutex
	logger   Logger
	metrics  *ProcessorMetrics
}

// Logger interface for logging
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// DefaultLogger uses standard log
type DefaultLogger struct{}

func (l *DefaultLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf("[INFO] %s %v", msg, formatKV(keysAndValues))
}

func (l *DefaultLogger) Debug(msg string, keysAndValues ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, formatKV(keysAndValues))
}

func (l *DefaultLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, formatKV(keysAndValues))
}

func formatKV(kv []interface{}) string {
	if len(kv)%2 != 0 {
		kv = append(kv, "MISSING")
	}
	var parts []string
	for i := 0; i < len(kv); i += 2 {
		parts = append(parts, fmt.Sprintf("%s=%v", kv[i], kv[i+1]))
	}
	return strings.Join(parts, " ")
}

// Job represents a message to be processed
type Job struct {
	ID         string
	Topic      string
	Event      eventpkg.EventSchema
	RetryCount int
	CreatedAt  time.Time
}

// Result represents the processing result
type Result struct {
	Job       *Job
	Success   bool
	Error     error
	Retryable bool
	Duration  time.Duration
}

// HandlerFunc processes a specific event type
type HandlerFunc func(ctx context.Context, event eventpkg.EventSchema) error

// ProcessorMetrics holds processor metrics
type ProcessorMetrics struct {
	mu            sync.RWMutex
	JobsProcessed int64
	JobsFailed    int64
	JobsRetried   int64
	JobsDLQ       int64
	AvgDuration   time.Duration
}

// NewProcessor creates a new message processor
func NewProcessor(config ProcessorConfig, logger Logger) (*Processor, error) {
	// Set defaults
	if config.Workers <= 0 {
		config.Workers = 5
	}
	if config.QueueSize <= 0 {
		config.QueueSize = 100
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 10
	}
	if config.BatchTimeout <= 0 {
		config.BatchTimeout = 1 * time.Second
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = 1 * time.Second
	}
	if config.DLQTopic == "" {
		config.DLQTopic = eventpkg.TopicDLQ
	}

	ctx, cancel := context.WithCancel(context.Background())

	if logger == nil {
		logger = &DefaultLogger{}
	}

	p := &Processor{
		config:   config,
		handlers: make(map[string]HandlerFunc),
		jobs:     make(chan *Job, config.QueueSize),
		results:  make(chan *Result, config.QueueSize),
		ctx:      ctx,
		cancel:   cancel,
		logger:   logger,
		metrics:  &ProcessorMetrics{},
	}

	return p, nil
}

// RegisterHandler registers a handler for a specific event type
func (p *Processor) RegisterHandler(eventType eventpkg.EventType, handler HandlerFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[string(eventType)] = handler
}

// Start begins processing messages
func (p *Processor) Start() {
	// Start workers
	for i := 0; i < p.config.Workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	// Start result processor
	p.wg.Add(1)
	go p.resultProcessor()
}

// Stop gracefully shuts down the processor
func (p *Processor) Stop() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	p.cancel()

	// Wait for all workers to finish
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout waiting for workers to stop")
	}
}

// Submit adds a job to the processing queue
func (p *Processor) Submit(event eventpkg.EventSchema) error {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return fmt.Errorf("processor is closed")
	}
	p.mu.RUnlock()

	job := &Job{
		ID:         event.ID,
		Topic:      eventpkg.GetTopicForEvent(event.Type),
		Event:      event,
		RetryCount: 0,
		CreatedAt:  time.Now(),
	}

	select {
	case p.jobs <- job:
		return nil
	case <-p.ctx.Done():
		return p.ctx.Err()
	default:
		return fmt.Errorf("job queue is full")
	}
}

// SubmitAsync adds a job asynchronously (non-blocking)
func (p *Processor) SubmitAsync(evt eventpkg.EventSchema) bool {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return false
	}
	p.mu.RUnlock()

	job := &Job{
		ID:         evt.ID,
		Topic:      eventpkg.GetTopicForEvent(evt.Type),
		Event:      evt,
		RetryCount: 0,
		CreatedAt:  time.Now(),
	}

	select {
	case p.jobs <- job:
		return true
	default:
		return false
	}
}

func (p *Processor) worker(id int) {
	defer p.wg.Done()

	p.logger.Info("Worker started", "worker_id", id)

	for {
		select {
		case job := <-p.jobs:
			if job == nil {
				return
			}
			result := p.processJob(p.ctx, job)
			p.results <- result

		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Processor) processJob(ctx context.Context, job *Job) *Result {
	start := time.Now()

	p.mu.RLock()
	handler, exists := p.handlers[string(job.Event.Type)]
	p.mu.RUnlock()

	if !exists {
		return &Result{
			Job:       job,
			Success:   false,
			Error:     fmt.Errorf("no handler registered for event type: %s", job.Event.Type),
			Retryable: false,
			Duration:  time.Since(start),
		}
	}

	err := handler(ctx, job.Event)
	duration := time.Since(start)

	if err == nil {
		return &Result{
			Job:      job,
			Success:  true,
			Duration: duration,
		}
	}

	// Determine if error is retryable
	retryable := p.isRetryableError(err)

	return &Result{
		Job:       job,
		Success:   false,
		Error:     err,
		Retryable: retryable,
		Duration:  duration,
	}
}

func (p *Processor) resultProcessor() {
	defer p.wg.Done()

	for {
		select {
		case result := <-p.results:
			if result == nil {
				return
			}
			p.handleResult(result)

		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Processor) handleResult(result *Result) {
	p.metrics.mu.Lock()
	p.metrics.JobsProcessed++
	if result.Success {
		p.logger.Debug("Job processed successfully",
			"job_id", result.Job.ID,
			"event_type", result.Job.Event.Type,
			"duration", result.Duration,
		)
	} else {
		p.metrics.JobsFailed++
		p.logger.Error("Job processing failed",
			"job_id", result.Job.ID,
			"event_type", result.Job.Event.Type,
			"error", result.Error,
			"retryable", result.Retryable,
		)

		if result.Retryable && result.Job.RetryCount < p.config.MaxRetries {
			p.retryJob(result.Job)
		} else if p.config.DLQEnabled {
			p.sendToDLQ(result.Job, result.Error)
		}
	}
	p.metrics.mu.Unlock()
}

func (p *Processor) retryJob(job *Job) {
	job.RetryCount++

	p.metrics.JobsRetried++
	p.logger.Info("Retrying job",
		"job_id", job.ID,
		"retry_count", job.RetryCount,
		"max_retries", p.config.MaxRetries,
	)

	// Delay before retry
	time.Sleep(p.config.RetryDelay * time.Duration(job.RetryCount))

	// Re-submit the job
	select {
	case p.jobs <- job:
	case <-p.ctx.Done():
	}
}

func (p *Processor) sendToDLQ(job *Job, processingErr error) {
	p.metrics.JobsDLQ++

	p.logger.Error("Sending job to DLQ",
		"job_id", job.ID,
		"event_type", job.Event.Type,
		"error", processingErr,
	)

	// In a real implementation, this would publish to a DLQ topic
	// For now, we just log it
	dlqMessage := map[string]interface{}{
		"original_event": job.Event,
		"error":          processingErr.Error(),
		"retry_count":    job.RetryCount,
		"timestamp":      time.Now(),
	}

	data, _ := json.Marshal(dlqMessage)
	p.logger.Error("DLQ message", "data", string(data))
}

func (p *Processor) isRetryableError(err error) bool {
	// List of non-retryable errors
	nonRetryableErrors := []string{
		"invalid",
		"not found",
		"already exists",
		"unauthorized",
		"forbidden",
	}

	errStr := err.Error()
	for _, nr := range nonRetryableErrors {
		if contains(errStr, nr) {
			return false
		}
	}

	return true
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetMetrics returns processor metrics
func (p *Processor) GetMetrics() ProcessorMetrics {
	p.metrics.mu.RLock()
	defer p.metrics.mu.RUnlock()
	return ProcessorMetrics{
		JobsProcessed: p.metrics.JobsProcessed,
		JobsFailed:    p.metrics.JobsFailed,
		JobsRetried:   p.metrics.JobsRetried,
		JobsDLQ:       p.metrics.JobsDLQ,
	}
}

// GetQueueStats returns current queue statistics
func (p *Processor) GetQueueStats() (pending int, capacity int) {
	return len(p.jobs), cap(p.jobs)
}

// IsHealthy returns true if the processor is running
func (p *Processor) IsHealthy() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return !p.closed
}
