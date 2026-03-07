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
	inventoryHandler "github.com/ak-repo/wim/internal/handler/inventory"
	orderHandler "github.com/ak-repo/wim/internal/handler/order"
	productHandler "github.com/ak-repo/wim/internal/handler/product"
	warehouseHandler "github.com/ak-repo/wim/internal/handler/warehouse"
	"github.com/ak-repo/wim/internal/middleware"
	"github.com/ak-repo/wim/internal/repository/postgres"
	inventorySvc "github.com/ak-repo/wim/internal/service/inventory"
	orderSvc "github.com/ak-repo/wim/internal/service/order"
	productSvc "github.com/ak-repo/wim/internal/service/product"
	warehouseSvc "github.com/ak-repo/wim/internal/service/warehouse"
	"github.com/ak-repo/wim/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	log.Info("starting warehouse inventory API")

	db, err := postgres.NewConnection(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	defer db.Close()

	redisClient, err := postgres.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatal("failed to connect to redis", "error", err)
	}
	defer redisClient.Close()

	repositories := postgres.NewRepositories(db, redisClient)

	productService := productSvc.NewService(repositories.Product, repositories.Barcode)
	warehouseService := warehouseSvc.NewService(repositories.Warehouse, repositories.Location)
	inventoryService := inventorySvc.NewService(repositories.Inventory, repositories.StockMovement, repositories.Batch)
	orderService := orderSvc.NewService(
		repositories.PurchaseOrder,
		repositories.SalesOrder,
		repositories.Inventory,
		repositories.StockMovement,
	)

	productHandler := productHandler.NewHandler(productService)
	warehouseHandler := warehouseHandler.NewHandler(warehouseService)
	inventoryHandler := inventoryHandler.NewHandler(inventoryService)
	orderHandler := orderHandler.NewHandler(orderService)

	router := gin.Default()
	router.Use(middleware.Logger(log))
	router.Use(middleware.ErrorHandler(log))
	router.Use(middleware.Recovery())

	setupRoutes(router, productHandler, warehouseHandler, inventoryHandler, orderHandler)

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
	inventoryHandler *inventoryHandler.Handler,
	orderHandler *orderHandler.Handler,
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

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
