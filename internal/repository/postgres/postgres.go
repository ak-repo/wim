package postgres

import (
	"context"
	"fmt"

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
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
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
	}
}
