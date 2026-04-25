package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

/*
========================
STRUCT DEFINITIONS
========================
*/

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
	Host string
	Port int
}

type AuthConfig struct {
	JWTSecret       string
	JWTIssuer       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type DatabaseConfig struct {
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
	URL                      string
}

func (d DatabaseConfig) DSN() string {
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

	ProducerAsync bool
	BatchSize     int
	BatchTimeout  time.Duration
	RequiredAcks  int
	Compression   string
	Idempotent    bool

	MinBytes       int
	MaxBytes       int
	MaxWait        time.Duration
	AutoCommit     bool
	CommitInterval time.Duration

	EnableSASL    bool
	SASLMechanism string
	SASLUsername  string
	SASLPassword  string

	EnableTLS   bool
	TLSCAFile   string
	TLSCertFile string
	TLSKeyFile  string
}

type WorkerConfig struct {
	PoolSize   int
	QueueSize  int
	RetryCount int
	RetryDelay time.Duration
	BatchSize  int
}

/*
========================
LOAD FUNCTION
========================
*/

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	v := viper.New()

	// YAML
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	// Bind ENV
	bindEnv(v)

	// Convert Kafka brokers
	if brokers := v.GetString("kafka.brokers"); brokers != "" {
		v.Set("kafka.brokers", strings.Split(brokers, ","))
	}

	// Validate
	if err := validate(v); err != nil {
		return nil, err
	}

	return buildConfig(v), nil
}

/*
========================
ENV BINDING
========================
*/

func bindEnv(v *viper.Viper) {
	// DATABASE
	v.BindEnv("database.host", "POSTGRES_HOST")
	v.BindEnv("database.port", "POSTGRES_PORT")
	v.BindEnv("database.user", "POSTGRES_USER")
	v.BindEnv("database.password", "POSTGRES_PASSWORD")
	v.BindEnv("database.database", "POSTGRES_DB")
	v.BindEnv("database.max_conns", "DATABASE_MAX_CONNS")
	v.BindEnv("database.url","DATABASE_URL")

	// REDIS
	v.BindEnv("redis.host", "REDIS_HOST")
	v.BindEnv("redis.port", "REDIS_PORT")
	v.BindEnv("redis.password", "REDIS_PASSWORD")
	v.BindEnv("redis.db", "REDIS_DB")

	// AUTH
	v.BindEnv("auth.jwt_secret", "AUTH_JWT_SECRET")
	v.BindEnv("auth.jwt_issuer", "AUTH_JWT_ISSUER")

	// KAFKA
	v.BindEnv("kafka.brokers", "KAFKA_BROKERS")
	v.BindEnv("kafka.topic", "KAFKA_TOPIC")
	v.BindEnv("kafka.group_id", "KAFKA_GROUP_ID")

	// OPTIONAL
	v.BindEnv("server.host", "SERVER_HOST")
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("log_level", "LOG_LEVEL")
}

/*
========================
VALIDATION
========================
*/

func validate(v *viper.Viper) error {
	required := map[string]string{
		"database.host":      "POSTGRES_HOST",
		"database.port":      "POSTGRES_PORT",
		"database.user":      "POSTGRES_USER",
		"database.password":  "POSTGRES_PASSWORD",
		"database.database":  "POSTGRES_DB",
		"database.max_conns": "DATABASE_MAX_CONNS",

		"redis.host":     "REDIS_HOST",
		"redis.port":     "REDIS_PORT",
		"redis.password": "REDIS_PASSWORD",

		"auth.jwt_secret": "AUTH_JWT_SECRET",
		"auth.jwt_issuer": "AUTH_JWT_ISSUER",

		"kafka.brokers":  "KAFKA_BROKERS",
		"kafka.topic":    "KAFKA_TOPIC",
		"kafka.group_id": "KAFKA_GROUP_ID",
	}

	for key, env := range required {
		if !v.IsSet(key) || v.Get(key) == "" {
			return fmt.Errorf("missing required config (%s) from env: %s", key, env)
		}
	}

	if v.GetInt("database.max_conns") < 1 {
		return fmt.Errorf("database.max_conns must be >= 1")
	}

	return nil
}

/*
========================
BUILD CONFIG
========================
*/

func buildConfig(v *viper.Viper) *Config {
	return &Config{
		Server: ServerConfig{
			Host: v.GetString("server.host"),
			Port: v.GetInt("server.port"),
		},
		Auth: AuthConfig{
			JWTSecret:       v.GetString("auth.jwt_secret"),
			JWTIssuer:       v.GetString("auth.jwt_issuer"),
			AccessTokenTTL:  v.GetDuration("auth.access_token_ttl"),
			RefreshTokenTTL: v.GetDuration("auth.refresh_token_ttl"),
		},
		Database: DatabaseConfig{
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
