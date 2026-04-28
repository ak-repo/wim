package auth

import "context"

type claimsContextKey string

const claimsKey claimsContextKey = "auth.claims"

func WithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(Claims)
	return claims, ok
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return 0, false
	}
	return claims.Subject, true
}

func RoleFromContext(ctx context.Context) (string, bool) {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return "", false
	}
	return claims.Role, true
}
