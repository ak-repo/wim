package postgres

import (
	"context"
	"fmt"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Product, error)
	List(ctx context.Context, filter ProductFilter) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductFilter struct {
	Category string
	IsActive *bool
	Search   string
	Limit    int
	Offset   int
}

type productRepo struct {
	db *DB
}

func NewProductRepository(db *DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (id, sku, name, description, category, unit_of_measure, 
			weight, length, width, height, barcode, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.db.Pool.Exec(ctx, query,
		product.ID, product.SKU, product.Name, product.Description, product.Category,
		product.UnitOfMeasure, product.Weight, product.Length, product.Width, product.Height,
		product.Barcode, product.IsActive, product.CreatedAt, product.UpdatedAt,
	)
	return err
}

func (r *productRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, sku, name, description, category, unit_of_measure, 
			weight, length, width, height, barcode, is_active, created_at, updated_at
		FROM products WHERE id = $1`

	var p domain.Product
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.UnitOfMeasure,
		&p.Weight, &p.Length, &p.Width, &p.Height, &p.Barcode, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	query := `
		SELECT id, sku, name, description, category, unit_of_measure, 
			weight, length, width, height, barcode, is_active, created_at, updated_at
		FROM products WHERE sku = $1`

	var p domain.Product
	err := r.db.Pool.QueryRow(ctx, query, sku).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.UnitOfMeasure,
		&p.Weight, &p.Length, &p.Width, &p.Height, &p.Barcode, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) List(ctx context.Context, filter ProductFilter) ([]*domain.Product, error) {
	query := `
		SELECT id, sku, name, description, category, unit_of_measure, 
			weight, length, width, height, barcode, is_active, created_at, updated_at
		FROM products WHERE 1=1`

	args := []interface{}{}
	argNum := 1

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argNum)
		args = append(args, filter.Category)
		argNum++
	}
	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argNum)
		args = append(args, *filter.IsActive)
		argNum++
	}
	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+filter.Search+"%")
		argNum++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argNum)
		args = append(args, filter.Limit)
		argNum++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argNum)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Description, &p.Category, &p.UnitOfMeasure,
			&p.Weight, &p.Length, &p.Width, &p.Height, &p.Barcode, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

func (r *productRepo) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products SET 
			sku = $2, name = $3, description = $4, category = $5, unit_of_measure = $6,
			weight = $7, length = $8, width = $9, height = $10, barcode = $11, 
			is_active = $12, updated_at = $13
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		product.ID, product.SKU, product.Name, product.Description, product.Category,
		product.UnitOfMeasure, product.Weight, product.Length, product.Width, product.Height,
		product.Barcode, product.IsActive, product.UpdatedAt,
	)
	return err
}

func (r *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
