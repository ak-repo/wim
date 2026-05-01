package handler

import "github.com/ak-repo/wim/internal/service"

type Handler struct {
	Auth      *AuthHandler
	Health    *HealthHandler
	User      *UserHandler
	Customer  *CustomerHandler
	Product   *ProductHandler
	Warehouse *WarehouseHandler
	Location  *LocationHandler
}

func NewHandlers(services *service.Services) *Handler {
	return &Handler{
		Auth:      NewAuthHandler(services),
		Health:    NewHealthHandler(),
		User:      NewUserHandler(services),
		Customer:  NewCustomerHandler(services),
		Product:   NewProductHandler(services),
		Warehouse: NewWarehouseHandler(services),
		Location:  NewLocationHandler(services),
	}
}
