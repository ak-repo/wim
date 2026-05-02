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

var ErrProductCategoryNotFound = errors.New("product category not found")

type ProductCategoryRepository interface {
	List(ctx context.Context, params *model.ProductCategoryParams) (model.ProductCategoryDTOs, error)
	Count(ctx context.Context, params *model.ProductCategoryParams) (int, error)
	GetByID(ctx context.Context, id int) (*model.ProductCategoryDTO, error)
	Create(ctx context.Context, req *model.ProductCategoryRequest) (int, error)
	Update(ctx context.Context, id int, req *model.ProductCategoryRequest) error
	Delete(ctx context.Context, id int) error
}

type productCategoryRepository struct {
	db *db.DB
}

func NewProductCategoryRepository(db *db.DB) ProductCategoryRepository {
	return &productCategoryRepository{db: db}
}

func (r *productCategoryRepository) List(ctx context.Context, params *model.ProductCategoryParams) (model.ProductCategoryDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, name, description, is_active, created_at, updated_at, deleted_at
		FROM product_categories
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

	return scanProductCategories(ctx, r.db, query, args...)
}

func (r *productCategoryRepository) Count(ctx context.Context, params *model.ProductCategoryParams) (int, error) {
	var (
		count      int
		args       []any
		conditions []string
	)

	query := `SELECT COUNT(*) FROM product_categories`

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count product categories")
	}

	return count, nil
}

func (r *productCategoryRepository) GetByID(ctx context.Context, id int) (*model.ProductCategoryDTO, error) {
	return scanProductCategory(ctx, r.db, `
		SELECT id, name, description, is_active, created_at, updated_at, deleted_at
		FROM product_categories WHERE id = $1
	`, id)
}

func (r *productCategoryRepository) Create(ctx context.Context, req *model.ProductCategoryRequest) (int, error) {
	query := `
		INSERT INTO product_categories (name, description, is_active, created_at, updated_at)
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
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create product category")
	}

	return id, nil
}

func (r *productCategoryRepository) Update(ctx context.Context, id int, req *model.ProductCategoryRequest) error {
	query := `
		UPDATE product_categories
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
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update product category")
	}

	if result.RowsAffected() == 0 {
		return ErrProductCategoryNotFound
	}

	return nil
}

func (r *productCategoryRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.Pool.Exec(ctx, `DELETE FROM product_categories WHERE id = $1`, id)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete product category")
	}

	if result.RowsAffected() == 0 {
		return ErrProductCategoryNotFound
	}

	return nil
}

func scanProductCategory(ctx context.Context, database *db.DB, query string, args ...any) (*model.ProductCategoryDTO, error) {
	var row model.ProductCategoryDTO
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
			return nil, ErrProductCategoryNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product category")
	}
	row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)

	return &row, nil
}

func scanProductCategories(ctx context.Context, database *db.DB, query string, args ...any) (model.ProductCategoryDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list product categories")
	}
	defer rows.Close()

	var productCategories model.ProductCategoryDTOs

	for rows.Next() {
		var row model.ProductCategoryDTO
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
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan product category row")
		}
		row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)

		productCategories = append(productCategories, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate product categories")
	}

	return productCategories, nil
}