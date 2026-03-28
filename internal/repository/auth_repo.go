package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	StoreRefreshToken(ctx context.Context, token *model.RefreshTokenDTO) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshTokenDTO, error)
	RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID, revokedAt time.Time) error
}

type authRepository struct {
	db *db.DB
}

func NewAuthRepository(database *db.DB) AuthRepository {
	return &authRepository{db: database}
}

func (r *authRepository) StoreRefreshToken(ctx context.Context, token *model.RefreshTokenDTO) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO refresh_tokens (
			id, user_id, token_hash, expires_at, revoked_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.RevokedAt, token.CreatedAt, token.UpdatedAt)
	if err != nil {
		return fmt.Errorf("store refresh token: %w", err)
	}

	return nil
}

func (r *authRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshTokenDTO, error) {
	var row model.RefreshTokenDTO
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at, updated_at
		FROM refresh_tokens WHERE token_hash = $1
	`, tokenHash).Scan(
		&row.ID,
		&row.UserID,
		&row.TokenHash,
		&row.ExpiresAt,
		&row.RevokedAt,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("get refresh token by hash: %w", err)
	}

	return &row, nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID, revokedAt time.Time) error {
	result, err := r.db.Pool.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = $2, updated_at = $2
		WHERE id = $1 AND revoked_at IS NULL
	`, tokenID, revokedAt)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}
