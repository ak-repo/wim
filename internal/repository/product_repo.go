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

var ErrProductNotFound = errors.New("product not found")

type ProductRepository interface {
	Create(ctx context.Context, product *model.ProductRequest) (int, error)
	GetByID(ctx context.Context, productID int) (*model.ProductDTO, error)
	GetBySKU(ctx context.Context, sku string) (*model.ProductDTO, error)
	ExistsBySKU(ctx context.Context, sku string) (bool, error)
	Update(ctx context.Context, productID int, product *model.ProductRequest) error
	Delete(ctx context.Context, productID int) error
	List(ctx context.Context, params *model.ProductParams) (model.ProductDTOs, error)
	Count(ctx context.Context, params *model.ProductParams) (int, error)
}

type productRepository struct {
	db *db.DB
}

func NewProductRepository(database *db.DB) ProductRepository {
	return &productRepository{db: database}
}

// CREATE
func (r *productRepository) Create(ctx context.Context, product *model.ProductRequest) (int, error) {
	query := `
		INSERT INTO products (
			ref_code, sku, name, description, category, unit_of_measure,
			weight, length, width, height, barcode, is_active, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW())
		RETURNING id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query,
		product.RefCode,
		product.SKU,
		product.Name,
		product.Description,
		product.Category,
		product.UnitOfMeasure,
		product.Weight,
		product.Length,
		product.Width,
		product.Height,
		product.Barcode,
		product.IsActive,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, apperrors.ErrAlreadyExists
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create product")
	}

	return id, nil
}

// GET BY ID
func (r *productRepository) GetByID(ctx context.Context, productID int) (*model.ProductDTO, error) {
	query := `
		SELECT id, ref_code, sku, name, description, category, unit_of_measure,
		       weight, length, width, height, barcode, is_active, created_at, updated_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	return scanProduct(ctx, r.db, query, productID)
}

// GET BY SKU
func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*model.ProductDTO, error) {
	query := `
		SELECT id, ref_code, sku, name, description, category, unit_of_measure,
		       weight, length, width, height, barcode, is_active, created_at, updated_at
		FROM products
		WHERE sku = $1 AND deleted_at IS NULL
	`

	return scanProduct(ctx, r.db, query, sku)
}

// EXISTS
func (r *productRepository) ExistsBySKU(ctx context.Context, sku string) (bool, error) {
	var exists bool

	err := r.db.Pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM products WHERE sku = $1 AND deleted_at IS NULL)`,
		sku,
	).Scan(&exists)

	if err != nil {
		return false, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check sku existence")
	}

	return exists, nil
}

// UPDATE (PATCH STYLE)
func (r *productRepository) Update(ctx context.Context, productID int, product *model.ProductRequest) error {
	query := `
		UPDATE products
		SET name = COALESCE($2, name),
			description = COALESCE($3, description),
			category = COALESCE($4, category),
			unit_of_measure = COALESCE($5, unit_of_measure),
			weight = COALESCE($6, weight),
			length = COALESCE($7, length),
			width = COALESCE($8, width),
			height = COALESCE($9, height),
			barcode = COALESCE($10, barcode),
			is_active = COALESCE($11, is_active),
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		productID,
		product.Name,
		product.Description,
		product.Category,
		product.UnitOfMeasure,
		product.Weight,
		product.Length,
		product.Width,
		product.Height,
		product.Barcode,
		product.IsActive,
	)

	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update product")
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

// DELETE
func (r *productRepository) Delete(ctx context.Context, productID int) error {
	result, err := r.db.Pool.Exec(ctx,
		`UPDATE products SET deleted_at = NOW() WHERE id = $1`,
		productID,
	)

	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete product")
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

// LIST
func (r *productRepository) List(ctx context.Context, params *model.ProductParams) (model.ProductDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, ref_code, sku, name, description, category, unit_of_measure,
		       weight, length, width, height, barcode, is_active, created_at, updated_at
		FROM products
	`

	// Base condition: only get non-deleted records
	conditions = append(conditions, "deleted_at IS NULL")

	// filters
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if params.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", len(args)+1))
		args = append(args, params.Category)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// pagination
	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		len(args)+1, len(args)+2)

	args = append(args, params.Limit, offset)

	return scanProducts(ctx, r.db, query, args...)
}

// COUNT
func (r *productRepository) Count(ctx context.Context, params *model.ProductParams) (int, error) {
	var (
		count      int
		args       []any
		conditions []string
	)

	query := `SELECT COUNT(*) FROM products`

	// Base condition: only count non-deleted records
	conditions = append(conditions, "deleted_at IS NULL")

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if params.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", len(args)+1))
		args = append(args, params.Category)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count products")
	}

	return count, nil
}

// SCAN SINGLE
func scanProduct(ctx context.Context, database *db.DB, query string, args ...any) (*model.ProductDTO, error) {
	var row model.ProductDTO
	var isActive sql.NullBool
	var createdAt, updatedAt sql.NullTime

	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.RefCode,
		&row.SKU,
		&row.Name,
		&row.Description,
		&row.Category,
		&row.UnitOfMeasure,
		&row.Weight,
		&row.Length,
		&row.Width,
		&row.Height,
		&row.Barcode,
		&isActive,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product")
	}
	row.ApplyProductNullScalars(isActive, createdAt, updatedAt)

	return &row, nil
}

// SCAN LIST
func scanProducts(ctx context.Context, database *db.DB, query string, args ...any) (model.ProductDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list products")
	}
	defer rows.Close()

	var products model.ProductDTOs

	for rows.Next() {
		var row model.ProductDTO
		var isActive sql.NullBool
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&row.ID,
			&row.RefCode,
			&row.SKU,
			&row.Name,
			&row.Description,
			&row.Category,
			&row.UnitOfMeasure,
			&row.Weight,
			&row.Length,
			&row.Width,
			&row.Height,
			&row.Barcode,
			&isActive,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan product row")
		}
		row.ApplyProductNullScalars(isActive, createdAt, updatedAt)

		products = append(products, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate products")
	}

	return products, nil
}
