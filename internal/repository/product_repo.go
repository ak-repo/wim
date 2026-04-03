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

var ErrProductNotFound = errors.New("product not found")

type ProductRepository interface {
	Create(ctx context.Context, product *model.CreateProductRequest) (int, error)
	GetByID(ctx context.Context, productID int) (*model.ProductDTO, error)
	GetBySKU(ctx context.Context, sku string) (*model.ProductDTO, error)
	ExistsBySKU(ctx context.Context, sku string) (bool, error)
	Update(ctx context.Context, productID int, product *model.UpdateProductRequest) error
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

func (r *productRepository) Create(ctx context.Context, product *model.CreateProductRequest) (int, error) {
	query := `
		INSERT INTO products (
			sku, name, description, category, unit_of_measure, 
			weight, length, width, height, barcode, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, true, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query,
		product.SKU, product.Name, product.Description, product.Category, product.UnitOfMeasure,
		product.Weight, product.Length, product.Width, product.Height, product.Barcode,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, apperrors.ErrAlreadyExists
		}
		return 0, fmt.Errorf("failed to create product: %w", err)
	}

	return id, nil
}

func (r *productRepository) GetByID(ctx context.Context, productID int) (*model.ProductDTO, error) {
	row, err := scanProduct(ctx, r.db, `
		SELECT id, sku, name, description, category, unit_of_measure, 
			weight, length, width, height, barcode, is_active, created_at, updated_at
		FROM products WHERE id = $1
	`, productID)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*model.ProductDTO, error) {
	row, err := scanProduct(ctx, r.db, `
		SELECT id, sku, name, description, category, unit_of_measure, 
			weight, length, width, height, barcode, is_active, created_at, updated_at
		FROM products WHERE sku = $1
	`, sku)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *productRepository) ExistsBySKU(ctx context.Context, sku string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1)`, sku).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check product by sku: %w", err)
	}

	return exists, nil
}

func (r *productRepository) Update(ctx context.Context, productID int, product *model.UpdateProductRequest) error {
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
		return fmt.Errorf("failed to update product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, productID int) error {
	result, err := r.db.Pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

func scanProduct(ctx context.Context, database *db.DB, query string, args ...any) (*model.ProductDTO, error) {
	var row model.ProductDTO
	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
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
		&row.IsActive,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.ProductDTO{}, ErrProductNotFound
		}
		return &model.ProductDTO{}, fmt.Errorf("scan product: %w", err)
	}

	return &row, nil
}

func (r *productRepository) List(ctx context.Context, params *model.ProductParams) (model.ProductDTOs, error) {
	args := []interface{}{}
	conditions := []string{}
	query := `
		SELECT id, sku, name, description, category, unit_of_measure, 
			weight, length, width, height, barcode, is_active, created_at, updated_at
		FROM products
	`

	// Active filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	// Category filter
	if params.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", len(args)+1))
		args = append(args, params.Category)
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

	rows, err := scanProducts(ctx, r.db, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	return rows, nil
}

func (r *productRepository) Count(ctx context.Context, params *model.ProductParams) (int, error) {
	var count int
	args := []interface{}{}
	conditions := []string{}

	query := `SELECT COUNT(*) FROM products`

	// Active filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	// Category filter
	if params.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", len(args)+1))
		args = append(args, params.Category)
	}

	// Apply WHERE if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count products: %w", err)
	}

	return count, nil
}

func scanProducts(ctx context.Context, database *db.DB, query string, args ...any) (model.ProductDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scan products: %w", err)
	}
	defer rows.Close()

	var products model.ProductDTOs
	for rows.Next() {
		var row model.ProductDTO
		if err := rows.Scan(
			&row.ID,
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
			&row.IsActive,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan product row: %w", err)
		}
		products = append(products, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product rows: %w", err)
	}

	return products, nil
}
