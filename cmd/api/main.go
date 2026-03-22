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

	"github.com/ak-repo/wim/internal/config"
	"github.com/ak-repo/wim/internal/db"
	dbutil "github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/http/handler"
	"github.com/ak-repo/wim/internal/http/router"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/internal/service"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	log.Println("starting warehouse inventory API")

	pgDB, err := db.NewConnection(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	defer pgDB.Close()

	sqlDB, err := dbutil.OpenSQLConnection(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to open migration database connection", "error", err)
	}
	defer sqlDB.Close()

	log.Println("running database migrations")
	if err := dbutil.RunMigrations(sqlDB); err != nil {
		log.Fatal("database migration failed", "error", err)
	}
	log.Println("database migrations completed")

	redisClient, err := db.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatal("failed to connect to redis", "error", err)
	}
	defer redisClient.Close()

	repos := repository.NewRepositories(repository.Dependencies{
		DB:    pgDB,
		Redis: redisClient,
	})

	services := service.NewServices(service.Dependencies{
		Repositories: repos,
	})

	handlers := handler.NewHandlers(services)

	router := router.SetupRoutes(handlers)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
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
