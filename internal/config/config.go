package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
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

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
	MaxConns int32
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
}

type WorkerConfig struct {
	PoolSize   int
	QueueSize  int
	RetryCount int
	RetryDelay time.Duration
	BatchSize  int
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/warehouse-inventory/")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("database.max_conns", 25)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("worker.pool_size", 5)
	viper.SetDefault("worker.queue_size", 100)
	viper.SetDefault("worker.retry_count", 3)
	viper.SetDefault("worker.retry_delay", "1s")
	viper.SetDefault("worker.batch_size", 10)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
		}
	}

	viper.AutomaticEnv()

	return &Config{
		Server: ServerConfig{
			Port: viper.GetInt("server.port"),
			Host: viper.GetString("server.host"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			Database: viper.GetString("database.database"),
			SSLMode:  viper.GetString("database.ssl_mode"),
			MaxConns: int32(viper.GetInt("database.max_conns")),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetInt("redis.port"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		Kafka: KafkaConfig{
			Brokers: viper.GetStringSlice("kafka.brokers"),
			Topic:   viper.GetString("kafka.topic"),
			GroupID: viper.GetString("kafka.group_id"),
		},
		Worker: WorkerConfig{
			PoolSize:   viper.GetInt("worker.pool_size"),
			QueueSize:  viper.GetInt("worker.queue_size"),
			RetryCount: viper.GetInt("worker.retry_count"),
			RetryDelay: viper.GetDuration("worker.retry_delay"),
			BatchSize:  viper.GetInt("worker.batch_size"),
		},
		LogLevel: viper.GetString("log_level"),
	}
}
