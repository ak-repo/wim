package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/ak-repo/wim/internal/http/handler"
	"github.com/ak-repo/wim/pkg/auth"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func SetupRoutes(handlers *handler.Handler, tokenManager auth.TokenManager) http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)

	allowedOrigins := []string{"http://localhost:5174", "http://localhost:3050"}
	if raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); raw != "" {
		parts := strings.Split(raw, ",")
		allowedOrigins = allowedOrigins[:0]
		for _, part := range parts {
			origin := strings.TrimSpace(part)
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", handlers.Health.Check)

	apiPrefix := strings.TrimSpace(os.Getenv("WIM_API_PREFIX"))
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
