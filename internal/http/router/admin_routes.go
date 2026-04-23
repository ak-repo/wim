package router

import (
	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/internal/http/handler"
	wimMiddleware "github.com/ak-repo/wim/internal/http/middleware"
	"github.com/ak-repo/wim/pkg/auth"
	"github.com/go-chi/chi"
)

func AdminRoutes(r chi.Router, handlers *handler.Handler, tokenManager auth.TokenManager) {
	publicRoutes := r.Route("/adminPublic", func(public chi.Router) {})
	privateRoutes := r.Route("/admin", func(private chi.Router) {
		private.Use(wimMiddleware.RequireAuth(tokenManager))
		private.Use(wimMiddleware.RoleBasedAccessControl(constants.RoleAdmin))
	})

	publicRoutes.Post("/login", handlers.Auth.Login)
	publicRoutes.Post("/register", handlers.Auth.Register)

	// Users Routes
	privateRoutes.Route("/users", func(users chi.Router) {
		users.Get("/", handlers.User.ListUsers)
		users.Post("/", handlers.User.CreateUser)
		users.Put("/:id", handlers.User.UpdateUser)
		users.Delete("/:id", handlers.User.DeleteUser)
		users.Get("/:id", handlers.User.GetUserByID)
	})

	// Product Routes
	privateRoutes.Route("/products", func(products chi.Router) {
		products.Post("/", handlers.Product.CreateProduct)
		products.Get("/", handlers.Product.ListProducts)
		products.Get("/sku/{sku}", handlers.Product.GetProductBySKU)
		products.Get("/{id}", handlers.Product.GetProductByID)
		products.Patch("/{id}", handlers.Product.UpdateProduct)
		products.Delete("/{id}", handlers.Product.DeleteProduct)
	})

	// Product Category Routes
	privateRoutes.Route("/product-categories", func(categories chi.Router) {
		categories.Post("/", handlers.ProductCategory.CreateProductCategory)
		categories.Get("/", handlers.ProductCategory.ListProductCategories)
		categories.Get("/{id}", handlers.ProductCategory.GetProductCategoryByID)
		categories.Patch("/{id}", handlers.ProductCategory.UpdateProductCategory)
		categories.Delete("/{id}", handlers.ProductCategory.DeleteProductCategory)
	})

	// location Routes
	privateRoutes.Route("/locations", func(locations chi.Router) {
		locations.Post("/", handlers.Location.CreateLocation)
		locations.Get("/", handlers.Location.ListLocations)
		locations.Get("/code/{code}", handlers.Location.GetLocationByCode)
		locations.Get("/warehouse/{warehouseId}", handlers.Location.ListLocationsByWarehouse)
		locations.Get("/{id}", handlers.Location.GetLocationByID)
		locations.Patch("/{id}", handlers.Location.UpdateLocation)
		locations.Delete("/{id}", handlers.Location.DeleteLocation)
	})

	// Warehouse Routes
	privateRoutes.Route("/warehouses", func(warehouses chi.Router) {
		warehouses.Post("/", handlers.Warehouse.CreateWarehouse)
		warehouses.Get("/", handlers.Warehouse.ListWarehouses)
		warehouses.Get("/code/{code}", handlers.Warehouse.GetWarehouseByCode)
		warehouses.Get("/{id}", handlers.Warehouse.GetWarehouseByID)
		warehouses.Put("/{id}", handlers.Warehouse.UpdateWarehouse)
		warehouses.Delete("/{id}", handlers.Warehouse.DeleteWarehouse)
	})

	// Inventory Routes
	privateRoutes.Route("/inventory", func(inventory chi.Router) {
		inventory.Post("/adjust", handlers.Inventory.AdjustInventory)
		inventory.Get("/", handlers.Inventory.ListInventory)
		inventory.Get("/{id}", handlers.Inventory.GetInventoryByID)
	})

	// Stock Movement Routes
	privateRoutes.Route("/stock-movements", func(movements chi.Router) {
		movements.Get("/", handlers.Inventory.ListStockMovements)
	})

	// Sales Order Routes
	privateRoutes.Route("/sales-orders", func(orders chi.Router) {
		orders.Post("/", handlers.SalesOrder.CreateSalesOrder)
		orders.Get("/", handlers.SalesOrder.ListSalesOrders)
		orders.Get("/ref", handlers.SalesOrder.GetSalesOrderByRefCode)
		orders.Get("/{id}", handlers.SalesOrder.GetSalesOrderByID)
		orders.Put("/{id}", handlers.SalesOrder.UpdateSalesOrder)
		orders.Patch("/{id}/cancel", handlers.SalesOrder.CancelSalesOrder)
		orders.Patch("/{id}/allocate", handlers.SalesOrder.AllocateSalesOrder)
		orders.Patch("/{id}/deallocate", handlers.SalesOrder.DeallocateSalesOrder)
		orders.Patch("/{id}/ship", handlers.SalesOrder.ShipSalesOrder)
	})

	// Dashboard Routes
	privateRoutes.Route("/dashboard", func(dashboard chi.Router) {
		dashboard.Get("/counts", handlers.Dashboard.TotalCounts)
	})

	// User Role Routes
	privateRoutes.Route("/user-roles", func(roles chi.Router) {
		roles.Post("/", handlers.UserRole.CreateUserRole)
		roles.Get("/", handlers.UserRole.ListUserRoles)
		roles.Get("/{id}", handlers.UserRole.GetUserRoleByID)
		roles.Patch("/{id}", handlers.UserRole.UpdateUserRole)
		roles.Delete("/{id}", handlers.UserRole.DeleteUserRole)
	})

}
