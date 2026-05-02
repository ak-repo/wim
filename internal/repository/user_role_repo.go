package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/model"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUserRoleNotFound = errors.New("user role not found")

type UserRoleRepository interface {
	List(ctx context.Context, params *model.UserRoleParams) (model.UserRoleDTOs, error)
	Count(ctx context.Context, params *model.UserRoleParams) (int, error)
	GetByID(ctx context.Context, id int) (*model.UserRoleDTO, error)
	Create(ctx context.Context, req *model.UserRoleRequest) (int, error)
	Update(ctx context.Context, id int, req *model.UserRoleRequest) error
	Delete(ctx context.Context, id int) error
}

type userRoleRepository struct {
	db *db.DB
}

func NewUserRoleRepository(db *db.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) List(ctx context.Context, params *model.UserRoleParams) (model.UserRoleDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, name, description, is_active, created_at, updated_at, deleted_at
		FROM user_roles
	`

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, params.Limit, offset)

	return scanUserRoles(ctx, r.db, query, args...)
}

func (r *userRoleRepository) Count(ctx context.Context, params *model.UserRoleParams) (int, error) {
	var (
		count      int
		args       []any
		conditions []string
	)

	query := `SELECT COUNT(*) FROM user_roles`

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count user roles")
	}

	return count, nil
}

func (r *userRoleRepository) GetByID(ctx context.Context, id int) (*model.UserRoleDTO, error) {
	return scanUserRole(ctx, r.db, `
		SELECT id, name, description, is_active, created_at, updated_at, deleted_at
		FROM user_roles WHERE id = $1
	`, id)
}

func (r *userRoleRepository) Create(ctx context.Context, req *model.UserRoleRequest) (int, error) {
	query := `
		INSERT INTO user_roles (name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query,
		req.Name,
		req.Description,
		req.IsActive,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, apperrors.ErrAlreadyExists
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create user role")
	}

	return id, nil
}

func (r *userRoleRepository) Update(ctx context.Context, id int, req *model.UserRoleRequest) error {
	query := `
		UPDATE user_roles
		SET name = COALESCE($2, name),
			description = COALESCE($3, description),
			is_active = COALESCE($4, is_active),
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		id,
		req.Name,
		req.Description,
		req.IsActive,
	)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update user role")
	}

	if result.RowsAffected() == 0 {
		return ErrUserRoleNotFound
	}

	return nil
}

func (r *userRoleRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.Pool.Exec(ctx, `DELETE FROM user_roles WHERE id = $1`, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete user role")
	}

	if result.RowsAffected() == 0 {
		return ErrUserRoleNotFound
	}

	return nil
}

func scanUserRole(ctx context.Context, database *db.DB, query string, args ...any) (*model.UserRoleDTO, error) {
	var row model.UserRoleDTO
	var isActive sql.NullBool
	var createdAt, updatedAt, deletedAt sql.NullTime

	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.Name,
		&row.Description,
		&isActive,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserRoleNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load user role")
	}
	row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)

	return &row, nil
}

func scanUserRoles(ctx context.Context, database *db.DB, query string, args ...any) (model.UserRoleDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list user roles")
	}
	defer rows.Close()

	var userRoles model.UserRoleDTOs

	for rows.Next() {
		var row model.UserRoleDTO
		var isActive sql.NullBool
		var createdAt, updatedAt, deletedAt sql.NullTime

		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Description,
			&isActive,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan user role row")
		}
		row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)

		userRoles = append(userRoles, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate user roles")
	}

	return userRoles, nil
}