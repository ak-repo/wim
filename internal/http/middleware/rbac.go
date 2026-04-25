package middleware

import (
	"errors"
	"net/http"

	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/httpx"
	"github.com/ak-repo/wim/pkg/auth"
)

const opRBAC errs.Op = "middleware/RBAC.RoleBasedAccessControl"

// Role Based Access Control
func RoleBasedAccessControl(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := auth.ClaimsFromContext(r.Context())
			if !ok {
				err := errs.E(opRBAC, errs.Unauthorized, errors.New("unauthorized"), errs.WithCode(errs.CodeUnauthorized))
				httpx.WriteError(w, r, err)
				return
			}

			if claims.Role == constants.RoleSuperAdmin {
				next.ServeHTTP(w, r)
				return
			}
			if claims.Role != role {
				err := errs.E(opRBAC, errs.Forbidden, errors.New("you are not authorized to access this resource"), errs.WithCode(errs.CodeForbidden))
				httpx.WriteError(w, r, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
