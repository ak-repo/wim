package service

import (
	"github.com/ak-repo/wim/internal/repository"
)

type Service struct {
	User UserService
}

type Dependencies struct {
	Repositories *repository.Repositories
}

func NewServices(deps Dependencies) *Service {
	return &Service{
		User: NewUserService(deps.Repositories.User),
	}
}
