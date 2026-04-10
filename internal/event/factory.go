package event

import (
	"time"

	"github.com/ak-repo/wim/config"
	"github.com/segmentio/kafka-go"
)

// NewProducerFromConfig creates a new Producer from config
func NewProducerFromConfig(cfg config.KafkaConfig) (*Producer, error) {
	// Parse compression
	var compression kafka.Compression
	switch cfg.Compression {
	case "gzip":
		compression = kafka.Gzip
	case "snappy":
		compression = kafka.Snappy
	case "lz4":
		compression = kafka.Lz4
	default:
		compression = kafka.Compression(0) // None
	}

	// Parse required acks
	var requiredAcks kafka.RequiredAcks
	switch cfg.RequiredAcks {
	case 0:
		requiredAcks = kafka.RequireNone
	case 1:
		requiredAcks = kafka.RequireOne
	case -1:
		requiredAcks = kafka.RequireAll
	default:
		requiredAcks = kafka.RequireOne
	}

	// Build SASL config
	var saslConfig *SASLConfig
	if cfg.EnableSASL {
		saslConfig = &SASLConfig{
			Mechanism: cfg.SASLMechanism,
			Username:  cfg.SASLUsername,
			Password:  cfg.SASLPassword,
		}
	}

	// Build TLS config
	var tlsConfig *TLSConfig
	if cfg.EnableTLS {
		tlsConfig = &TLSConfig{
			Enabled:  true,
			CAFile:   cfg.TLSCAFile,
			CertFile: cfg.TLSCertFile,
			KeyFile:  cfg.TLSKeyFile,
		}
	}

	producerConfig := ProducerConfig{
		Brokers:         cfg.Brokers,
		DefaultTopic:    cfg.Topic,
		BatchSize:       cfg.BatchSize,
		BatchTimeout:    cfg.BatchTimeout,
		RequiredAcks:    requiredAcks,
		Compression:     compression,
		MaxRetries:      3,
		RetryBackoff:    100 * time.Millisecond,
		TLS:             tlsConfig,
		SASL:            saslConfig,
		Idempotent:      cfg.Idempotent,
		Async:           cfg.ProducerAsync,
		QueueBufferSize: 1000,
	}

	return NewProducer(producerConfig)
}

// NewConsumerFromConfig creates a new Consumer from config
func NewConsumerFromConfig(cfg config.KafkaConfig) (*Consumer, error) {
	// Build SASL config
	var saslConfig *SASLConfig
	if cfg.EnableSASL {
		saslConfig = &SASLConfig{
			Mechanism: cfg.SASLMechanism,
			Username:  cfg.SASLUsername,
			Password:  cfg.SASLPassword,
		}
	}

	// Build TLS config
	var tlsConfig *TLSConfig
	if cfg.EnableTLS {
		tlsConfig = &TLSConfig{
			Enabled:  true,
			CAFile:   cfg.TLSCAFile,
			CertFile: cfg.TLSCertFile,
			KeyFile:  cfg.TLSKeyFile,
		}
	}

	consumerConfig := ConsumerConfig{
		Brokers:        cfg.Brokers,
		GroupID:        cfg.GroupID,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		MaxWait:        cfg.MaxWait,
		AutoCommit:     cfg.AutoCommit,
		CommitInterval: cfg.CommitInterval,
		TLS:            tlsConfig,
		SASL:           saslConfig,
		MaxRetries:     3,
		RetryDelay:     time.Second,
	}

	return NewConsumer(consumerConfig)
}

// InitKafka initializes Kafka topics and returns a producer
func InitKafka(cfg config.KafkaConfig) (*Producer, error) {
	// Create topic manager and ensure topics exist
	tm := NewTopicManager(cfg.Brokers)

	if err := tm.EnsureAllTopics(nil); err != nil {
		// Log but don't fail - topics might already exist
		// or broker might not be available yet
	}

	// Create producer
	return NewProducerFromConfig(cfg)
}
