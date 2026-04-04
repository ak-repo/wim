package model

import (
	"database/sql"
	"time"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
	Contact  string `json:"contact,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	UserID       int          `db:"user_id"`
	TokenHash    string       `db:"token_hash"`
	ExpiresAt    time.Time    `db:"expires_at"`
	RevokedAt    sql.NullTime `db:"revoked_at"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	RefreshToken string       `json:"refreshToken"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RefreshTokenDTO struct {
	ID        int          `db:"id"`
	UserID    int          `db:"user_id"`
	TokenHash string       `db:"token_hash"`
	ExpiresAt time.Time    `db:"expires_at"`
	RevokedAt sql.NullTime `db:"revoked_at"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}
