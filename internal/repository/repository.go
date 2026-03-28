package repository

import "github.com/ak-repo/wim/internal/db"

type Repositories struct {
	User UserRepository
	Auth AuthRepository
}

type Dependencies struct {
	DB    *db.DB
	Redis *db.Redis
}

func NewRepositories(deps Dependencies) *Repositories {
	return &Repositories{
		User: NewUserRepository(deps.DB),
		Auth: NewAuthRepository(deps.DB),
	}
}
