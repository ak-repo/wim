package handler

import "github.com/ak-repo/wim/internal/service"

type Handler struct {
	Auth   *AuthHandler
	Health *HealthHandler
	User   *UserHandler
}

func NewHandlers(services *service.Services) *Handler {
	return &Handler{
		Auth:   NewAuthHandler(services),
		Health: NewHealthHandler(),
		User:   NewUserHandler(services),
	}
}
