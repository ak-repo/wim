package postgres

import (
	"context"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	GetByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*domain.Location, error)
	GetByCode(ctx context.Context, code string, warehouseID uuid.UUID) (*domain.Location, error)
	Update(ctx context.Context, location *domain.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type locationRepo struct {
	db *DB
}

func NewLocationRepository(db *DB) LocationRepository {
	return &locationRepo{db: db}
}

func (r *locationRepo) Create(ctx context.Context, location *domain.Location) error {
	query := `
		INSERT INTO locations (id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.Pool.Exec(ctx, query,
		location.ID, location.WarehouseID, location.Zone, location.Aisle, location.Rack, location.Bin,
		location.LocationCode, location.LocationType, location.IsPickFace, location.MaxWeight,
		location.IsActive, location.CreatedAt, location.UpdatedAt,
	)
	return err
}

func (r *locationRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	query := `
		SELECT id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at
		FROM locations WHERE id = $1`

	var loc domain.Location
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&loc.ID, &loc.WarehouseID, &loc.Zone, &loc.Aisle, &loc.Rack, &loc.Bin,
		&loc.LocationCode, &loc.LocationType, &loc.IsPickFace, &loc.MaxWeight,
		&loc.IsActive, &loc.CreatedAt, &loc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *locationRepo) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*domain.Location, error) {
	query := `
		SELECT id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at
		FROM locations WHERE warehouse_id = $1
		ORDER BY zone, aisle, rack, bin`

	rows, err := r.db.Pool.Query(ctx, query, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locs []*domain.Location
	for rows.Next() {
		var loc domain.Location
		err := rows.Scan(
			&loc.ID, &loc.WarehouseID, &loc.Zone, &loc.Aisle, &loc.Rack, &loc.Bin,
			&loc.LocationCode, &loc.LocationType, &loc.IsPickFace, &loc.MaxWeight,
			&loc.IsActive, &loc.CreatedAt, &loc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		locs = append(locs, &loc)
	}
	return locs, nil
}

func (r *locationRepo) GetByCode(ctx context.Context, code string, warehouseID uuid.UUID) (*domain.Location, error) {
	query := `
		SELECT id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at
		FROM locations WHERE location_code = $1 AND warehouse_id = $2`

	var loc domain.Location
	err := r.db.Pool.QueryRow(ctx, query, code, warehouseID).Scan(
		&loc.ID, &loc.WarehouseID, &loc.Zone, &loc.Aisle, &loc.Rack, &loc.Bin,
		&loc.LocationCode, &loc.LocationType, &loc.IsPickFace, &loc.MaxWeight,
		&loc.IsActive, &loc.CreatedAt, &loc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *locationRepo) Update(ctx context.Context, location *domain.Location) error {
	query := `
		UPDATE locations SET 
			zone = $2, aisle = $3, rack = $4, bin = $5, location_code = $6,
			location_type = $7, is_pick_face = $8, max_weight = $9, is_active = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		location.ID, location.Zone, location.Aisle, location.Rack, location.Bin,
		location.LocationCode, location.LocationType, location.IsPickFace, location.MaxWeight,
		location.IsActive, location.UpdatedAt,
	)
	return err
}

func (r *locationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM locations WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
