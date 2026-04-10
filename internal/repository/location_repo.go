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

var ErrLocationNotFound = errors.New("location not found")

type LocationRepository interface {
	Create(ctx context.Context, location *model.LocationRequest) (int, error)
	GetByID(ctx context.Context, locationID int) (*model.LocationDTO, error)
	GetByCode(ctx context.Context, code string) (*model.LocationDTO, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	Update(ctx context.Context, locationID int, location *model.LocationRequest) error
	Delete(ctx context.Context, locationID int) error
	List(ctx context.Context, params *model.LocationParams) (model.LocationDTOs, error)
	Count(ctx context.Context, params *model.LocationParams) (int, error)
	ListByWarehouse(ctx context.Context, warehouseID int) (model.LocationDTOs, error)
}

type locationRepository struct {
	db *db.DB
}

func NewLocationRepository(database *db.DB) LocationRepository {
	return &locationRepository{
		db: database,
	}
}

func (r *locationRepository) Create(ctx context.Context, location *model.LocationRequest) (int, error) {
	var warehouseCode string
	err := r.db.Pool.QueryRow(ctx, `
		SELECT code FROM warehouses WHERE id = $1
	`, location.WarehouseID).Scan(&warehouseCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrWarehouseNotFound
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to lookup warehouse")
	}

	query := `
		INSERT INTO locations (
			 ref_code, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING id
	`

	var id int
	err = r.db.Pool.QueryRow(ctx, query,
		location.RefCode,
		location.WarehouseID,
		location.Zone,
		location.Aisle,
		location.Rack,
		location.Bin,
		location.LocationCode,
		location.LocationType,
		location.IsPickFace,
		location.MaxWeight,
		location.IsActive,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, apperrors.ErrAlreadyExists
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create location")
	}

	return id, nil
}

func (r *locationRepository) GetByID(ctx context.Context, locationID int) (*model.LocationDTO, error) {
	return scanLocation(ctx, r.db, `
		SELECT id, ref_code, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at, deleted_at
		FROM locations WHERE id = $1 AND deleted_at IS NULL
	`, locationID)
}

func (r *locationRepository) GetByCode(ctx context.Context, code string) (*model.LocationDTO, error) {
	return scanLocation(ctx, r.db, `
		SELECT id, ref_code, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at, deleted_at
		FROM locations WHERE location_code = $1 AND deleted_at IS NULL
	`, code)
}

func (r *locationRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM locations WHERE location_code = $1 AND deleted_at IS NULL)`, code).Scan(&exists)
	if err != nil {
		return false, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check location by code")
	}

	return exists, nil
}

func (r *locationRepository) Update(ctx context.Context, locationID int, location *model.LocationRequest) error {
	query := `
		UPDATE locations
		SET zone = COALESCE($2, zone),
			aisle = COALESCE($3, aisle),
			rack = COALESCE($4, rack),
			bin = COALESCE($5, bin),
			location_code = COALESCE($6, location_code),
			location_type = COALESCE($7, location_type),
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
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update location")
	}

	if result.RowsAffected() == 0 {
		return ErrLocationNotFound
	}

	return nil
}

func (r *locationRepository) Delete(ctx context.Context, locationID int) error {
	result, err := r.db.Pool.Exec(ctx, `UPDATE locations SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, locationID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete location")
	}

	if result.RowsAffected() == 0 {
		return ErrLocationNotFound
	}

	return nil
}

func scanLocation(ctx context.Context, database *db.DB, query string, args ...any) (*model.LocationDTO, error) {
	var row model.LocationDTO
	var isActive, isPickFace sql.NullBool
	var createdAt, updatedAt, deletedAt sql.NullTime

	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.RefCode,
		&row.WarehouseID,
		&row.Zone,
		&row.Aisle,
		&row.Rack,
		&row.Bin,
		&row.LocationCode,
		&row.LocationType,
		&isPickFace,
		&row.MaxWeight,
		&isActive,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err == nil {
		row.ApplyNullScalars(isActive, isPickFace, createdAt, updatedAt, deletedAt)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load location")
	}

	return &row, nil
}

func (r *locationRepository) List(ctx context.Context, params *model.LocationParams) (model.LocationDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, ref_code, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at, deleted_at
		FROM locations
	`

	// Base condition: only get non-deleted records
	conditions = append(conditions, "deleted_at IS NULL")

	// Active filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	// Warehouse filter
	if params.WarehouseID != 0 {
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

	return scanLocations(ctx, r.db, query, args...)
}

func (r *locationRepository) Count(ctx context.Context, params *model.LocationParams) (int, error) {
	var count int
	var args []any
	var conditions []string

	query := `SELECT COUNT(*) FROM locations`

	// Base condition: only count non-deleted records
	conditions = append(conditions, "deleted_at IS NULL")

	// Active filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	// Warehouse filter
	if params.WarehouseID != 0 {
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
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count locations")
	}

	return count, nil
}

func (r *locationRepository) ListByWarehouse(ctx context.Context, warehouseID int) (model.LocationDTOs, error) {
	query := `
		SELECT id, ref_code, warehouse_id, zone, aisle, rack, bin, location_code, location_type, is_pick_face, max_weight, is_active, created_at, updated_at, deleted_at
		FROM locations WHERE warehouse_id = $1 AND deleted_at IS NULL
		ORDER BY zone, aisle, rack, bin
	`
	return scanLocations(ctx, r.db, query, warehouseID)
}

func scanLocations(ctx context.Context, database *db.DB, query string, args ...any) (model.LocationDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query locations")
	}
	defer rows.Close()

	var locations model.LocationDTOs
	for rows.Next() {
		var row model.LocationDTO
		var isActive, isPickFace sql.NullBool
		var createdAt, updatedAt, deletedAt sql.NullTime

		err := rows.Scan(
			&row.ID,
			&row.RefCode,
			&row.WarehouseID,
			&row.Zone,
			&row.Aisle,
			&row.Rack,
			&row.Bin,
			&row.LocationCode,
			&row.LocationType,
			&isPickFace,
			&row.MaxWeight,
			&isActive,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan location row")
		}
		row.ApplyNullScalars(isActive, isPickFace, createdAt, updatedAt, deletedAt)
		locations = append(locations, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate locations")
	}

	return locations, nil
}
