package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/model"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/jackc/pgx/v5"
)

var ErrInventoryNotFound = errors.New("inventory not found")

type InventoryRepository interface {
	GetByID(ctx context.Context, inventoryID int) (*model.InventoryDTO, error)
	GetByKey(ctx context.Context, productID, warehouseID, locationID int, batchID *int) (*model.InventoryDTO, error)
	List(ctx context.Context, params *model.InventoryParams) (model.InventoryDTOs, error)
	Count(ctx context.Context, params *model.InventoryParams) (int, error)
	Adjust(ctx context.Context, productID, warehouseID, locationID int, batchID *int, delta int, movementType, referenceType string, referenceID *int, performedBy *int, notes string) error
	ListMovements(ctx context.Context, params *model.StockMovementParams) (model.StockMovementDTOs, error)
	CountMovements(ctx context.Context, params *model.StockMovementParams) (int, error)
}

type inventoryRepository struct {
	db *db.DB
}

func NewInventoryRepository(database *db.DB) InventoryRepository {
	return &inventoryRepository{db: database}
}

func (r *inventoryRepository) GetByID(ctx context.Context, inventoryID int) (*model.InventoryDTO, error) {
	query := `
        SELECT id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_qty, version, created_at, updated_at
        FROM inventories
        WHERE id = $1
    `
	return scanInventory(ctx, r.db, query, inventoryID)
}

func (r *inventoryRepository) GetByKey(ctx context.Context, productID, warehouseID, locationID int, batchID *int) (*model.InventoryDTO, error) {
	query := `
        SELECT id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_qty, version, created_at, updated_at
        FROM inventories
        WHERE product_id = $1
          AND warehouse_id = $2
          AND location_id = $3
          AND ((batch_id IS NULL AND $4 IS NULL) OR batch_id = $4)
    `
	return scanInventory(ctx, r.db, query, productID, warehouseID, locationID, batchID)
}

func (r *inventoryRepository) List(ctx context.Context, params *model.InventoryParams) (model.InventoryDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
        SELECT id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_qty, version, created_at, updated_at
        FROM inventories
    `

	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", len(args)+1))
		args = append(args, *params.ProductID)
	}
	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("location_id = $%d", len(args)+1))
		args = append(args, *params.LocationID)
	}
	if params.BatchID != nil {
		conditions = append(conditions, fmt.Sprintf("batch_id = $%d", len(args)+1))
		args = append(args, *params.BatchID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, params.Limit, offset)

	return scanInventories(ctx, r.db, query, args...)
}

func (r *inventoryRepository) Count(ctx context.Context, params *model.InventoryParams) (int, error) {
	var (
		args       []any
		conditions []string
	)

	query := `SELECT COUNT(*) FROM inventories`

	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", len(args)+1))
		args = append(args, *params.ProductID)
	}
	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("location_id = $%d", len(args)+1))
		args = append(args, *params.LocationID)
	}
	if params.BatchID != nil {
		conditions = append(conditions, fmt.Sprintf("batch_id = $%d", len(args)+1))
		args = append(args, *params.BatchID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count inventory")
	}
	return count, nil
}

func (r *inventoryRepository) Adjust(ctx context.Context, productID, warehouseID, locationID int, batchID *int, delta int, movementType, referenceType string, referenceID *int, performedBy *int, notes string) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start inventory transaction")
	}
	defer tx.Rollback(ctx)

	selectQuery := `
        SELECT id, quantity, reserved_qty
        FROM inventories
        WHERE product_id = $1
          AND warehouse_id = $2
          AND location_id = $3
          AND ((batch_id IS NULL AND $4 IS NULL) OR batch_id = $4)
        FOR UPDATE
    `

	var (
		inventoryID int
		quantity    int
		reservedQty int
	)

	scanErr := tx.QueryRow(ctx, selectQuery, productID, warehouseID, locationID, batchID).
		Scan(&inventoryID, &quantity, &reservedQty)

	if scanErr != nil {
		if errors.Is(scanErr, pgx.ErrNoRows) {
			if delta < 0 {
				return apperrors.ErrInsufficientStock
			}

			insertQuery := `
                INSERT INTO inventories (product_id, warehouse_id, location_id, batch_id, quantity, reserved_qty, version, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, 0, 1, NOW(), NOW())
            `
			if _, err := tx.Exec(ctx, insertQuery, productID, warehouseID, locationID, batchID, delta); err != nil {
				return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create inventory")
			}
		} else {
			return apperrors.Wrap(scanErr, apperrors.CodeDatabase, "failed to load inventory")
		}
	} else {
		newQty := quantity + delta
		if newQty < 0 || newQty < reservedQty {
			return apperrors.ErrInsufficientStock
		}

		updateQuery := `
            UPDATE inventories
            SET quantity = $2,
                version = version + 1,
                updated_at = NOW()
            WHERE id = $1
        `
		if _, err := tx.Exec(ctx, updateQuery, inventoryID, newQty); err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update inventory")
		}
	}

	movementQuantity := int(math.Abs(float64(delta)))

	var locationIDFrom *int
	var locationIDTo *int
	if delta > 0 {
		locationIDTo = &locationID
	} else {
		locationIDFrom = &locationID
	}

	var notesValue any
	if strings.TrimSpace(notes) != "" {
		notesValue = notes
	}

	var referenceTypeValue any
	if strings.TrimSpace(referenceType) != "" {
		referenceTypeValue = referenceType
	}

	movementQuery := `
        INSERT INTO stock_movements (
            movement_type, product_id, warehouse_id, location_id_from, location_id_to,
            batch_id, quantity, reference_type, reference_id, performed_by, notes, created_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW())
    `
	if _, err := tx.Exec(ctx, movementQuery,
		movementType,
		productID,
		warehouseID,
		locationIDFrom,
		locationIDTo,
		batchID,
		movementQuantity,
		referenceTypeValue,
		referenceID,
		performedBy,
		notesValue,
	); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create stock movement")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit inventory transaction")
	}

	return nil
}

func (r *inventoryRepository) ListMovements(ctx context.Context, params *model.StockMovementParams) (model.StockMovementDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
        SELECT id, movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id,
               quantity, reference_type, reference_id, performed_by, notes, created_at
        FROM stock_movements
    `

	if params.MovementType != nil {
		conditions = append(conditions, fmt.Sprintf("movement_type = $%d", len(args)+1))
		args = append(args, *params.MovementType)
	}
	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", len(args)+1))
		args = append(args, *params.ProductID)
	}
	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("(location_id_from = $%d OR location_id_to = $%d)", len(args)+1, len(args)+1))
		args = append(args, *params.LocationID)
	}
	if params.BatchID != nil {
		conditions = append(conditions, fmt.Sprintf("batch_id = $%d", len(args)+1))
		args = append(args, *params.BatchID)
	}
	if params.ReferenceType != nil {
		conditions = append(conditions, fmt.Sprintf("reference_type = $%d", len(args)+1))
		args = append(args, *params.ReferenceType)
	}
	if params.ReferenceID != nil {
		conditions = append(conditions, fmt.Sprintf("reference_id = $%d", len(args)+1))
		args = append(args, *params.ReferenceID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, params.Limit, offset)

	return scanStockMovements(ctx, r.db, query, args...)
}

func (r *inventoryRepository) CountMovements(ctx context.Context, params *model.StockMovementParams) (int, error) {
	var (
		args       []any
		conditions []string
	)

	query := `SELECT COUNT(*) FROM stock_movements`

	if params.MovementType != nil {
		conditions = append(conditions, fmt.Sprintf("movement_type = $%d", len(args)+1))
		args = append(args, *params.MovementType)
	}
	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", len(args)+1))
		args = append(args, *params.ProductID)
	}
	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("(location_id_from = $%d OR location_id_to = $%d)", len(args)+1, len(args)+1))
		args = append(args, *params.LocationID)
	}
	if params.BatchID != nil {
		conditions = append(conditions, fmt.Sprintf("batch_id = $%d", len(args)+1))
		args = append(args, *params.BatchID)
	}
	if params.ReferenceType != nil {
		conditions = append(conditions, fmt.Sprintf("reference_type = $%d", len(args)+1))
		args = append(args, *params.ReferenceType)
	}
	if params.ReferenceID != nil {
		conditions = append(conditions, fmt.Sprintf("reference_id = $%d", len(args)+1))
		args = append(args, *params.ReferenceID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count stock movements")
	}
	return count, nil
}

func scanInventory(ctx context.Context, database *db.DB, query string, args ...any) (*model.InventoryDTO, error) {
	var row model.InventoryDTO
	var batchID sql.NullInt64
	var createdAt, updatedAt sql.NullTime

	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.ProductID,
		&row.WarehouseID,
		&row.LocationID,
		&batchID,
		&row.Quantity,
		&row.ReservedQty,
		&row.Version,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInventoryNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load inventory")
	}

	if batchID.Valid {
		batch := int(batchID.Int64)
		row.BatchID = &batch
	}
	row.ApplyNullScalars(createdAt, updatedAt)

	return &row, nil
}

func scanInventories(ctx context.Context, database *db.DB, query string, args ...any) (model.InventoryDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query inventory")
	}
	defer rows.Close()

	var inventories model.InventoryDTOs
	for rows.Next() {
		var row model.InventoryDTO
		var batchID sql.NullInt64
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&row.ID,
			&row.ProductID,
			&row.WarehouseID,
			&row.LocationID,
			&batchID,
			&row.Quantity,
			&row.ReservedQty,
			&row.Version,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan inventory row")
		}

		if batchID.Valid {
			batch := int(batchID.Int64)
			row.BatchID = &batch
		}
		row.ApplyNullScalars(createdAt, updatedAt)
		inventories = append(inventories, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate inventory rows")
	}

	return inventories, nil
}

func scanStockMovements(ctx context.Context, database *db.DB, query string, args ...any) (model.StockMovementDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query stock movements")
	}
	defer rows.Close()

	var movements model.StockMovementDTOs
	for rows.Next() {
		var row model.StockMovementDTO
		var locationIDFrom sql.NullInt64
		var locationIDTo sql.NullInt64
		var batchID sql.NullInt64
		var referenceType sql.NullString
		var referenceID sql.NullInt64
		var performedBy sql.NullInt64
		var notes sql.NullString
		var createdAt sql.NullTime

		err := rows.Scan(
			&row.ID,
			&row.MovementType,
			&row.ProductID,
			&row.WarehouseID,
			&locationIDFrom,
			&locationIDTo,
			&batchID,
			&row.Quantity,
			&referenceType,
			&referenceID,
			&performedBy,
			&notes,
			&createdAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan stock movement row")
		}

		if locationIDFrom.Valid {
			v := int(locationIDFrom.Int64)
			row.LocationIDFrom = &v
		}
		if locationIDTo.Valid {
			v := int(locationIDTo.Int64)
			row.LocationIDTo = &v
		}
		if batchID.Valid {
			v := int(batchID.Int64)
			row.BatchID = &v
		}
		if referenceType.Valid {
			row.ReferenceType = &referenceType.String
		}
		if referenceID.Valid {
			v := int(referenceID.Int64)
			row.ReferenceID = &v
		}
		if performedBy.Valid {
			v := int(performedBy.Int64)
			row.PerformedBy = &v
		}
		if notes.Valid {
			row.Notes = &notes.String
		}
		row.ApplyNullScalars(createdAt)

		movements = append(movements, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate stock movements")
	}

	return movements, nil
}
