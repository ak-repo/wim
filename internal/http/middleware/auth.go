package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/httpx"
	"github.com/ak-repo/wim/pkg/auth"
)

const opAuth errs.Op = "middleware/Auth.RequireAuth"

func RequireAuth(tokenManager auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := strings.TrimSpace(r.Header.Get("Authorization"))
			if authorization == "" {
				err := errs.E(opAuth, errs.Unauthorized, errors.New("missing authorization header"), errs.WithCode(errs.CodeUnauthorized))
				httpx.WriteError(w, r, err)
				return
			}

			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				err := errs.E(opAuth, errs.Unauthorized, errors.New("invalid authorization header"), errs.WithCode(errs.CodeUnauthorized))
				httpx.WriteError(w, r, err)
				return
			}

			claims, err := tokenManager.ParseJWTToken(r.Context(), strings.TrimSpace(parts[1]))
			if err != nil {
				err := errs.E(opAuth, errs.Unauthorized, errors.New("invalid or expired access token"), errs.WithCode(errs.CodeUnauthorized))
				httpx.WriteError(w, r, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithClaims(r.Context(), claims)))
		})
	}
}
