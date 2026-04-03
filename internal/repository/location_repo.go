package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/model"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrLocationNotFound = errors.New("location not found")

type LocationRepository interface {
	Create(ctx context.Context, location *model.CreateLocationRequest) (uuid.UUID, error)
	GetByID(ctx context.Context, locationID uuid.UUID) (*model.LocationDTO, error)
	GetByCode(ctx context.Context, code string) (*model.LocationDTO, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	Update(ctx context.Context, locationID uuid.UUID, location *model.UpdateLocationRequest) error
	Delete(ctx context.Context, locationID uuid.UUID) error
	List(ctx context.Context, params *model.LocationParams) (model.LocationDTOs, error)
	Count(ctx context.Context, params *model.LocationParams) (int, error)
	ListByWarehouse(ctx context.Context, warehouseID uuid.UUID) (model.LocationDTOs, error)
}

type locationRepository struct {
	db *db.DB
}

func NewLocationRepository(database *db.DB) LocationRepository {
	return &locationRepository{db: database}
}

func (r *locationRepository) Create(ctx context.Context, location *model.CreateLocationRequest) (uuid.UUID, error) {
	query := `
		INSERT INTO locations (
			id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, true, NOW(), NOW())
		RETURNING id
	`

	id := uuid.New()
	_, err := r.db.Pool.Exec(ctx, query,
		id, location.WarehouseID, location.Zone, location.Aisle, location.Rack,
		location.Bin, location.LocationCode, location.LocationType, location.IsPickFace, location.MaxWeight,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, apperrors.ErrAlreadyExists
		}
		return uuid.Nil, fmt.Errorf("failed to create location: %w", err)
	}

	return id, nil
}

func (r *locationRepository) GetByID(ctx context.Context, locationID uuid.UUID) (*model.LocationDTO, error) {
	row, err := scanLocation(ctx, r.db, `
		SELECT id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at
		FROM locations WHERE id = $1
	`, locationID)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *locationRepository) GetByCode(ctx context.Context, code string) (*model.LocationDTO, error) {
	row, err := scanLocation(ctx, r.db, `
		SELECT id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at
		FROM locations WHERE location_code = $1
	`, code)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *locationRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM locations WHERE location_code = $1)`, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check location by code: %w", err)
	}

	return exists, nil
}

func (r *locationRepository) Update(ctx context.Context, locationID uuid.UUID, location *model.UpdateLocationRequest) error {
	query := `
		UPDATE locations
		SET zone = COALESCE(NULLIF($2, ''), zone),
			aisle = COALESCE(NULLIF($3, ''), aisle),
			rack = COALESCE(NULLIF($4, ''), rack),
			bin = COALESCE(NULLIF($5, ''), bin),
			location_code = COALESCE(NULLIF($6, ''), location_code),
			location_type = COALESCE(NULLIF($7, ''), location_type),
			is_pick_face = COALESCE($8, is_pick_face),
			max_weight = COALESCE($9, max_weight),
			is_active = COALESCE($10, is_active),
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		locationID,
		location.Zone,
		location.Aisle,
		location.Rack,
		location.Bin,
		location.LocationCode,
		location.LocationType,
		location.IsPickFace,
		location.MaxWeight,
		location.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrLocationNotFound
	}

	return nil
}

func (r *locationRepository) Delete(ctx context.Context, locationID uuid.UUID) error {
	result, err := r.db.Pool.Exec(ctx, `DELETE FROM locations WHERE id = $1`, locationID)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrLocationNotFound
	}

	return nil
}

func scanLocation(ctx context.Context, database *db.DB, query string, args ...any) (*model.LocationDTO, error) {
	var row model.LocationDTO
	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.WarehouseID,
		&row.Zone,
		&row.Aisle,
		&row.Rack,
		&row.Bin,
		&row.LocationCode,
		&row.LocationType,
		&row.IsPickFace,
		&row.MaxWeight,
		&row.IsActive,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.LocationDTO{}, ErrLocationNotFound
		}
		return &model.LocationDTO{}, fmt.Errorf("scan location: %w", err)
	}

	return &row, nil
}

func (r *locationRepository) List(ctx context.Context, params *model.LocationParams) (model.LocationDTOs, error) {
	args := []interface{}{}
	conditions := []string{}
	query := `
		SELECT id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at
		FROM locations
	`

	// Active filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	// Warehouse filter
	if params.WarehouseID != uuid.Nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, params.WarehouseID)
	}

	// Zone filter
	if params.Zone != "" {
		conditions = append(conditions, fmt.Sprintf("zone = $%d", len(args)+1))
		args = append(args, params.Zone)
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

	rows, err := scanLocations(ctx, r.db, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list locations: %w", err)
	}

	return rows, nil
}

func (r *locationRepository) Count(ctx context.Context, params *model.LocationParams) (int, error) {
	var count int
	args := []interface{}{}
	conditions := []string{}

	query := `SELECT COUNT(*) FROM locations`

	// Active filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	// Warehouse filter
	if params.WarehouseID != uuid.Nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, params.WarehouseID)
	}

	// Zone filter
	if params.Zone != "" {
		conditions = append(conditions, fmt.Sprintf("zone = $%d", len(args)+1))
		args = append(args, params.Zone)
	}

	// Apply WHERE if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count locations: %w", err)
	}

	return count, nil
}

func (r *locationRepository) ListByWarehouse(ctx context.Context, warehouseID uuid.UUID) (model.LocationDTOs, error) {
	query := `
		SELECT id, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at
		FROM locations WHERE warehouse_id = $1
		ORDER BY zone, aisle, rack, bin
	`
	rows, err := scanLocations(ctx, r.db, query, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("list locations by warehouse: %w", err)
	}
	return rows, nil
}

func scanLocations(ctx context.Context, database *db.DB, query string, args ...any) (model.LocationDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scan locations: %w", err)
	}
	defer rows.Close()

	var locations model.LocationDTOs
	for rows.Next() {
		var row model.LocationDTO
		if err := rows.Scan(
			&row.ID,
			&row.WarehouseID,
			&row.Zone,
			&row.Aisle,
			&row.Rack,
			&row.Bin,
			&row.LocationCode,
			&row.LocationType,
			&row.IsPickFace,
			&row.MaxWeight,
			&row.IsActive,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan location row: %w", err)
		}
		locations = append(locations, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate location rows: %w", err)
	}

	return locations, nil
}
