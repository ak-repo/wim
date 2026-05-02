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
		users.Get("/{id}", handlers.User.GetUserByID)
		users.Post("/",handlers.User.CreateUser)
		users.Put("/{id}",handlers.User.UpdateUser)
		users.Delete("/{id}",handlers.User.DeleteUser)
	})

	// Customer Routes
	privateRoutes.Route("/customers", func(customers chi.Router) {
		customers.Post("/", handlers.Customer.CreateCustomer)
		customers.Get("/", handlers.Customer.ListCustomers)
		customers.Get("/email/{email}", handlers.Customer.GetCustomerByEmail)
		customers.Get("/{id}", handlers.Customer.GetCustomerByID)
		customers.Patch("/{id}", handlers.Customer.UpdateCustomer)
		customers.Delete("/{id}", handlers.Customer.DeleteCustomer)
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

	// Customer Types Routes
	privateRoutes.Route("/customer-types", func(customerTypes chi.Router) {
		customerTypes.Get("/", handlers.CustomerType.ListCustomerTypes)
		customerTypes.Post("/", handlers.CustomerType.CreateCustomerType)
		customerTypes.Get("/{id}", handlers.CustomerType.GetCustomerTypeByID)
		customerTypes.Put("/{id}", handlers.CustomerType.UpdateCustomerType)
		customerTypes.Delete("/{id}", handlers.CustomerType.DeleteCustomerType)
	})

	// User Roles Routes
	privateRoutes.Route("/user-roles", func(userRoles chi.Router) {
		userRoles.Get("/", handlers.UserRole.ListUserRoles)
		userRoles.Post("/", handlers.UserRole.CreateUserRole)
		userRoles.Get("/{id}", handlers.UserRole.GetUserRoleByID)
		userRoles.Put("/{id}", handlers.UserRole.UpdateUserRole)
		userRoles.Delete("/{id}", handlers.UserRole.DeleteUserRole)
	})

	// Product Categories Routes
	privateRoutes.Route("/product-categories", func(productCategories chi.Router) {
		productCategories.Get("/", handlers.ProductCategory.ListProductCategories)
		productCategories.Post("/", handlers.ProductCategory.CreateProductCategory)
		productCategories.Get("/{id}", handlers.ProductCategory.GetProductCategoryByID)
		productCategories.Put("/{id}", handlers.ProductCategory.UpdateProductCategory)
		productCategories.Delete("/{id}", handlers.ProductCategory.DeleteProductCategory)
	})
}
