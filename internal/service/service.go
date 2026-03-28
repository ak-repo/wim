package service

import (
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/pkg/auth"
)

type Services struct {
	User UserService
	Auth AuthService
}

type Dependencies struct {
	Repositories   *repository.Repositories
	PasswordHasher auth.PasswordHasher
	TokenManager   auth.TokenManager
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		User: NewUserService(deps.Repositories),
		Auth: NewAuthService(deps.Repositories, deps.TokenManager, deps.PasswordHasher),
	}
}
