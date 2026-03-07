package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ak-repo/wim/internal/config"
	"github.com/ak-repo/wim/internal/queue"
	"github.com/ak-repo/wim/internal/worker"
	"github.com/ak-repo/wim/pkg/logger"

	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	appLog := logger.New(cfg.LogLevel)
	appLog.Info("starting warehouse inventory worker")

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    cfg.Kafka.Topic,
		GroupID:  cfg.Kafka.GroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer kafkaReader.Close()

	workerPool := worker.NewPool(ctx, cfg.Worker, appLog)
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
		log.Fatal(err)
	}

	appLog.Info("worker exited")
}
