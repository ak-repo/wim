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

var ErrProductCategoryNotFound = errors.New("product category not found")

type ProductCategoryRepository interface {
	Create(ctx context.Context, category *model.ProductCategoryRequest) (int, error)
	GetByID(ctx context.Context, categoryID int) (*model.ProductCategoryDTO, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Update(ctx context.Context, categoryID int, category *model.ProductCategoryRequest) error
	Delete(ctx context.Context, categoryID int) error
	List(ctx context.Context, params *model.ProductCategoryParams) (model.ProductCategoryDTOs, error)
	Count(ctx context.Context, params *model.ProductCategoryParams) (int, error)
}

type productCategoryRepository struct {
	db *db.DB
}

func NewProductCategoryRepository(database *db.DB) ProductCategoryRepository {
	return &productCategoryRepository{db: database}
}

func (r *productCategoryRepository) Create(ctx context.Context, category *model.ProductCategoryRequest) (int, error) {
	query := `
		INSERT INTO product_categories (
			ref_code, name, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query,
		category.RefCode,
		category.Name,
		category.IsActive,
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

func (r *productCategoryRepository) GetByID(ctx context.Context, categoryID int) (*model.ProductCategoryDTO, error) {
	query := `
		SELECT id, name, ref_code, is_active
		FROM product_categories
		WHERE id = $1 AND deleted_at IS NULL
	`

	var row model.ProductCategoryDTO
	err := r.db.Pool.QueryRow(ctx, query, categoryID).Scan(
		&row.ID,
		&row.Name,
		&row.RefCode,
		&row.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductCategoryNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product category")
	}

	return &row, nil
}

func (r *productCategoryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM product_categories WHERE name = $1 AND deleted_at IS NULL)`
	if err := r.db.Pool.QueryRow(ctx, query, name).Scan(&exists); err != nil {
		return false, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check product category by name")
	}
	return exists, nil
}

func (r *productCategoryRepository) Update(ctx context.Context, categoryID int, category *model.ProductCategoryRequest) error {
	query := `
		UPDATE product_categories
		SET name = COALESCE($2, name),
			is_active = COALESCE($3, is_active),
			updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.db.Pool.Exec(ctx, query,
		categoryID,
		category.Name,
		category.IsActive,
	)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update product category")
	}

	if result.RowsAffected() == 0 {
		return ErrProductCategoryNotFound
	}

	return nil
}

func (r *productCategoryRepository) Delete(ctx context.Context, categoryID int) error {
	result, err := r.db.Pool.Exec(ctx,
		`UPDATE product_categories SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		categoryID,
	)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete product category")
	}
	if result.RowsAffected() == 0 {
		return ErrProductCategoryNotFound
	}
	return nil
}

func (r *productCategoryRepository) List(ctx context.Context, params *model.ProductCategoryParams) (model.ProductCategoryDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, name, ref_code, is_active
		FROM product_categories
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
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list product categories")
	}
	defer rows.Close()

	var categories model.ProductCategoryDTOs
	for rows.Next() {
		var row model.ProductCategoryDTO
		err := rows.Scan(&row.ID, &row.Name, &row.RefCode, &row.IsActive)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan product category row")
		}
		categories = append(categories, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate product categories")
	}

	return categories, nil
}

func (r *productCategoryRepository) Count(ctx context.Context, params *model.ProductCategoryParams) (int, error) {
	var count int
	var args []any
	var conditions []string

	query := `SELECT COUNT(*) FROM product_categories`

	conditions = append(conditions, "deleted_at IS NULL")

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count product categories")
	}

	return count, nil
}
