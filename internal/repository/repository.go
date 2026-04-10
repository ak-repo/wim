package repository

import "github.com/ak-repo/wim/internal/db"

type Repositories struct {
	User       UserRepository
	Auth       AuthRepository
	Product    ProductRepository
	Warehouse  WarehouseRepository
	Location   LocationRepository
	Inventory  InventoryRepository
	RefCode    RefCodeGenerator
	SalesOrder SalesOrderRepository
}

type Dependencies struct {
	DB    *db.DB
	Redis *db.Redis
}

func NewRepositories(deps Dependencies) *Repositories {
	return &Repositories{
		User:       NewUserRepository(deps.DB),
		Auth:       NewAuthRepository(deps.DB),
		Product:    NewProductRepository(deps.DB),
		Warehouse:  NewWarehouseRepository(deps.DB),
		Location:   NewLocationRepository(deps.DB),
		Inventory:  NewInventoryRepository(deps.DB),
		SalesOrder: NewSalesOrderRepository(deps.DB),
		RefCode:    NewRefCodeGenerator(deps.DB, deps.Redis),
	}
}
