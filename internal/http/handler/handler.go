package handler

import "github.com/ak-repo/wim/internal/service"

type Handler struct {
	Auth            *AuthHandler
	Health          *HealthHandler
	User            *UserHandler
	Product         *ProductHandler
	ProductCategory *ProductCategoryHandler
	Warehouse       *WarehouseHandler
	Location        *LocationHandler
	Inventory       *InventoryHandler
	SalesOrder      *SalesOrderHandler
	Dashboard       *DashboardHandler
	UserRole        *UserRoleHandler
}

func NewHandlers(services *service.Services) *Handler {
	return &Handler{
		Auth:            NewAuthHandler(services),
		Health:          NewHealthHandler(),
		User:            NewUserHandler(services),
		Product:         NewProductHandler(services),
		ProductCategory: NewProductCategoryHandler(services),
		Warehouse:       NewWarehouseHandler(services),
		Location:        NewLocationHandler(services),
		Inventory:       NewInventoryHandler(services),
		SalesOrder:      NewSalesOrderHandler(services),
		Dashboard:       NewDashboardHandler(services),
		UserRole:        NewUserRoleHandler(services),
	}
}
