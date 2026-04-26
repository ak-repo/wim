package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ak-repo/wim/config"
	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/event"
	"github.com/ak-repo/wim/internal/http/handler"
	"github.com/ak-repo/wim/internal/http/router"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", "error", err)
	}

	log.Println("starting warehouse inventory API")

	pgDB, err := db.NewConnection(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	defer pgDB.Close()

	redisClient, err := db.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatal("failed to connect to redis", "error", err)
	}
	defer redisClient.Close()

	repos := repository.NewRepositories(repository.Dependencies{
		DB:    pgDB,
		Redis: redisClient,
	})

	passwordHasher := auth.NewBcryptPasswordHasher(bcrypt.DefaultCost)
	tokenManager := auth.NewJWTTokenManager(cfg.Auth.JWTSecret, cfg.Auth.JWTIssuer, cfg.Auth.AccessTokenTTL)

	// Initialize event publisher (Kafka if configured, otherwise mock)
	var eventPublisher event.EventPublisher
	if len(cfg.Kafka.Brokers) > 0 {
		var err error
		eventPublisher, err = event.InitKafka(cfg.Kafka)
		if err != nil {
			log.Printf("Failed to initialize Kafka producer, falling back to mock: %v", err)
			eventPublisher = event.NewMockPublisher()
		}
	} else {
		eventPublisher = event.NewMockPublisher()
	}
	defer eventPublisher.Close()

	services := service.NewServices(service.Dependencies{
		Repositories:   repos,
		PasswordHasher: passwordHasher,
		TokenManager:   tokenManager,
		EventPublisher: eventPublisher,
	})

	handlers := handler.NewHandlers(services)

	router := router.SetupRoutes(handlers, tokenManager, cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	go func() {
		log.Println("server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced to shutdown", "error", err)
	}

	log.Println("server exited")
}
