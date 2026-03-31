package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/model"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUserNotFound = errors.New("user not found")
var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type UserRepository interface {
	Create(ctx context.Context, user *model.UserRequest) error
	GetByID(ctx context.Context, userID uuid.UUID) (*model.UserDTO, error)
	GetByEmail(ctx context.Context, email string) (*model.UserDTO, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *model.UserRequest) error
	Delete(ctx context.Context, userID uuid.UUID) error

	List(ctx context.Context, params *model.UserParams) (model.UserDTOs, error)
	Count(ctx context.Context, params *model.UserParams) (int, error)
}

type userRepository struct {
	db *db.DB
}

func NewUserRepository(database *db.DB) UserRepository {
	return &userRepository{db: database}
}

func (r *userRepository) Create(ctx context.Context, user *model.UserRequest) error {

	query := `
		INSERT INTO users (
			id, username, email, password_hash, role, contact, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Pool.Exec(ctx, query, user.ID, user.Username, user.Email, user.PasswordHash, user.Role, user.Contact, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userID uuid.UUID) (*model.UserDTO, error) {
	row, err := scanUser(ctx, r.db, `
		SELECT id, username, email, password_hash, role, contact, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.UserDTO, error) {
	row, err := scanUser(ctx, r.db, `
		SELECT id, username, email, password_hash, role, contact, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`, email)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user by email: %w", err)
	}

	return exists, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.UserRequest) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, password_hash = $4, role = $5, contact = $6, is_active = $7, updated_at = $8
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, user.ID, user.Username, user.Email, user.PasswordHash, user.Role, user.Contact, user.IsActive, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	result, err := r.db.Pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func scanUser(ctx context.Context, database *db.DB, query string, args ...any) (*model.UserDTO, error) {
	var row model.UserDTO
	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.Username,
		&row.Email,
		&row.PasswordHash,
		&row.Role,
		&row.Contact,
		&row.IsActive,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.UserDTO{}, ErrUserNotFound
		}
		return &model.UserDTO{}, fmt.Errorf("scan user: %w", err)
	}

	return &row, nil
}

func (r *userRepository) List(ctx context.Context, params *model.UserParams) (model.UserDTOs, error) {
	args := []interface{}{}
	conditions := []string{}
	query := `
		SELECT id, username, email, password_hash, role, contact, is_active, created_at, updated_at
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

	rows, err := scanUsers(ctx, r.db, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return rows, nil
}

func (r *userRepository) Count(ctx context.Context, params *model.UserParams) (int, error) {
	var count int
	args := []interface{}{}
	conditions := []string{}

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

	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}

func scanUsers(ctx context.Context, database *db.DB, query string, args ...any) (model.UserDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scan users: %w", err)
	}
	defer rows.Close()

	var users model.UserDTOs
	for rows.Next() {
		var row model.UserDTO
		if err := rows.Scan(
			&row.ID,
			&row.Username,
			&row.Email,
			&row.PasswordHash,
			&row.Role,
			&row.Contact,
			&row.IsActive,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}
		users = append(users, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user rows: %w", err)
	}

	return users, nil
}
