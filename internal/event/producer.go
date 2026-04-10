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

// ProducerConfig holds configuration for the Kafka producer
type ProducerConfig struct {
	Brokers         []string
	DefaultTopic    string
	BatchSize       int
	BatchTimeout    time.Duration
	RequiredAcks    kafka.RequiredAcks
	Compression     kafka.Compression
	MaxRetries      int
	RetryBackoff    time.Duration
	TLS             *TLSConfig
	SASL            *SASLConfig
	Idempotent      bool
	Async           bool
	QueueBufferSize int
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled  bool
	CAFile   string
	CertFile string
	KeyFile  string
}

// SASLConfig holds SASL configuration
type SASLConfig struct {
	Mechanism string // "plain", "scram-sha-256", "scram-sha-512"
	Username  string
	Password  string
}

// Producer is a production-grade Kafka producer
type Producer struct {
	config  ProducerConfig
	writers map[string]*kafka.Writer
	mu      sync.RWMutex
	asyncCh chan *asyncMessage
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	closed  bool
	metrics *ProducerMetrics
}

// ProducerMetrics holds producer metrics
type ProducerMetrics struct {
	mu              sync.RWMutex
	MessagesSent    int64
	MessagesFailed  int64
	MessagesDropped int64
	BytesSent       int64
}

type asyncMessage struct {
	ctx     context.Context
	topic   string
	key     []byte
	value   []byte
	headers map[string]string
	respCh  chan error
}

// NewProducer creates a new production-grade Kafka producer
func NewProducer(config ProducerConfig) (*Producer, error) {
	if len(config.Brokers) == 0 {
		return nil, errors.New("at least one broker required")
	}

	// Set defaults
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}
	if config.BatchTimeout <= 0 {
		config.BatchTimeout = 100 * time.Millisecond
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	if config.RetryBackoff <= 0 {
		config.RetryBackoff = 100 * time.Millisecond
	}
	if config.QueueBufferSize <= 0 {
		config.QueueBufferSize = 1000
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &Producer{
		config:  config,
		writers: make(map[string]*kafka.Writer),
		asyncCh: make(chan *asyncMessage, config.QueueBufferSize),
		ctx:     ctx,
		cancel:  cancel,
		metrics: &ProducerMetrics{},
	}

	// Start async worker if async mode is enabled
	if config.Async {
		p.wg.Add(1)
		go p.asyncWorker()
	}

	return p, nil
}

// createDialer creates a kafka.Dialer with TLS and SASL configuration
func (p *Producer) createDialer() (*kafka.Dialer, error) {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	if p.config.SASL != nil {
		switch p.config.SASL.Mechanism {
		case "plain":
			dialer.SASLMechanism = plain.Mechanism{
				Username: p.config.SASL.Username,
				Password: p.config.SASL.Password,
			}
		default:
			return nil, fmt.Errorf("unsupported SASL mechanism: %s", p.config.SASL.Mechanism)
		}
	}

	return dialer, nil
}

// getWriter gets or creates a writer for a specific topic
func (p *Producer) getWriter(topic string) (*kafka.Writer, error) {
	p.mu.RLock()
	writer, exists := p.writers[topic]
	p.mu.RUnlock()

	if exists {
		return writer, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring lock
	if writer, exists := p.writers[topic]; exists {
		return writer, nil
	}

	dialer, err := p.createDialer()
	if err != nil {
		return nil, err
	}

	writer = &kafka.Writer{
		Addr:         kafka.TCP(p.config.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		BatchSize:    p.config.BatchSize,
		BatchTimeout: p.config.BatchTimeout,
		RequiredAcks: p.config.RequiredAcks,
		Compression:  p.config.Compression,
		MaxAttempts:  p.config.MaxRetries,
		Transport:    &kafka.Transport{Dial: dialer.DialFunc},
	}

	if p.config.Idempotent {
		writer.RequiredAcks = kafka.RequireAll
	}

	p.writers[topic] = writer
	return writer, nil
}

// Publish sends an event synchronously
func (p *Producer) Publish(ctx context.Context, event Event) error {
	if p.closed {
		return errors.New("producer is closed")
	}

	// Build event schema
	schema := EventSchema{
		Version:   EventVersion,
		ID:        event.ID,
		Type:      event.Type,
		Source:    EventSource,
		Timestamp: event.Timestamp,
		Payload:   event.Payload,
		Metadata:  make(map[string]string),
	}

	// Get topic for event
	topic := GetTopicForEvent(event.Type)
	if topic == "" {
		topic = p.config.DefaultTopic
	}

	payload, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Try to publish with retries
	return p.publishWithRetry(ctx, topic, []byte(event.ID), payload, nil)
}

// PublishAsync sends an event asynchronously (returns immediately)
func (p *Producer) PublishAsync(ctx context.Context, event Event) error {
	if p.closed {
		return errors.New("producer is closed")
	}

	if !p.config.Async {
		// Fall back to sync publish
		return p.Publish(ctx, event)
	}

	schema := EventSchema{
		Version:   EventVersion,
		ID:        event.ID,
		Type:      event.Type,
		Source:    EventSource,
		Timestamp: event.Timestamp,
		Payload:   event.Payload,
		Metadata:  make(map[string]string),
	}

	topic := GetTopicForEvent(event.Type)
	if topic == "" {
		topic = p.config.DefaultTopic
	}

	payload, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &asyncMessage{
		ctx:    ctx,
		topic:  topic,
		key:    []byte(event.ID),
		value:  payload,
		respCh: nil, // Fire and forget for true async
	}

	select {
	case p.asyncCh <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Queue is full, drop the message and increment metrics
		p.metrics.mu.Lock()
		p.metrics.MessagesDropped++
		p.metrics.mu.Unlock()
		return errors.New("async queue is full, message dropped")
	}
}

// PublishAsyncWithCallback sends an event asynchronously with a callback
func (p *Producer) PublishAsyncWithCallback(ctx context.Context, event Event, callback func(error)) error {
	if p.closed {
		return errors.New("producer is closed")
	}

	if !p.config.Async {
		// Fall back to sync publish
		err := p.Publish(ctx, event)
		if callback != nil {
			callback(err)
		}
		return err
	}

	schema := EventSchema{
		Version:   EventVersion,
		ID:        event.ID,
		Type:      event.Type,
		Source:    EventSource,
		Timestamp: event.Timestamp,
		Payload:   event.Payload,
		Metadata:  make(map[string]string),
	}

	topic := GetTopicForEvent(event.Type)
	if topic == "" {
		topic = p.config.DefaultTopic
	}

	payload, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &asyncMessage{
		ctx:    ctx,
		topic:  topic,
		key:    []byte(event.ID),
		value:  payload,
		respCh: make(chan error, 1),
	}

	go func() {
		if callback != nil {
			callback(<-msg.respCh)
		}
		close(msg.respCh)
	}()

	select {
	case p.asyncCh <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		p.metrics.mu.Lock()
		p.metrics.MessagesDropped++
		p.metrics.mu.Unlock()
		return errors.New("async queue is full, message dropped")
	}
}

func (p *Producer) publishWithRetry(ctx context.Context, topic string, key, value []byte, headers map[string]string) error {
	writer, err := p.getWriter(topic)
	if err != nil {
		return err
	}

	kafkaHeaders := make([]kafka.Header, 0, len(headers))
	for k, v := range headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: k, Value: []byte(v)})
	}

	var lastErr error
	for attempt := 0; attempt < p.config.MaxRetries; attempt++ {
		err = writer.WriteMessages(ctx, kafka.Message{
			Key:     key,
			Value:   value,
			Headers: kafkaHeaders,
		})

		if err == nil {
			p.metrics.mu.Lock()
			p.metrics.MessagesSent++
			p.metrics.BytesSent += int64(len(key) + len(value))
			p.metrics.mu.Unlock()
			return nil
		}

		lastErr = err
		if attempt < p.config.MaxRetries-1 {
			time.Sleep(p.config.RetryBackoff * time.Duration(attempt+1))
		}
	}

	p.metrics.mu.Lock()
	p.metrics.MessagesFailed++
	p.metrics.mu.Unlock()

	return fmt.Errorf("failed after %d attempts: %w", p.config.MaxRetries, lastErr)
}

func (p *Producer) asyncWorker() {
	defer p.wg.Done()

	for {
		select {
		case msg := <-p.asyncCh:
			err := p.publishWithRetry(msg.ctx, msg.topic, msg.key, msg.value, msg.headers)
			if msg.respCh != nil {
				select {
				case msg.respCh <- err:
				case <-p.ctx.Done():
				}
			}
		case <-p.ctx.Done():
			return
		}
	}
}

// GetMetrics returns current producer metrics
func (p *Producer) GetMetrics() ProducerMetrics {
	p.metrics.mu.RLock()
	defer p.metrics.mu.RUnlock()
	return ProducerMetrics{
		MessagesSent:    p.metrics.MessagesSent,
		MessagesFailed:  p.metrics.MessagesFailed,
		MessagesDropped: p.metrics.MessagesDropped,
		BytesSent:       p.metrics.BytesSent,
	}
}

// HealthCheck checks if the producer can connect to Kafka
func (p *Producer) HealthCheck(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", p.config.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()
	return nil
}

// Close gracefully shuts down the producer
func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	p.cancel()

	// Wait for async worker to finish
	if p.config.Async {
		close(p.asyncCh)
		p.wg.Wait()
	}

	// Close all writers
	var errs []error
	for topic, writer := range p.writers {
		if err := writer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close writer for topic %s: %w", topic, err))
		}
	}

	if len(errs) > 0 {
		return errors.New("errors closing writers")
	}

	return nil
}

// Ensure Producer implements EventPublisher
var _ EventPublisher = (*Producer)(nil)

// Subscribe is not supported by the producer
func (p *Producer) Subscribe(topic string, handler EventHandler) error {
	return errors.New("subscribe is not supported by producer, use Consumer instead")
}
