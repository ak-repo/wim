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

var ErrUserNotFound = errors.New("user not found")
var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type UserRepository interface {
	Create(ctx context.Context, user *model.UserRequest) (int, error)
	GetByID(ctx context.Context, userID int) (*model.UserDTO, error)
	GetByEmail(ctx context.Context, email string) (*model.UserDTO, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	Update(ctx context.Context, userID int, user *model.UserRequest) error
	Delete(ctx context.Context, userID int) error

	List(ctx context.Context, params *model.UserParams) (model.UserDTOs, error)
	Count(ctx context.Context, params *model.UserParams) (int, error)
}

type userRepository struct {
	db *db.DB
}

func NewUserRepository(database *db.DB) UserRepository {
	return &userRepository{
		db: database,
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.UserRequest) (int, error) {
	query := `
		INSERT INTO users (
			 ref_code, username, email, password_hash, role, contact, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query,
		user.RefCode,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Contact,
		user.IsActive,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, apperrors.ErrAlreadyExists
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create user")
	}

	return id, nil
}

func (r *userRepository) GetByID(ctx context.Context, userID int) (*model.UserDTO, error) {
	return scanUser(ctx, r.db, `
		SELECT id, ref_code, username, email, password_hash, role, contact, is_active, created_at, updated_at, deleted_at
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`, userID)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.UserDTO, error) {
	return scanUser(ctx, r.db, `
		SELECT id, ref_code, username, email, password_hash, role, contact, is_active, created_at, updated_at, deleted_at
		FROM users WHERE email = $1
	`, email)
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return false, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check user by email")
	}

	return exists, nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username).Scan(&exists)
	if err != nil {
		return false, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check user by username")
	}

	return exists, nil
}

func (r *userRepository) Update(ctx context.Context, userID int, user *model.UserRequest) error {
	query := `
		UPDATE users
		SET username = COALESCE($2, username),
			email = COALESCE($3, email),
			password_hash = COALESCE($4, password_hash),
			role = COALESCE($5, role),
			contact = COALESCE($6, contact),
			is_active = COALESCE($7, is_active),
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		userID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Contact,
		user.IsActive,
	)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update user")
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, userID int) error {
	result, err := r.db.Pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete user")
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func scanUser(ctx context.Context, database *db.DB, query string, args ...any) (*model.UserDTO, error) {
	var row model.UserDTO
	var isActive sql.NullBool
	var createdAt, updatedAt, deletedAt sql.NullTime

	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.RefCode,
		&row.Username,
		&row.Email,
		&row.PasswordHash,
		&row.Role,
		&row.Contact,
		&isActive,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err == nil {
		row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load user")
	}

	return &row, nil
}

func (r *userRepository) List(ctx context.Context, params *model.UserParams) (model.UserDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, ref_code, username, email, password_hash, role, contact, is_active, created_at, updated_at, deleted_at
		FROM users
	`

	// Active users filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	// Apply WHERE if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Pagination
	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		len(args)+1, len(args)+2,
	)
	args = append(args, params.Limit, offset)

	return scanUsers(ctx, r.db, query, args...)
}

func (r *userRepository) Count(ctx context.Context, params *model.UserParams) (int, error) {
	var count int
	var args []any
	var conditions []string

	query := `SELECT COUNT(*) FROM users`

	// Active users filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	// Apply WHERE if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count users")
	}

	return count, nil
}

func scanUsers(ctx context.Context, database *db.DB, query string, args ...any) (model.UserDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query users")
	}
	defer rows.Close()

	var users model.UserDTOs
	for rows.Next() {
		var row model.UserDTO
		var isActive sql.NullBool
		var createdAt, updatedAt, deletedAt sql.NullTime

		err := rows.Scan(
			&row.ID,
			&row.RefCode,
			&row.Username,
			&row.Email,
			&row.PasswordHash,
			&row.Role,
			&row.Contact,
			&isActive,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan user row")
		}
		row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)
		users = append(users, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate users")
	}

	return users, nil
}
