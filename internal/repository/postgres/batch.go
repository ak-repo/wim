package postgres

import (
	"context"
	"fmt"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type BatchRepository interface {
	Create(ctx context.Context, batch *domain.Batch) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Batch, error)
	GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.Batch, error)
	GetExpiringSoon(ctx context.Context, days int) ([]*domain.Batch, error)
	Update(ctx context.Context, batch *domain.Batch) error
}

type batchRepo struct {
	db *DB
}

func NewBatchRepository(db *DB) BatchRepository {
	return &batchRepo{db: db}
}

func (r *batchRepo) Create(ctx context.Context, batch *domain.Batch) error {
	query := `
		INSERT INTO batches (id, batch_number, product_id, supplier_id, manufacturing_date, expiry_date, origin_country, quantity_initial, quantity_remaining, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Pool.Exec(ctx, query,
		batch.ID, batch.BatchNumber, batch.ProductID, batch.SupplierID,
		batch.ManufacturingDate, batch.ExpiryDate, batch.OriginCountry,
		batch.QuantityInitial, batch.QuantityRemaining, batch.IsActive,
		batch.CreatedAt, batch.UpdatedAt,
	)
	return err
}

func (r *batchRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Batch, error) {
	query := `
		SELECT id, batch_number, product_id, supplier_id, manufacturing_date, expiry_date, origin_country, quantity_initial, quantity_remaining, is_active, created_at, updated_at
		FROM batches WHERE id = $1`

	var b domain.Batch
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.BatchNumber, &b.ProductID, &b.SupplierID,
		&b.ManufacturingDate, &b.ExpiryDate, &b.OriginCountry,
		&b.QuantityInitial, &b.QuantityRemaining, &b.IsActive,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *batchRepo) GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.Batch, error) {
	query := `
		SELECT id, batch_number, product_id, supplier_id, manufacturing_date, expiry_date, origin_country, quantity_initial, quantity_remaining, is_active, created_at, updated_at
		FROM batches WHERE product_id = $1 AND is_active = true
		ORDER BY expiry_date ASC`

	rows, err := r.db.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batches []*domain.Batch
	for rows.Next() {
		var b domain.Batch
		err := rows.Scan(
			&b.ID, &b.BatchNumber, &b.ProductID, &b.SupplierID,
			&b.ManufacturingDate, &b.ExpiryDate, &b.OriginCountry,
			&b.QuantityInitial, &b.QuantityRemaining, &b.IsActive,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &b)
	}
	return batches, nil
}

func (r *batchRepo) GetExpiringSoon(ctx context.Context, days int) ([]*domain.Batch, error) {
	query := `
		SELECT id, batch_number, product_id, supplier_id, manufacturing_date, expiry_date, origin_country, quantity_initial, quantity_remaining, is_active, created_at, updated_at
		FROM batches 
		WHERE is_active = true AND expiry_date IS NOT NULL 
		AND expiry_date <= NOW() + ($1 || ' days')::interval
		ORDER BY expiry_date ASC`

	rows, err := r.db.Pool.Query(ctx, query, fmt.Sprintf("%d", days))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batches []*domain.Batch
	for rows.Next() {
		var b domain.Batch
		err := rows.Scan(
			&b.ID, &b.BatchNumber, &b.ProductID, &b.SupplierID,
			&b.ManufacturingDate, &b.ExpiryDate, &b.OriginCountry,
			&b.QuantityInitial, &b.QuantityRemaining, &b.IsActive,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &b)
	}
	return batches, nil
}

func (r *batchRepo) Update(ctx context.Context, batch *domain.Batch) error {
	query := `
		UPDATE batches SET 
			batch_number = $2, product_id = $3, supplier_id = $4, manufacturing_date = $5,
			expiry_date = $6, origin_country = $7, quantity_initial = $8, 
			quantity_remaining = $9, is_active = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		batch.ID, batch.BatchNumber, batch.ProductID, batch.SupplierID,
		batch.ManufacturingDate, batch.ExpiryDate, batch.OriginCountry,
		batch.QuantityInitial, batch.QuantityRemaining, batch.IsActive,
		batch.UpdatedAt,
	)
	return err
}
