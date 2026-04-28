package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ak-repo/wim/internal/db"
	apperrors "github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/model"
	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	StoreRefreshToken(ctx context.Context, token *model.RefreshTokenDTO) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshTokenDTO, error)
	RevokeRefreshToken(ctx context.Context, tokenID int, revokedAt time.Time) error
}

type authRepository struct {
	db *db.DB
}

func NewAuthRepository(database *db.DB) AuthRepository {
	return &authRepository{db: database}
}

func (r *authRepository) StoreRefreshToken(ctx context.Context, token *model.RefreshTokenDTO) error {
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO refresh_tokens (
			ref_code, user_id, token_hash, expires_at, revoked_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, token.RefCode, token.UserID, token.TokenHash, token.ExpiresAt, token.RevokedAt, token.CreatedAt, token.UpdatedAt).Scan(&token.ID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to store refresh token")
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
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load refresh token")
	}

	return &row, nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, tokenID int, revokedAt time.Time) error {
	result, err := r.db.Pool.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = $2, updated_at = $2
		WHERE id = $1 AND revoked_at IS NULL
	`, tokenID, revokedAt)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to revoke refresh token")
	}

	if result.RowsAffected() == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}
