package service

import (
	"github.com/ak-repo/wim/internal/event"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/pkg/auth"
)

type Services struct {
	User       UserService
	Auth       AuthService
	Product    ProductService
	Warehouse  WarehouseService
	Location   LocationService
	Inventory  InventoryService
	SalesOrder SalesOrderService
}

type Dependencies struct {
	Repositories   *repository.Repositories
	PasswordHasher auth.PasswordHasher
	TokenManager   auth.TokenManager
	EventPublisher event.EventPublisher
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		User:       NewUserService(deps.Repositories, deps.PasswordHasher),
		Auth:       NewAuthService(deps.Repositories, deps.TokenManager, deps.PasswordHasher),
		Product:    NewProductService(deps.Repositories),
		Warehouse:  NewWarehouseService(deps.Repositories),
		Location:   NewLocationService(deps.Repositories),
		Inventory:  NewInventoryService(deps.Repositories),
		SalesOrder: NewSalesOrderService(deps.Repositories, deps.EventPublisher),
	}
}
