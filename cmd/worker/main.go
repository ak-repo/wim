package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ak-repo/wim/internal/config"
	dbutil "github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/queue"
	"github.com/ak-repo/wim/internal/repository/postgres"
	batchSvc "github.com/ak-repo/wim/internal/service/batch"
	inventorySvc "github.com/ak-repo/wim/internal/service/inventory"
	reportSvc "github.com/ak-repo/wim/internal/service/report"
	"github.com/ak-repo/wim/internal/worker"
	"github.com/ak-repo/wim/pkg/logger"
	"github.com/ak-repo/wim/pkg/tracing"

	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	appLog := logger.New(cfg.LogLevel)
	appLog.Info("starting warehouse inventory worker")

	shutdownTracer, err := tracing.Init(ctx, "warehouse-inventory-worker")
	if err != nil {
		appLog.Fatal("failed to initialize tracing", "error", err)
	}
	defer shutdownTracer(context.Background())

	pgDB, err := postgres.NewConnection(ctx, cfg.Database)
	if err != nil {
		appLog.Fatal("failed to connect to database", "error", err)
	}
	defer pgDB.Close()

	sqlDB, err := dbutil.OpenSQLConnection(ctx, cfg.Database)
	if err != nil {
		appLog.Fatal("failed to open migration database connection", "error", err)
	}
	defer sqlDB.Close()

	appLog.Info("running database migrations")
	if err := dbutil.RunMigrations(sqlDB); err != nil {
		appLog.Fatal("database migration failed", "error", err)
	}
	appLog.Info("database migrations completed")

	redisClient, err := postgres.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		appLog.Fatal("failed to connect to redis", "error", err)
	}
	defer redisClient.Close()

	repositories := postgres.NewRepositories(pgDB, redisClient)

	batchService := batchSvc.NewService(repositories.Batch)
	inventoryService := inventorySvc.NewService(repositories.Inventory, repositories.StockMovement, repositories.Batch, repositories.AuditLog, nil)
	reportService := reportSvc.NewService(repositories.Inventory, repositories.StockMovement, repositories.Batch)

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    cfg.Kafka.Topic,
		GroupID:  cfg.Kafka.GroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer kafkaReader.Close()

	workerPool := worker.NewPool(ctx, cfg.Worker, appLog, worker.Dependencies{
		Inventory: inventoryService,
		Batch:     batchService,
		Report:    reportService,
	})
	processor := queue.NewProcessor(appLog, workerPool)

	signalCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		appLog.Info("shutting down worker")
		cancel()
	}()

	if err := processor.Run(signalCtx, kafkaReader); err != nil {
		appLog.Error("worker error", "error", err)
		appLog.Fatal("worker stopped", "error", err)
	}

	appLog.Info("worker exited")
}
