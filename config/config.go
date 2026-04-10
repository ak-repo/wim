package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Auth     AuthConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	Worker   WorkerConfig
	LogLevel string
}

type ServerConfig struct {
	Port int
	Host string
}

type AuthConfig struct {
	JWTSecret       string
	JWTIssuer       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type DatabaseConfig struct {
	URL      string
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
	MaxConns int32

	ConnectRetries           int
	ConnectRetryInitialDelay time.Duration
	ConnectRetryMaxDelay     time.Duration
}

func (d DatabaseConfig) DSN() string {
	if d.URL != "" {
		return d.URL
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
		d.SSLMode,
	)
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
	// Producer settings
	ProducerAsync bool
	BatchSize     int
	BatchTimeout  time.Duration
	RequiredAcks  int    // 0, 1, -1 (all)
	Compression   string // "none", "gzip", "snappy", "lz4"
	Idempotent    bool
	// Consumer settings
	MinBytes       int
	MaxBytes       int
	MaxWait        time.Duration
	AutoCommit     bool
	CommitInterval time.Duration
	// Security settings
	EnableSASL    bool
	SASLMechanism string
	SASLUsername  string
	SASLPassword  string
	EnableTLS     bool
	TLSCAFile     string
	TLSCertFile   string
	TLSKeyFile    string
}

type WorkerConfig struct {
	PoolSize   int
	QueueSize  int
	RetryCount int
	RetryDelay time.Duration
	BatchSize  int
}

func Load() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/warehouse-inventory/")

	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("auth.jwt_secret", "change-me-in-production")
	v.SetDefault("auth.jwt_issuer", "wim")
	v.SetDefault("auth.access_token_ttl", "168h")
	v.SetDefault("auth.refresh_token_ttl", "720h")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "wim_user")
	v.SetDefault("database.password", "wim_pass")
	v.SetDefault("database.database", "warehouse_inventory")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_conns", 25)
	v.SetDefault("database.connect_retries", 6)
	v.SetDefault("database.connect_retry_initial_delay", "1s")
	v.SetDefault("database.connect_retry_max_delay", "30s")

	v.SetDefault("log_level", "info")
	v.SetDefault("worker.pool_size", 5)
	v.SetDefault("worker.queue_size", 100)
	v.SetDefault("worker.retry_count", 3)
	v.SetDefault("worker.retry_delay", "1s")
	v.SetDefault("worker.batch_size", 10)

	v.SetDefault("kafka.batch_size", 100)
	v.SetDefault("kafka.batch_timeout", "100ms")
	v.SetDefault("kafka.required_acks", 1)
	v.SetDefault("kafka.compression", "snappy")
	v.SetDefault("kafka.idempotent", true)
	v.SetDefault("kafka.min_bytes", 1024)
	v.SetDefault("kafka.max_bytes", 10485760) // 10MB
	v.SetDefault("kafka.max_wait", "500ms")
	v.SetDefault("kafka.auto_commit", false)
	v.SetDefault("kafka.commit_interval", "5s")
	v.SetDefault("kafka.enable_sasl", false)
	v.SetDefault("kafka.sasl_mechanism", "plain")
	v.SetDefault("kafka.enable_tls", false)

	_ = v.ReadInConfig()

	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.MergeInConfig()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	_ = v.BindEnv("database.url", "DATABASE_URL")
	_ = v.BindEnv("kafka.brokers", "KAFKA_BROKERS")

	return &Config{
		Server: ServerConfig{
			Port: v.GetInt("server.port"),
			Host: v.GetString("server.host"),
		},
		Auth: AuthConfig{
			JWTSecret:       v.GetString("auth.jwt_secret"),
			JWTIssuer:       v.GetString("auth.jwt_issuer"),
			AccessTokenTTL:  v.GetDuration("auth.access_token_ttl"),
			RefreshTokenTTL: v.GetDuration("auth.refresh_token_ttl"),
		},
		Database: DatabaseConfig{
			URL:      v.GetString("database.url"),
			Host:     v.GetString("database.host"),
			Port:     v.GetInt("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			Database: v.GetString("database.database"),
			SSLMode:  v.GetString("database.ssl_mode"),
			MaxConns: int32(v.GetInt("database.max_conns")),

			ConnectRetries:           v.GetInt("database.connect_retries"),
			ConnectRetryInitialDelay: v.GetDuration("database.connect_retry_initial_delay"),
			ConnectRetryMaxDelay:     v.GetDuration("database.connect_retry_max_delay"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("redis.host"),
			Port:     v.GetInt("redis.port"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},
		Kafka: KafkaConfig{
			Brokers:        v.GetStringSlice("kafka.brokers"),
			Topic:          v.GetString("kafka.topic"),
			GroupID:        v.GetString("kafka.group_id"),
			ProducerAsync:  v.GetBool("kafka.producer_async"),
			BatchSize:      v.GetInt("kafka.batch_size"),
			BatchTimeout:   v.GetDuration("kafka.batch_timeout"),
			RequiredAcks:   v.GetInt("kafka.required_acks"),
			Compression:    v.GetString("kafka.compression"),
			Idempotent:     v.GetBool("kafka.idempotent"),
			MinBytes:       v.GetInt("kafka.min_bytes"),
			MaxBytes:       v.GetInt("kafka.max_bytes"),
			MaxWait:        v.GetDuration("kafka.max_wait"),
			AutoCommit:     v.GetBool("kafka.auto_commit"),
			CommitInterval: v.GetDuration("kafka.commit_interval"),
			EnableSASL:     v.GetBool("kafka.enable_sasl"),
			SASLMechanism:  v.GetString("kafka.sasl_mechanism"),
			SASLUsername:   v.GetString("kafka.sasl_username"),
			SASLPassword:   v.GetString("kafka.sasl_password"),
			EnableTLS:      v.GetBool("kafka.enable_tls"),
			TLSCAFile:      v.GetString("kafka.tls_ca_file"),
			TLSCertFile:    v.GetString("kafka.tls_cert_file"),
			TLSKeyFile:     v.GetString("kafka.tls_key_file"),
		},
		Worker: WorkerConfig{
			PoolSize:   v.GetInt("worker.pool_size"),
			QueueSize:  v.GetInt("worker.queue_size"),
			RetryCount: v.GetInt("worker.retry_count"),
			RetryDelay: v.GetDuration("worker.retry_delay"),
			BatchSize:  v.GetInt("worker.batch_size"),
		},
		LogLevel: v.GetString("log_level"),
	}
}
