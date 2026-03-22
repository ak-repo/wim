package handler

import "github.com/ak-repo/wim/internal/service"

type Handler struct {
	Auth   *AuthHandler
	Health *HealthHandler
}

func NewHandlers(services *service.Service) *Handler {
	return &Handler{
		Auth:   NewAuthHandler(services.User),
		Health: NewHealthHandler(),
	}
}
