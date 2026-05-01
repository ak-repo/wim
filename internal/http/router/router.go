package router

import (
	"net/http"
	"strings"

	"github.com/ak-repo/wim/config"
	"github.com/ak-repo/wim/internal/http/handler"
	"github.com/ak-repo/wim/pkg/auth"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func SetupRoutes(handlers *handler.Handler, tokenManager auth.TokenManager, cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)

	allowedOrigins := cfg.Server.AllowOrigins
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", handlers.Health.Check)

	apiPrefix := strings.TrimSpace(cfg.Server.APIPrefix)
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}
	if !strings.HasPrefix(apiPrefix, "/") {
		apiPrefix = "/" + apiPrefix
	}

	r.Route(apiPrefix, func(api chi.Router) {
		AdminRoutes(api, handlers, tokenManager)
		// api.With(wimMiddleware.RequireAuth(tokenManager)).Get("/me", handlers.Auth.Me)
	})

	return r
}
