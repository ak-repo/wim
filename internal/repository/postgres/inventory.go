package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

var ErrConcurrentUpdate = fmt.Errorf("concurrent update detected")

type InventoryRepository interface {
	Create(ctx context.Context, inv *domain.Inventory) error
	List(ctx context.Context, limit, offset int) ([]*domain.Inventory, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error)
	GetByProductWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) ([]*domain.Inventory, error)
	GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]*domain.Inventory, error)
	GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.Inventory, error)
	GetTotalQuantity(ctx context.Context, productID, warehouseID uuid.UUID) (int, error)
	Update(ctx context.Context, inv *domain.Inventory) error
	UpdateWithVersion(ctx context.Context, inv *domain.Inventory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func (r *inventoryRepo) List(ctx context.Context, limit, offset int) ([]*domain.Inventory, error) {
	query := `
		SELECT id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_quantity, version, last_movement_id, created_at, updated_at
		FROM inventory
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invs []*domain.Inventory
	for rows.Next() {
		var inv domain.Inventory
		err := rows.Scan(
			&inv.ID, &inv.ProductID, &inv.WarehouseID, &inv.LocationID, &inv.BatchID,
			&inv.Quantity, &inv.ReservedQty, &inv.Version, &inv.LastMovementID, &inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invs = append(invs, &inv)
	}

	return invs, nil
}

type inventoryRepo struct {
	db    *DB
	redis *Redis
}

func NewInventoryRepository(db *DB, redis *Redis) InventoryRepository {
	return &inventoryRepo{db: db, redis: redis}
}

func (r *inventoryRepo) Create(ctx context.Context, inv *domain.Inventory) error {
	query := `
		INSERT INTO inventory (id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_quantity, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Pool.Exec(ctx, query,
		inv.ID, inv.ProductID, inv.WarehouseID, inv.LocationID, inv.BatchID,
		inv.Quantity, inv.ReservedQty, inv.Version, inv.CreatedAt, inv.UpdatedAt,
	)
	return err
}

func (r *inventoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	cacheKey := "inv:id:" + id.String()
	if r.redis != nil {
		cached, err := r.redis.Client.Get(ctx, cacheKey).Result()
		if err == nil {
			var inv domain.Inventory
			if err := json.Unmarshal([]byte(cached), &inv); err == nil {
				return &inv, nil
			}
		}
	}

	query := `
		SELECT id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_quantity, version, last_movement_id, created_at, updated_at
		FROM inventory WHERE id = $1`

	var inv domain.Inventory
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&inv.ID, &inv.ProductID, &inv.WarehouseID, &inv.LocationID, &inv.BatchID,
		&inv.Quantity, &inv.ReservedQty, &inv.Version, &inv.LastMovementID, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if r.redis != nil {
		data, _ := json.Marshal(inv)
		r.redis.Client.Set(ctx, cacheKey, data, 30)
	}

	return &inv, nil
}

func (r *inventoryRepo) GetByProductWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) ([]*domain.Inventory, error) {
	query := `
		SELECT id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_quantity, version, last_movement_id, created_at, updated_at
		FROM inventory 
		WHERE product_id = $1 AND warehouse_id = $2`

	rows, err := r.db.Pool.Query(ctx, query, productID, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invs []*domain.Inventory
	for rows.Next() {
		var inv domain.Inventory
		err := rows.Scan(
			&inv.ID, &inv.ProductID, &inv.WarehouseID, &inv.LocationID, &inv.BatchID,
			&inv.Quantity, &inv.ReservedQty, &inv.Version, &inv.LastMovementID, &inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invs = append(invs, &inv)
	}

	return invs, nil
}

func (r *inventoryRepo) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]*domain.Inventory, error) {
	query := `
		SELECT id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_quantity, version, last_movement_id, created_at, updated_at
		FROM inventory 
		WHERE warehouse_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Pool.Query(ctx, query, warehouseID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invs []*domain.Inventory
	for rows.Next() {
		var inv domain.Inventory
		err := rows.Scan(
			&inv.ID, &inv.ProductID, &inv.WarehouseID, &inv.LocationID, &inv.BatchID,
			&inv.Quantity, &inv.ReservedQty, &inv.Version, &inv.LastMovementID, &inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invs = append(invs, &inv)
	}

	return invs, nil
}

func (r *inventoryRepo) GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.Inventory, error) {
	query := `
		SELECT id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_quantity, version, last_movement_id, created_at, updated_at
		FROM inventory 
		WHERE product_id = $1`

	rows, err := r.db.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invs []*domain.Inventory
	for rows.Next() {
		var inv domain.Inventory
		err := rows.Scan(
			&inv.ID, &inv.ProductID, &inv.WarehouseID, &inv.LocationID, &inv.BatchID,
			&inv.Quantity, &inv.ReservedQty, &inv.Version, &inv.LastMovementID, &inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invs = append(invs, &inv)
	}

	return invs, nil
}

func (r *inventoryRepo) GetTotalQuantity(ctx context.Context, productID, warehouseID uuid.UUID) (int, error) {
	cacheKey := "inv:total:" + productID.String() + ":" + warehouseID.String()
	if r.redis != nil {
		cached, err := r.redis.Client.Get(ctx, cacheKey).Result()
		if err == nil {
			var qty int
			if _, err := fmt.Sscan(cached, &qty); err == nil {
				return qty, nil
			}
		}
	}

	query := `
		SELECT COALESCE(SUM(quantity - reserved_quantity), 0)
		FROM inventory 
		WHERE product_id = $1 AND warehouse_id = $2`

	var qty int
	err := r.db.Pool.QueryRow(ctx, query, productID, warehouseID).Scan(&qty)
	if err != nil {
		return 0, err
	}

	if r.redis != nil {
		r.redis.Client.Set(ctx, cacheKey, fmt.Sprintf("%d", qty), 30)
	}

	return qty, nil
}

func (r *inventoryRepo) Update(ctx context.Context, inv *domain.Inventory) error {
	query := `
		UPDATE inventory SET 
			product_id = $2, warehouse_id = $3, location_id = $4, batch_id = $5,
			quantity = $6, reserved_quantity = $7, version = version + 1, 
			last_movement_id = $8, updated_at = $9
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		inv.ID, inv.ProductID, inv.WarehouseID, inv.LocationID, inv.BatchID,
		inv.Quantity, inv.ReservedQty, inv.LastMovementID, inv.UpdatedAt,
	)
	return err
}

func (r *inventoryRepo) UpdateWithVersion(ctx context.Context, inv *domain.Inventory) error {
	query := `
		UPDATE inventory SET 
			quantity = $2, reserved_quantity = $3, version = version + 1, 
			last_movement_id = $4, updated_at = $5
		WHERE id = $1 AND version = $6`

	result, err := r.db.Pool.Exec(ctx, query,
		inv.ID, inv.Quantity, inv.ReservedQty, inv.LastMovementID, inv.UpdatedAt, inv.Version,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrConcurrentUpdate
	}

	return nil
}

func (r *inventoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM inventory WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
