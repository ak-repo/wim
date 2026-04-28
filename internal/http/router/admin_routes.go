package router

import (
	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/internal/http/handler"
	wimMiddleware "github.com/ak-repo/wim/internal/http/middleware"
	"github.com/ak-repo/wim/pkg/auth"
	"github.com/go-chi/chi"
)

func AdminRoutes(r chi.Router, handlers *handler.Handler, tokenManager auth.TokenManager) {
	r.Route("/adminPublic", func(public chi.Router) {
		public.Post("/login", handlers.Auth.Login)
		public.Post("/register", handlers.Auth.Register)
	})

	r.Route("/admin", func(admin chi.Router) {
		admin.Use(wimMiddleware.RequireAuth(tokenManager))
		admin.Get("/me", handlers.Auth.Me)

		admin.Route("/products", func(products chi.Router) {
			products.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
			products.Post("/", handlers.Product.CreateProduct)
			products.Get("/", handlers.Product.ListProducts)
			products.Get("/sku/{sku}", handlers.Product.GetProductBySKU)
			products.Get("/{id}", handlers.Product.GetProductByID)
			products.Patch("/{id}", handlers.Product.UpdateProduct)
			products.Delete("/{id}", handlers.Product.DeleteProduct)
		})

		admin.Route("/purchase-orders", func(orders chi.Router) {
			orders.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
			orders.Post("/", handlers.PurchaseOrder.CreatePurchaseOrder)
			orders.Get("/", handlers.PurchaseOrder.ListPurchaseOrders)
			orders.Get("/ref", handlers.PurchaseOrder.GetPurchaseOrderByRefCode)
			orders.Get("/{id}", handlers.PurchaseOrder.GetPurchaseOrderByID)
			orders.Post("/{id}/receive", handlers.PurchaseOrder.ReceivePurchaseOrder)
			orders.Post("/{id}/put-away", handlers.PurchaseOrder.PutAwayPurchaseOrder)
		})

		admin.Route("/product-categories", func(categories chi.Router) {
			categories.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
			categories.Post("/", handlers.ProductCategory.CreateProductCategory)
			categories.Get("/", handlers.ProductCategory.ListProductCategories)
			categories.Get("/{id}", handlers.ProductCategory.GetProductCategoryByID)
			categories.Patch("/{id}", handlers.ProductCategory.UpdateProductCategory)
			categories.Delete("/{id}", handlers.ProductCategory.DeleteProductCategory)
		})

		admin.Route("/locations", func(locations chi.Router) {
			locations.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
			locations.Post("/", handlers.Location.CreateLocation)
			locations.Get("/", handlers.Location.ListLocations)
			locations.Get("/code/{code}", handlers.Location.GetLocationByCode)
			locations.Get("/warehouse/{warehouseId}", handlers.Location.ListLocationsByWarehouse)
			locations.Get("/{id}", handlers.Location.GetLocationByID)
			locations.Patch("/{id}", handlers.Location.UpdateLocation)
			locations.Delete("/{id}", handlers.Location.DeleteLocation)
		})

		admin.Route("/warehouses", func(warehouses chi.Router) {
			warehouses.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
			warehouses.Post("/", handlers.Warehouse.CreateWarehouse)
			warehouses.Get("/", handlers.Warehouse.ListWarehouses)
			warehouses.Get("/code/{code}", handlers.Warehouse.GetWarehouseByCode)
			warehouses.Get("/{id}", handlers.Warehouse.GetWarehouseByID)
			warehouses.Put("/{id}", handlers.Warehouse.UpdateWarehouse)
			warehouses.Delete("/{id}", handlers.Warehouse.DeleteWarehouse)
		})

		admin.Route("/inventory", func(inventory chi.Router) {
			inventory.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
			inventory.Post("/adjust", handlers.Inventory.AdjustInventory)
			inventory.Get("/", handlers.Inventory.ListInventory)
			inventory.Get("/{id}", handlers.Inventory.GetInventoryByID)
		})

		admin.Route("/stock-movements", func(movements chi.Router) {
			movements.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
			movements.Get("/", handlers.Inventory.ListStockMovements)
		})

		admin.Route("/sales-orders", func(orders chi.Router) {
			orders.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
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

		admin.Route("/picking-tasks", func(picking chi.Router) {
			picking.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
			picking.Post("/", handlers.Picking.CreatePickingTask)
			picking.Get("/", handlers.Picking.ListPickingTasks)
			picking.Get("/ref", handlers.Picking.GetPickingTaskByRefCode)
			picking.Get("/{id}", handlers.Picking.GetPickingTaskByID)
			picking.Post("/{id}/assign", handlers.Picking.AssignPickingTask)
			picking.Post("/{id}/start", handlers.Picking.StartPickingTask)
			picking.Post("/{id}/pick", handlers.Picking.PickItem)
			picking.Post("/{id}/complete", handlers.Picking.CompletePickingTask)
			picking.Post("/{id}/cancel", handlers.Picking.CancelPickingTask)
		})

		admin.Route("/dashboard", func(dashboard chi.Router) {
			dashboard.Use(wimMiddleware.RequireAnyRole(constants.RoleAdmin))
			dashboard.Get("/counts", handlers.Dashboard.TotalCounts)
		})

		admin.Route("/users", func(users chi.Router) {
			users.Use(wimMiddleware.RequireAnyRole(constants.RoleSuperAdmin))
			users.Get("/", handlers.User.ListUsers)
			users.Post("/", handlers.User.CreateUser)
			users.Put("/:id", handlers.User.UpdateUser)
			users.Delete("/:id", handlers.User.DeleteUser)
			users.Get("/:id", handlers.User.GetUserByID)
		})

		admin.Route("/user-roles", func(roles chi.Router) {
			roles.Use(wimMiddleware.RequireAnyRole(constants.RoleSuperAdmin))
			roles.Post("/", handlers.UserRole.CreateUserRole)
			roles.Get("/", handlers.UserRole.ListUserRoles)
			roles.Get("/{id}", handlers.UserRole.GetUserRoleByID)
			roles.Patch("/{id}", handlers.UserRole.UpdateUserRole)
			roles.Delete("/{id}", handlers.UserRole.DeleteUserRole)
		})
	})
}