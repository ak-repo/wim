package router

import (
	"net/http"

	"github.com/ak-repo/wim/internal/http/handler"
	"github.com/ak-repo/wim/internal/http/middleware"
	"github.com/go-chi/chi"
)

func SetupRoutes(handlers *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/api/v1/health", handlers.Health.Check)
	r.Post("/api/v1/auth/register", handlers.Auth.Register)
	r.Post("/api/v1/auth/login", handlers.Auth.Login)
	r.With(middleware.RequireAuth).Get("/api/v1/me", handlers.Auth.Me)

	return r
}
