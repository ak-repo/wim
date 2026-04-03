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

var ErrWarehouseNotFound = errors.New("warehouse not found")

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *model.CreateWarehouseRequest) (uuid.UUID, error)
	GetByID(ctx context.Context, warehouseID uuid.UUID) (*model.WarehouseDTO, error)
	GetByCode(ctx context.Context, code string) (*model.WarehouseDTO, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	Update(ctx context.Context, warehouseID uuid.UUID, warehouse *model.UpdateWarehouseRequest) error
	Delete(ctx context.Context, warehouseID uuid.UUID) error
	List(ctx context.Context, params *model.WarehouseParams) (model.WarehouseDTOs, error)
	Count(ctx context.Context, params *model.WarehouseParams) (int, error)
}

type warehouseRepository struct {
	db *db.DB
}

func NewWarehouseRepository(database *db.DB) WarehouseRepository {
	return &warehouseRepository{db: database}
}

func (r *warehouseRepository) Create(ctx context.Context, warehouse *model.CreateWarehouseRequest) (uuid.UUID, error) {
	query := `
		INSERT INTO warehouses (
			id, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true, NOW(), NOW())
		RETURNING id
	`

	id := uuid.New()
	_, err := r.db.Pool.Exec(ctx, query,
		id, warehouse.Code, warehouse.Name, warehouse.AddressLine1, warehouse.AddressLine2,
		warehouse.City, warehouse.State, warehouse.PostalCode, warehouse.Country,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, apperrors.ErrAlreadyExists
		}
		return uuid.Nil, fmt.Errorf("failed to create warehouse: %w", err)
	}

	return id, nil
}

func (r *warehouseRepository) GetByID(ctx context.Context, warehouseID uuid.UUID) (*model.WarehouseDTO, error) {
	row, err := scanWarehouse(ctx, r.db, `
		SELECT id, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at
		FROM warehouses WHERE id = $1
	`, warehouseID)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *warehouseRepository) GetByCode(ctx context.Context, code string) (*model.WarehouseDTO, error) {
	row, err := scanWarehouse(ctx, r.db, `
		SELECT id, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at
		FROM warehouses WHERE code = $1
	`, code)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *warehouseRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM warehouses WHERE code = $1)`, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check warehouse by code: %w", err)
	}

	return exists, nil
}

func (r *warehouseRepository) Update(ctx context.Context, warehouseID uuid.UUID, warehouse *model.UpdateWarehouseRequest) error {
	query := `
		UPDATE warehouses
		SET name = COALESCE(NULLIF($2, ''), name),
			address_line1 = COALESCE(NULLIF($3, ''), address_line1),
			address_line2 = COALESCE(NULLIF($4, ''), address_line2),
			city = COALESCE(NULLIF($5, ''), city),
			state = COALESCE(NULLIF($6, ''), state),
			postal_code = COALESCE(NULLIF($7, ''), postal_code),
			country = COALESCE(NULLIF($8, ''), country),
			is_active = COALESCE($9, is_active),
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		warehouseID,
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
		return fmt.Errorf("failed to update warehouse: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrWarehouseNotFound
	}

	return nil
}

func (r *warehouseRepository) Delete(ctx context.Context, warehouseID uuid.UUID) error {
	result, err := r.db.Pool.Exec(ctx, `DELETE FROM warehouses WHERE id = $1`, warehouseID)
	if err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrWarehouseNotFound
	}

	return nil
}

func scanWarehouse(ctx context.Context, database *db.DB, query string, args ...any) (*model.WarehouseDTO, error) {
	var row model.WarehouseDTO
	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.Code,
		&row.Name,
		&row.AddressLine1,
		&row.AddressLine2,
		&row.City,
		&row.State,
		&row.PostalCode,
		&row.Country,
		&row.IsActive,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.WarehouseDTO{}, ErrWarehouseNotFound
		}
		return &model.WarehouseDTO{}, fmt.Errorf("scan warehouse: %w", err)
	}

	return &row, nil
}

func (r *warehouseRepository) List(ctx context.Context, params *model.WarehouseParams) (model.WarehouseDTOs, error) {
	args := []interface{}{}
	conditions := []string{}
	query := `
		SELECT id, code, name, address_line1, address_line2, city, state, postal_code, country, is_active, created_at, updated_at
		FROM warehouses
	`

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

	rows, err := scanWarehouses(ctx, r.db, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list warehouses: %w", err)
	}

	return rows, nil
}

func (r *warehouseRepository) Count(ctx context.Context, params *model.WarehouseParams) (int, error) {
	var count int
	args := []interface{}{}
	conditions := []string{}

	query := `SELECT COUNT(*) FROM warehouses`

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
		return 0, fmt.Errorf("count warehouses: %w", err)
	}

	return count, nil
}

func scanWarehouses(ctx context.Context, database *db.DB, query string, args ...any) (model.WarehouseDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scan warehouses: %w", err)
	}
	defer rows.Close()

	var warehouses model.WarehouseDTOs
	for rows.Next() {
		var row model.WarehouseDTO
		if err := rows.Scan(
			&row.ID,
			&row.Code,
			&row.Name,
			&row.AddressLine1,
			&row.AddressLine2,
			&row.City,
			&row.State,
			&row.PostalCode,
			&row.Country,
			&row.IsActive,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan warehouse row: %w", err)
		}
		warehouses = append(warehouses, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate warehouse rows: %w", err)
	}

	return warehouses, nil
}
