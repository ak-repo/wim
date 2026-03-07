package postgres

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *domain.Warehouse) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error)
	GetByCode(ctx context.Context, code string) (*domain.Warehouse, error)
	List(ctx context.Context, filter WarehouseFilter) ([]*domain.Warehouse, error)
	Update(ctx context.Context, warehouse *domain.Warehouse) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type WarehouseFilter struct {
	IsActive *bool
	Search   string
	Limit    int
	Offset   int
}

type warehouseRepo struct {
	db *DB
}

func NewWarehouseRepository(db *DB) WarehouseRepository {
	return &warehouseRepo{db: db}
}

func (r *warehouseRepo) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	query := `
		INSERT INTO warehouses (id, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Pool.Exec(ctx, query,
		warehouse.ID, warehouse.Code, warehouse.Name, warehouse.AddressLine1, warehouse.AddressLine2,
		warehouse.City, warehouse.State, warehouse.PostalCode, warehouse.Country, warehouse.IsActive,
		warehouse.CreatedAt, warehouse.UpdatedAt,
	)
	return err
}

func (r *warehouseRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error) {
	query := `
		SELECT id, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at
		FROM warehouses WHERE id = $1`

	var w domain.Warehouse
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&w.ID, &w.Code, &w.Name, &w.AddressLine1, &w.AddressLine2,
		&w.City, &w.State, &w.PostalCode, &w.Country, &w.IsActive,
		&w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *warehouseRepo) GetByCode(ctx context.Context, code string) (*domain.Warehouse, error) {
	query := `
		SELECT id, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at
		FROM warehouses WHERE code = $1`

	var w domain.Warehouse
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(
		&w.ID, &w.Code, &w.Name, &w.AddressLine1, &w.AddressLine2,
		&w.City, &w.State, &w.PostalCode, &w.Country, &w.IsActive,
		&w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *warehouseRepo) List(ctx context.Context, filter WarehouseFilter) ([]*domain.Warehouse, error) {
	query := `
		SELECT id, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at
		FROM warehouses WHERE 1=1`

	args := []interface{}{}

	if filter.IsActive != nil {
		query += " AND is_active = $1"
		args = append(args, *filter.IsActive)
	}
	if filter.Search != "" {
		query += " AND (name ILIKE $2 OR code ILIKE $2)"
		args = append(args, "%"+filter.Search+"%")
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += " LIMIT $3"
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		query += " OFFSET $4"
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []*domain.Warehouse
	for rows.Next() {
		var w domain.Warehouse
		err := rows.Scan(
			&w.ID, &w.Code, &w.Name, &w.AddressLine1, &w.AddressLine2,
			&w.City, &w.State, &w.PostalCode, &w.Country, &w.IsActive,
			&w.CreatedAt, &w.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &w)
	}

	return warehouses, nil
}

func (r *warehouseRepo) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	query := `
		UPDATE warehouses SET 
			code = $2, name = $3, address_line1 = $4, address_line2 = $5, city = $6,
			state = $7, postal_code = $8, country = $9, is_active = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		warehouse.ID, warehouse.Code, warehouse.Name, warehouse.AddressLine1, warehouse.AddressLine2,
		warehouse.City, warehouse.State, warehouse.PostalCode, warehouse.Country, warehouse.IsActive,
		warehouse.UpdatedAt,
	)
	return err
}

func (r *warehouseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM warehouses WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
