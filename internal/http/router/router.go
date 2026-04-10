package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/internal/http/handler"
	wimMiddleware "github.com/ak-repo/wim/internal/http/middleware"
	"github.com/ak-repo/wim/pkg/auth"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func SetupRoutes(handlers *handler.Handler, tokenManager auth.TokenManager) http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)

	allowedOrigins := []string{"http://localhost:5174", "http://localhost:3050"}
	if raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); raw != "" {
		parts := strings.Split(raw, ",")
		allowedOrigins = allowedOrigins[:0]
		for _, part := range parts {
			origin := strings.TrimSpace(part)
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", handlers.Health.Check)

	apiPrefix := strings.TrimSpace(os.Getenv("WIM_API_PREFIX"))
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}
	if !strings.HasPrefix(apiPrefix, "/") {
		apiPrefix = "/" + apiPrefix
	}

	r.Route(apiPrefix, func(api chi.Router) {
		AdminRoutes(api, handlers, tokenManager)
		// api.With(wimMiddleware.RequireAuth(tokenManager)).Get("/me", handlers.Auth.Me)
	})

	return r
}

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

}
