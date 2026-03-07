package postgres

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type BarcodeRepository interface {
	Create(ctx context.Context, barcode *domain.Barcode) error
	GetByValue(ctx context.Context, value string) (*domain.Barcode, error)
	GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.Barcode, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type barcodeRepo struct {
	db *DB
}

func NewBarcodeRepository(db *DB) BarcodeRepository {
	return &barcodeRepo{db: db}
}

func (r *barcodeRepo) Create(ctx context.Context, barcode *domain.Barcode) error {
	query := `
		INSERT INTO barcodes (id, product_id, barcode_value, barcode_type, is_primary, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Pool.Exec(ctx, query,
		barcode.ID, barcode.ProductID, barcode.BarcodeValue, barcode.BarcodeType,
		barcode.IsPrimary, barcode.CreatedAt,
	)
	return err
}

func (r *barcodeRepo) GetByValue(ctx context.Context, value string) (*domain.Barcode, error) {
	query := `
		SELECT id, product_id, barcode_value, barcode_type, is_primary, created_at
		FROM barcodes WHERE barcode_value = $1`

	var b domain.Barcode
	err := r.db.Pool.QueryRow(ctx, query, value).Scan(
		&b.ID, &b.ProductID, &b.BarcodeValue, &b.BarcodeType,
		&b.IsPrimary, &b.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *barcodeRepo) GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.Barcode, error) {
	query := `
		SELECT id, product_id, barcode_value, barcode_type, is_primary, created_at
		FROM barcodes WHERE product_id = $1`

	rows, err := r.db.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barcodes []*domain.Barcode
	for rows.Next() {
		var b domain.Barcode
		err := rows.Scan(
			&b.ID, &b.ProductID, &b.BarcodeValue, &b.BarcodeType,
			&b.IsPrimary, &b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		barcodes = append(barcodes, &b)
	}
	return barcodes, nil
}

func (r *barcodeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM barcodes WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
