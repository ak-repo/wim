package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	Subject uuid.UUID
	Role    string
}

type TokenManager interface {
	IssueAccessToken(ctx context.Context, claims Claims) (string, error)
	IssueRefreshToken(ctx context.Context, claims Claims) (string, time.Duration, error)
	ParseJWTToken(ctx context.Context, token string) (Claims, error)
}

type JWTTokenManager struct {
	secretKey []byte
	issuer    string
	ttl       time.Duration
}

type accessTokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTTokenManager(secretKey, issuer string, ttl time.Duration) TokenManager {
	if ttl <= 0 {
		ttl = 60 * time.Minute
	}

	return &JWTTokenManager{
		secretKey: []byte(secretKey),
		issuer:    issuer,
		ttl:       ttl,
	}
}

func (m *JWTTokenManager) IssueAccessToken(ctx context.Context, claims Claims) (string, error) {
	now := time.Now().UTC()
	jwtClaims := accessTokenClaims{
		UserID: claims.Subject.String(),
		Role:   claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   claims.Subject.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(m.secretKey)
}

func (m *JWTTokenManager) IssueRefreshToken(ctx context.Context, claims Claims) (string, time.Duration, error) {
	now := time.Now().UTC()
	refreshTime := m.ttl * 24 // 1 Day
	jwtClaims := accessTokenClaims{
		Role: string(claims.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   claims.Subject.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenStr, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return tokenStr, refreshTime, nil
}

func (m *JWTTokenManager) ParseJWTToken(ctx context.Context, token string) (Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &accessTokenClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}

		return m.secretKey, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("parse access token: %w", ErrInvalidToken)
	}

	claims, ok := parsed.Claims.(*accessTokenClaims)
	if !ok || !parsed.Valid {
		return Claims{}, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return Claims{}, fmt.Errorf("parse access token subject: %w", ErrInvalidToken)
	}

	return Claims{Subject: userID, Role: claims.Role}, nil
}
