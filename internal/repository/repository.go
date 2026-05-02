package repository

import "github.com/ak-repo/wim/internal/db"

type Repositories struct {
	User           UserRepository
	Auth           AuthRepository
	Customer       CustomerRepository
	Product        ProductRepository
	Warehouse      WarehouseRepository
	Location       LocationRepository
	RefCode        RefCodeGenerator
	CustomerType   CustomerTypeRepository
	UserRole       UserRoleRepository
	ProductCategory ProductCategoryRepository
}

type Dependencies struct {
	DB    *db.DB
	Redis *db.Redis
}

func NewRepositories(deps Dependencies) *Repositories {
	return &Repositories{
		User:            NewUserRepository(deps.DB),
		Auth:            NewAuthRepository(deps.DB),
		Customer:        NewCustomerRepository(deps.DB),
		Product:         NewProductRepository(deps.DB),
		Warehouse:       NewWarehouseRepository(deps.DB),
		Location:        NewLocationRepository(deps.DB),
		RefCode:         NewRefCodeGenerator(deps.DB, deps.Redis),
		CustomerType:    NewCustomerTypeRepository(deps.DB),
		UserRole:        NewUserRoleRepository(deps.DB),
		ProductCategory: NewProductCategoryRepository(deps.DB),
	}
}
