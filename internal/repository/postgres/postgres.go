package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ak-repo/wim/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	Pool *pgxpool.Pool
}

type Redis struct {
	Client *redis.Client
}

func NewConnection(ctx context.Context, cfg config.DatabaseConfig) (*DB, error) {
	connStr := cfg.DSN()

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	retries := cfg.ConnectRetries
	if retries < 1 {
		retries = 1
	}

	delay := cfg.ConnectRetryInitialDelay
	if delay <= 0 {
		delay = time.Second
	}

	maxDelay := cfg.ConnectRetryMaxDelay
	if maxDelay <= 0 {
		maxDelay = 30 * time.Second
	}

	var pingErr error
	for attempt := 1; attempt <= retries; attempt++ {
		pingErr = pool.Ping(ctx)
		if pingErr == nil {
			break
		}

		if attempt == retries {
			pool.Close()
			return nil, fmt.Errorf("failed to ping database after %d attempts: %w", retries, pingErr)
		}

		time.Sleep(delay)
		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Redis{Client: client}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

type Repositories struct {
	Product       ProductRepository
	Warehouse     WarehouseRepository
	Location      LocationRepository
	Inventory     InventoryRepository
	StockMovement StockMovementRepository
	Batch         BatchRepository
	Barcode       BarcodeRepository
	PurchaseOrder PurchaseOrderRepository
	SalesOrder    SalesOrderRepository
	AuditLog      AuditLogRepository
	Transfer      TransferRepository
}

func NewRepositories(db *DB, redis *Redis) *Repositories {
	return &Repositories{
		Product:       NewProductRepository(db),
		Warehouse:     NewWarehouseRepository(db),
		Location:      NewLocationRepository(db),
		Inventory:     NewInventoryRepository(db, redis),
		StockMovement: NewStockMovementRepository(db),
		Batch:         NewBatchRepository(db),
		Barcode:       NewBarcodeRepository(db),
		PurchaseOrder: NewPurchaseOrderRepository(db),
		SalesOrder:    NewSalesOrderRepository(db),
		AuditLog:      NewAuditLogRepository(db),
		Transfer:      NewTransferRepository(db),
	}
}
