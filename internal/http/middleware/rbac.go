package middleware

import (
	"net/http"

	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/pkg/auth"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/ak-repo/wim/pkg/response"
)

// Role Based Access Controlw
func RoleBasedAccessControl(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := auth.ClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if claims.Role == constants.RoleSuperAdmin {
				next.ServeHTTP(w, r)
				return
			}
			if claims.Role != role {
				response.WriteError(w, http.StatusForbidden, apperrors.CodeForbidden, "you are not authorized to access this resource")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
