package postgres

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type StockMovementRepository interface {
	Create(ctx context.Context, movement *domain.StockMovement) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.StockMovement, error)
	GetByReference(ctx context.Context, refType string, refID uuid.UUID) ([]*domain.StockMovement, error)
	GetByProduct(ctx context.Context, productID uuid.UUID, limit int) ([]*domain.StockMovement, error)
	GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit int) ([]*domain.StockMovement, error)
}

type stockMovementRepo struct {
	db *DB
}

func NewStockMovementRepository(db *DB) StockMovementRepository {
	return &stockMovementRepo{db: db}
}

func (r *stockMovementRepo) Create(ctx context.Context, movement *domain.StockMovement) error {
	query := `
		INSERT INTO stock_movements (id, movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id, quantity, reference_type, reference_id, performed_by, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.Pool.Exec(ctx, query,
		movement.ID, movement.MovementType, movement.ProductID, movement.WarehouseID,
		movement.LocationIDFrom, movement.LocationIDTo, movement.BatchID, movement.Quantity,
		movement.ReferenceType, movement.ReferenceID, movement.PerformedBy, movement.Notes, movement.CreatedAt,
	)
	return err
}

func (r *stockMovementRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.StockMovement, error) {
	query := `
		SELECT id, movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id, quantity, reference_type, reference_id, performed_by, notes, created_at
		FROM stock_movements WHERE id = $1`

	var m domain.StockMovement
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.MovementType, &m.ProductID, &m.WarehouseID,
		&m.LocationIDFrom, &m.LocationIDTo, &m.BatchID, &m.Quantity,
		&m.ReferenceType, &m.ReferenceID, &m.PerformedBy, &m.Notes, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *stockMovementRepo) GetByReference(ctx context.Context, refType string, refID uuid.UUID) ([]*domain.StockMovement, error) {
	query := `
		SELECT id, movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id, quantity, reference_type, reference_id, performed_by, notes, created_at
		FROM stock_movements WHERE reference_type = $1 AND reference_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, refType, refID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []*domain.StockMovement
	for rows.Next() {
		var m domain.StockMovement
		err := rows.Scan(
			&m.ID, &m.MovementType, &m.ProductID, &m.WarehouseID,
			&m.LocationIDFrom, &m.LocationIDTo, &m.BatchID, &m.Quantity,
			&m.ReferenceType, &m.ReferenceID, &m.PerformedBy, &m.Notes, &m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		movements = append(movements, &m)
	}
	return movements, nil
}

func (r *stockMovementRepo) GetByProduct(ctx context.Context, productID uuid.UUID, limit int) ([]*domain.StockMovement, error) {
	query := `
		SELECT id, movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id, quantity, reference_type, reference_id, performed_by, notes, created_at
		FROM stock_movements WHERE product_id = $1
		ORDER BY created_at DESC LIMIT $2`

	rows, err := r.db.Pool.Query(ctx, query, productID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []*domain.StockMovement
	for rows.Next() {
		var m domain.StockMovement
		err := rows.Scan(
			&m.ID, &m.MovementType, &m.ProductID, &m.WarehouseID,
			&m.LocationIDFrom, &m.LocationIDTo, &m.BatchID, &m.Quantity,
			&m.ReferenceType, &m.ReferenceID, &m.PerformedBy, &m.Notes, &m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		movements = append(movements, &m)
	}
	return movements, nil
}

func (r *stockMovementRepo) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit int) ([]*domain.StockMovement, error) {
	query := `
		SELECT id, movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id, quantity, reference_type, reference_id, performed_by, notes, created_at
		FROM stock_movements WHERE warehouse_id = $1
		ORDER BY created_at DESC LIMIT $2`

	rows, err := r.db.Pool.Query(ctx, query, warehouseID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []*domain.StockMovement
	for rows.Next() {
		var m domain.StockMovement
		err := rows.Scan(
			&m.ID, &m.MovementType, &m.ProductID, &m.WarehouseID,
			&m.LocationIDFrom, &m.LocationIDTo, &m.BatchID, &m.Quantity,
			&m.ReferenceType, &m.ReferenceID, &m.PerformedBy, &m.Notes, &m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		movements = append(movements, &m)
	}
	return movements, nil
}
