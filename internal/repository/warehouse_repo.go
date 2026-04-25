package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/model"
	apperrors "github.com/ak-repo/wim/internal/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrWarehouseNotFound = errors.New("warehouse not found")

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *model.WarehouseRequest) (int, error)
	GetByID(ctx context.Context, warehouseID int) (*model.WarehouseDTO, error)
	GetByCode(ctx context.Context, code string) (*model.WarehouseDTO, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	Update(ctx context.Context, warehouseID int, warehouse *model.WarehouseRequest) error
	Delete(ctx context.Context, warehouseID int) error
	List(ctx context.Context, params *model.WarehouseParams) (model.WarehouseDTOs, error)
	Count(ctx context.Context, params *model.WarehouseParams) (int, error)
}

type warehouseRepository struct {
	db *db.DB
}

func NewWarehouseRepository(database *db.DB) WarehouseRepository {
	return &warehouseRepository{
		db: database,
	}
}

func (r *warehouseRepository) Create(ctx context.Context, warehouse *model.WarehouseRequest) (int, error) {
	query := `
		INSERT INTO warehouses (
			 ref_code, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query,
		warehouse.RefCode,
		warehouse.Code,
		warehouse.Name,
		warehouse.AddressLine1,
		warehouse.AddressLine2,
		warehouse.City,
		warehouse.State,
		warehouse.PostalCode,
		warehouse.Country,
		warehouse.IsActive,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, apperrors.ErrAlreadyExists
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create warehouse")
	}

	return id, nil
}

func (r *warehouseRepository) GetByID(ctx context.Context, warehouseID int) (*model.WarehouseDTO, error) {
	return scanWarehouse(ctx, r.db, `
		SELECT id, ref_code, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at, deleted_at
		FROM warehouses WHERE id = $1 AND deleted_at IS NULL
	`, warehouseID)
}

func (r *warehouseRepository) GetByCode(ctx context.Context, code string) (*model.WarehouseDTO, error) {
	return scanWarehouse(ctx, r.db, `
		SELECT id, ref_code, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at, deleted_at
		FROM warehouses WHERE code = $1 AND deleted_at IS NULL
	`, code)
}

func (r *warehouseRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM warehouses WHERE code = $1 AND deleted_at IS NULL)`, code).Scan(&exists)
	if err != nil {
		return false, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check warehouse by code")
	}

	return exists, nil
}

func (r *warehouseRepository) Update(ctx context.Context, warehouseID int, warehouse *model.WarehouseRequest) error {
	query := `
		UPDATE warehouses
		SET code = COALESCE($2, code),
			name = COALESCE($3, name),
			address_line1 = COALESCE($4, address_line1),
			address_line2 = COALESCE($5, address_line2),
			city = COALESCE($6, city),
			state = COALESCE($7, state),
			postal_code = COALESCE($8, postal_code),
			country = COALESCE($9, country),
			is_active = COALESCE($10, is_active),
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		warehouseID,
		warehouse.Code,
		warehouse.Name,
		warehouse.AddressLine1,
		warehouse.AddressLine2,
		warehouse.City,
		warehouse.State,
		warehouse.PostalCode,
		warehouse.Country,
		warehouse.IsActive,
	)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update warehouse")
	}

	if result.RowsAffected() == 0 {
		return ErrWarehouseNotFound
	}

	return nil
}

func (r *warehouseRepository) Delete(ctx context.Context, warehouseID int) error {
	result, err := r.db.Pool.Exec(ctx, `UPDATE warehouses SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, warehouseID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete warehouse")
	}

	if result.RowsAffected() == 0 {
		return ErrWarehouseNotFound
	}

	return nil
}

func scanWarehouse(ctx context.Context, database *db.DB, query string, args ...any) (*model.WarehouseDTO, error) {
	var row model.WarehouseDTO
	var isActive sql.NullBool
	var createdAt, updatedAt, deletedAt sql.NullTime

	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.RefCode,
		&row.Code,
		&row.Name,
		&row.AddressLine1,
		&row.AddressLine2,
		&row.City,
		&row.State,
		&row.PostalCode,
		&row.Country,
		&isActive,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err == nil {
		row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWarehouseNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load warehouse")
	}

	return &row, nil
}

func (r *warehouseRepository) List(ctx context.Context, params *model.WarehouseParams) (model.WarehouseDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, ref_code, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at, deleted_at
		FROM warehouses
	`

	// Base condition: only get non-deleted records
	conditions = append(conditions, "deleted_at IS NULL")

	// Active filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
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

	return scanWarehouses(ctx, r.db, query, args...)
}

func (r *warehouseRepository) Count(ctx context.Context, params *model.WarehouseParams) (int, error) {
	var count int
	var args []any
	var conditions []string

	query := `SELECT COUNT(*) FROM warehouses`

	// Base condition: only count non-deleted records
	conditions = append(conditions, "deleted_at IS NULL")

	// Active filter
	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	// Apply WHERE if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count warehouses")
	}

	return count, nil
}

func scanWarehouses(ctx context.Context, database *db.DB, query string, args ...any) (model.WarehouseDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query warehouses")
	}
	defer rows.Close()

	var warehouses model.WarehouseDTOs
	for rows.Next() {
		var row model.WarehouseDTO
		var isActive sql.NullBool
		var createdAt, updatedAt, deletedAt sql.NullTime

		err := rows.Scan(
			&row.ID,
			&row.RefCode,
			&row.Code,
			&row.Name,
			&row.AddressLine1,
			&row.AddressLine2,
			&row.City,
			&row.State,
			&row.PostalCode,
			&row.Country,
			&isActive,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan warehouse row")
		}
		row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)
		warehouses = append(warehouses, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate warehouses")
	}

	return warehouses, nil
}
