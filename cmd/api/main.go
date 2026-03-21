package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ak-repo/wim/internal/config"
	dbutil "github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/event"
	auditHandler "github.com/ak-repo/wim/internal/handler/audit"
	batchHandler "github.com/ak-repo/wim/internal/handler/batch"
	inventoryHandler "github.com/ak-repo/wim/internal/handler/inventory"
	orderHandler "github.com/ak-repo/wim/internal/handler/order"
	productHandler "github.com/ak-repo/wim/internal/handler/product"
	reportHandler "github.com/ak-repo/wim/internal/handler/report"
	transferHandler "github.com/ak-repo/wim/internal/handler/transfer"
	warehouseHandler "github.com/ak-repo/wim/internal/handler/warehouse"
	"github.com/ak-repo/wim/internal/middleware"
	"github.com/ak-repo/wim/internal/repository/postgres"
	auditSvc "github.com/ak-repo/wim/internal/service/audit"
	batchSvc "github.com/ak-repo/wim/internal/service/batch"
	inventorySvc "github.com/ak-repo/wim/internal/service/inventory"
	orderSvc "github.com/ak-repo/wim/internal/service/order"
	productSvc "github.com/ak-repo/wim/internal/service/product"
	reportSvc "github.com/ak-repo/wim/internal/service/report"
	transferSvc "github.com/ak-repo/wim/internal/service/transfer"
	warehouseSvc "github.com/ak-repo/wim/internal/service/warehouse"
	"github.com/ak-repo/wim/pkg/logger"
	"github.com/ak-repo/wim/pkg/tracing"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	log.Info("starting warehouse inventory API")

	shutdownTracer, err := tracing.Init(ctx, "warehouse-inventory-api")
	if err != nil {
		log.Fatal("failed to initialize tracing", "error", err)
	}
	defer shutdownTracer(context.Background())

	pgDB, err := postgres.NewConnection(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	defer pgDB.Close()

	sqlDB, err := dbutil.OpenSQLConnection(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to open migration database connection", "error", err)
	}
	defer sqlDB.Close()

	log.Info("running database migrations")
	if err := dbutil.RunMigrations(sqlDB); err != nil {
		log.Fatal("database migration failed", "error", err)
	}
	log.Info("database migrations completed")

	redisClient, err := postgres.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatal("failed to connect to redis", "error", err)
	}
	defer redisClient.Close()

	repositories := postgres.NewRepositories(pgDB, redisClient)
	publisher := event.NewKafkaPublisher(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer publisher.Close()

	productService := productSvc.NewService(repositories.Product, repositories.Barcode, repositories.AuditLog, publisher)
	warehouseService := warehouseSvc.NewService(repositories.Warehouse, repositories.Location, repositories.AuditLog, publisher)
	batchService := batchSvc.NewService(repositories.Batch)
	inventoryService := inventorySvc.NewService(repositories.Inventory, repositories.StockMovement, repositories.Batch, repositories.AuditLog, publisher)
	orderService := orderSvc.NewService(
		repositories.PurchaseOrder,
		repositories.SalesOrder,
		repositories.Inventory,
		repositories.StockMovement,
		repositories.AuditLog,
		publisher,
		pgDB,
	)
	transferService := transferSvc.NewService(repositories.Transfer, repositories.Inventory, repositories.StockMovement, repositories.AuditLog, publisher, pgDB)
	auditService := auditSvc.NewService(repositories.AuditLog)
	reportService := reportSvc.NewService(repositories.Inventory, repositories.StockMovement, repositories.Batch)

	productHandler := productHandler.NewHandler(productService)
	warehouseHandler := warehouseHandler.NewHandler(warehouseService)
	batchHandler := batchHandler.NewHandler(batchService)
	inventoryHandler := inventoryHandler.NewHandler(inventoryService)
	orderHandler := orderHandler.NewHandler(orderService)
	transferHandler := transferHandler.NewHandler(transferService)
	auditHandler := auditHandler.NewHandler(auditService)
	reportHandler := reportHandler.NewHandler(reportService)

	router := gin.Default()
	router.Use(middleware.Tracing("warehouse-inventory-api"))
	router.Use(middleware.Logger(log))
	router.Use(middleware.ErrorHandler(log))
	router.Use(middleware.Recovery())

	setupRoutes(router, productHandler, warehouseHandler, batchHandler, inventoryHandler, orderHandler, transferHandler, auditHandler, reportHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		log.Info("server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced to shutdown", "error", err)
	}

	log.Info("server exited")
}

func setupRoutes(
	r *gin.Engine,
	productHandler *productHandler.Handler,
	warehouseHandler *warehouseHandler.Handler,
	batchHandler *batchHandler.Handler,
	inventoryHandler *inventoryHandler.Handler,
	orderHandler *orderHandler.Handler,
	transferHandler *transferHandler.Handler,
	auditHandler *auditHandler.Handler,
	reportHandler *reportHandler.Handler,
) {
	api := r.Group("/api/v1")

	products := api.Group("/products")
	{
		products.GET("", productHandler.List)
		products.GET("/:id", productHandler.Get)
		products.POST("", productHandler.Create)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}

	warehouses := api.Group("/warehouses")
	{
		warehouses.GET("", warehouseHandler.List)
		warehouses.GET("/:id", warehouseHandler.Get)
		warehouses.POST("", warehouseHandler.Create)
		warehouses.PUT("/:id", warehouseHandler.Update)
		warehouses.DELETE("/:id", warehouseHandler.Delete)

		warehouses.GET("/:id/locations", warehouseHandler.ListLocations)
		warehouses.POST("/:id/locations", warehouseHandler.CreateLocation)
	}

	inventoryGroup := api.Group("/inventory")
	{
		inventoryGroup.GET("", inventoryHandler.List)
		inventoryGroup.GET("/warehouse/:warehouse_id", inventoryHandler.GetByWarehouse)
		inventoryGroup.GET("/product/:product_id", inventoryHandler.GetByProduct)
		inventoryGroup.POST("/adjust", inventoryHandler.Adjust)
	}

	batches := api.Group("/batches")
	{
		batches.GET("", batchHandler.ListByProduct)
		batches.GET("/expiring", batchHandler.ExpiringSoon)
		batches.GET("/:id", batchHandler.Get)
	}

	orders := api.Group("/orders")
	{
		purchaseOrders := orders.Group("/purchase")
		{
			purchaseOrders.GET("", orderHandler.ListPurchaseOrders)
			purchaseOrders.GET("/:id", orderHandler.GetPurchaseOrder)
			purchaseOrders.POST("", orderHandler.CreatePurchaseOrder)
			purchaseOrders.POST("/:id/receive", orderHandler.ReceivePurchaseOrder)
		}

		salesOrders := orders.Group("/sales")
		{
			salesOrders.GET("", orderHandler.ListSalesOrders)
			salesOrders.GET("/:id", orderHandler.GetSalesOrder)
			salesOrders.POST("", orderHandler.CreateSalesOrder)
			salesOrders.POST("/:id/allocate", orderHandler.AllocateSalesOrder)
			salesOrders.POST("/:id/ship", orderHandler.ShipSalesOrder)
		}
	}

	transfers := api.Group("/transfers")
	{
		transfers.GET("", transferHandler.List)
		transfers.GET("/:id", transferHandler.Get)
		transfers.POST("", transferHandler.Create)
		transfers.POST("/:id/approve", transferHandler.Approve)
		transfers.POST("/:id/ship", transferHandler.Ship)
		transfers.POST("/:id/receive", transferHandler.Receive)
	}

	reports := api.Group("/reports")
	{
		reports.GET("/inventory", reportHandler.Inventory)
		reports.GET("/movements", reportHandler.Movements)
		reports.GET("/expiry", reportHandler.Expiry)
	}

	audit := api.Group("/audit")
	{
		audit.GET("/logs", auditHandler.List)
	}

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
