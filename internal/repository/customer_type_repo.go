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

var ErrCustomerTypeNotFound = errors.New("customer type not found")

type CustomerTypeRepository interface {
	List(ctx context.Context, params *model.CustomerTypeParams) (model.CustomerTypeDTOs, error)
	Count(ctx context.Context, params *model.CustomerTypeParams) (int, error)
	GetByID(ctx context.Context, id int) (*model.CustomerTypeDTO, error)
	Create(ctx context.Context, req *model.CustomerTypeRequest) (int, error)
	Update(ctx context.Context, id int, req *model.CustomerTypeRequest) error
	Delete(ctx context.Context, id int) error
}

type customerTypeRepository struct {
	db *db.DB
}

func NewCustomerTypeRepository(db *db.DB) CustomerTypeRepository {
	return &customerTypeRepository{db: db}
}

func (r *customerTypeRepository) List(ctx context.Context, params *model.CustomerTypeParams) (model.CustomerTypeDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, name, description, is_active, created_at, updated_at, deleted_at
		FROM customer_types
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

	return scanCustomerTypes(ctx, r.db, query, args...)
}

func (r *customerTypeRepository) Count(ctx context.Context, params *model.CustomerTypeParams) (int, error) {
	var (
		count      int
		args       []any
		conditions []string
	)

	query := `SELECT COUNT(*) FROM customer_types`

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count customer types")
	}

	return count, nil
}

func (r *customerTypeRepository) GetByID(ctx context.Context, id int) (*model.CustomerTypeDTO, error) {
	return scanCustomerType(ctx, r.db, `
		SELECT id, name, description, is_active, created_at, updated_at, deleted_at
		FROM customer_types WHERE id = $1
	`, id)
}

func (r *customerTypeRepository) Create(ctx context.Context, req *model.CustomerTypeRequest) (int, error) {
	query := `
		INSERT INTO customer_types (name, description, is_active, created_at, updated_at)
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
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create customer type")
	}

	return id, nil
}

func (r *customerTypeRepository) Update(ctx context.Context, id int, req *model.CustomerTypeRequest) error {
	query := `
		UPDATE customer_types
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
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update customer type")
	}

	if result.RowsAffected() == 0 {
		return ErrCustomerTypeNotFound
	}

	return nil
}

func (r *customerTypeRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.Pool.Exec(ctx, `DELETE FROM customer_types WHERE id = $1`, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete customer type")
	}

	if result.RowsAffected() == 0 {
		return ErrCustomerTypeNotFound
	}

	return nil
}

func scanCustomerType(ctx context.Context, database *db.DB, query string, args ...any) (*model.CustomerTypeDTO, error) {
	var row model.CustomerTypeDTO
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
			return nil, ErrCustomerTypeNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load customer type")
	}
	row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)

	return &row, nil
}

func scanCustomerTypes(ctx context.Context, database *db.DB, query string, args ...any) (model.CustomerTypeDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list customer types")
	}
	defer rows.Close()

	var customerTypes model.CustomerTypeDTOs

	for rows.Next() {
		var row model.CustomerTypeDTO
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
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan customer type row")
		}
		row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)

		customerTypes = append(customerTypes, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate customer types")
	}

	return customerTypes, nil
}
