package repository

import (
	"context"
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
	Create(ctx context.Context, role *model.UserRoleRequest) (int, error)
	GetByID(ctx context.Context, roleID int) (*model.UserRoleDTO, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Update(ctx context.Context, roleID int, role *model.UserRoleRequest) error
	Delete(ctx context.Context, roleID int) error
	List(ctx context.Context, params *model.UserRoleParams) (model.UserRoleDTOs, error)
	Count(ctx context.Context, params *model.UserRoleParams) (int, error)
}

type userRoleRepository struct {
	db *db.DB
}

func NewUserRoleRepository(database *db.DB) UserRoleRepository {
	return &userRoleRepository{db: database}
}

func (r *userRoleRepository) Create(ctx context.Context, role *model.UserRoleRequest) (int, error) {
	query := `
		INSERT INTO user_roles (
			ref_code, name, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query,
		role.RefCode,
		role.Name,
		role.IsActive,
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

func (r *userRoleRepository) GetByID(ctx context.Context, roleID int) (*model.UserRoleDTO, error) {
	query := `
		SELECT id, name, ref_code, is_active
		FROM user_roles
		WHERE id = $1 AND deleted_at IS NULL
	`

	var row model.UserRoleDTO
	err := r.db.Pool.QueryRow(ctx, query, roleID).Scan(
		&row.ID,
		&row.Name,
		&row.RefCode,
		&row.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserRoleNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load user role")
	}

	return &row, nil
}

func (r *userRoleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_roles WHERE name = $1 AND deleted_at IS NULL)`
	if err := r.db.Pool.QueryRow(ctx, query, name).Scan(&exists); err != nil {
		return false, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check user role by name")
	}
	return exists, nil
}

func (r *userRoleRepository) Update(ctx context.Context, roleID int, role *model.UserRoleRequest) error {
	query := `
		UPDATE user_roles
		SET name = COALESCE($2, name),
			is_active = COALESCE($3, is_active),
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		roleID,
		role.Name,
		role.IsActive,
	)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update user role")
	}

	if result.RowsAffected() == 0 {
		return ErrUserRoleNotFound
	}

	return nil
}

func (r *userRoleRepository) Delete(ctx context.Context, roleID int) error {
	result, err := r.db.Pool.Exec(ctx,
		`UPDATE user_roles SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		roleID,
	)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete user role")
	}
	if result.RowsAffected() == 0 {
		return ErrUserRoleNotFound
	}
	return nil
}

func (r *userRoleRepository) List(ctx context.Context, params *model.UserRoleParams) (model.UserRoleDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, name, ref_code, is_active
		FROM user_roles
	`

	conditions = append(conditions, "deleted_at IS NULL")

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		len(args)+1, len(args)+2,
	)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list user roles")
	}
	defer rows.Close()

	var roles model.UserRoleDTOs
	for rows.Next() {
		var row model.UserRoleDTO
		err := rows.Scan(&row.ID, &row.Name, &row.RefCode, &row.IsActive)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan user role row")
		}
		roles = append(roles, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate user roles")
	}

	return roles, nil
}

func (r *userRoleRepository) Count(ctx context.Context, params *model.UserRoleParams) (int, error) {
	var count int
	var args []any
	var conditions []string

	query := `SELECT COUNT(*) FROM user_roles`

	conditions = append(conditions, "deleted_at IS NULL")

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count user roles")
	}

	return count, nil
}
