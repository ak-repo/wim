package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/internal/db"
	apperrors "github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/model"
	"github.com/jackc/pgx/v5"
)

var ErrPickingTaskNotFound = errors.New("picking task not found")

type PickingRepository interface {
	CreateTask(ctx context.Context, task *model.PickingTaskDTO, items []model.PickingTaskItemDTO) (*model.PickingTaskDTO, error)
	GetTaskByID(ctx context.Context, taskID int) (*model.PickingTaskDTO, error)
	GetTaskByRefCode(ctx context.Context, refCode string) (*model.PickingTaskDTO, error)
	ListTasks(ctx context.Context, params *model.PickingTaskParams) (model.PickingTaskDTOs, error)
	CountTasks(ctx context.Context, params *model.PickingTaskParams) (int, error)
	GetTaskItems(ctx context.Context, taskID int) ([]model.PickingTaskItemDTO, error)
	AssignTask(ctx context.Context, taskID int, assignedTo int, notes string) error
	StartTask(ctx context.Context, taskID int) error
	PickItem(ctx context.Context, taskID int, itemID int, quantity int, locationID int, batchID *int, performedBy int) error
	CompleteTask(ctx context.Context, taskID int, notes string) error
	CancelTask(ctx context.Context, taskID int, notes string) error
}

type pickingRepository struct {
	db *db.DB
}

func NewPickingRepository(database *db.DB) PickingRepository {
	return &pickingRepository{db: database}
}

func (r *pickingRepository) CreateTask(ctx context.Context, task *model.PickingTaskDTO, items []model.PickingTaskItemDTO) (*model.PickingTaskDTO, error) {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start picking task transaction")
	}
	defer tx.Rollback(ctx)

	// Generate ref code
	var refCode string
	if err := tx.QueryRow(ctx, `SELECT 'PK' || TO_CHAR(NOW(), 'YYMMDD') || LPAD(COALESCE(MAX(CAST(SUBSTRING(ref_code FROM 10) AS INTEGER)), 0) + 1::TEXT, 4, '0')
		FROM picking_tasks WHERE ref_code LIKE 'PK' || TO_CHAR(NOW(), 'YYMMDD') || '%'`).Scan(&refCode); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to generate picking task reference")
	}
	task.RefCode = refCode

	// Insert task
	query := `
		INSERT INTO picking_tasks (
			ref_code, sales_order_id, warehouse_id, status, priority, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, task.RefCode, task.SalesOrderID, task.WarehouseID, task.Status, task.Priority, task.CreatedBy).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create picking task")
	}

	// Insert items
	if len(items) > 0 {
		itemQuery := `
			INSERT INTO picking_task_items (
				picking_task_id, sales_order_item_id, product_id, location_id, batch_id, quantity_required, quantity_picked, status, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, 0, $7, NOW(), NOW())
			RETURNING id, created_at, updated_at
		`
		for i := range items {
			items[i].PickingTaskID = task.ID
			items[i].Status = string(model.PickingStatusPending)
			
			var locationID, batchID interface{}
			if items[i].LocationID.Valid {
				locationID = items[i].LocationID.Int64
			}
			if items[i].BatchID.Valid {
				batchID = items[i].BatchID.Int64
			}
			
			err = tx.QueryRow(ctx, itemQuery, items[i].PickingTaskID, items[i].SalesOrderItemID, items[i].ProductID, locationID, batchID, items[i].QuantityRequired).
				Scan(&items[i].ID, &items[i].CreatedAt, &items[i].UpdatedAt)
			if err != nil {
				return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create picking task item")
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit picking task transaction")
	}

	return task, nil
}

func (r *pickingRepository) GetTaskByID(ctx context.Context, taskID int) (*model.PickingTaskDTO, error) {
	return scanPickingTask(ctx, r.db, `
		SELECT id, ref_code, sales_order_id, warehouse_id, status, priority, assigned_to, started_at, completed_at, notes, created_by, created_at, updated_at
		FROM picking_tasks
		WHERE id = $1
	`, taskID)
}

func (r *pickingRepository) GetTaskByRefCode(ctx context.Context, refCode string) (*model.PickingTaskDTO, error) {
	return scanPickingTask(ctx, r.db, `
		SELECT id, ref_code, sales_order_id, warehouse_id, status, priority, assigned_to, started_at, completed_at, notes, created_by, created_at, updated_at
		FROM picking_tasks
		WHERE ref_code = $1
	`, refCode)
}

func (r *pickingRepository) ListTasks(ctx context.Context, params *model.PickingTaskParams) (model.PickingTaskDTOs, error) {
	var args []interface{}
	conditions := []string{"1 = 1"}

	query := `
		SELECT id, ref_code, sales_order_id, warehouse_id, status, priority, assigned_to, started_at, completed_at, notes, created_by, created_at, updated_at
		FROM picking_tasks
	`

	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *params.Status)
	}
	if params.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", len(args)+1))
		args = append(args, *params.Priority)
	}
	if params.AssignedTo != nil {
		conditions = append(conditions, fmt.Sprintf("assigned_to = $%d", len(args)+1))
		args = append(args, *params.AssignedTo)
	}

	whereClause := strings.Join(conditions, " AND ")
	query += " WHERE " + whereClause

	offset := (params.Page - 1) * params.Limit
	limit := params.Limit
	if limit <= 0 {
		limit = 10
	}
	
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	return scanPickingTasks(ctx, r.db, query, args...)
}

func (r *pickingRepository) CountTasks(ctx context.Context, params *model.PickingTaskParams) (int, error) {
	var args []interface{}
	conditions := []string{"1 = 1"}

	query := `SELECT COUNT(*) FROM picking_tasks`

	if params.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("warehouse_id = $%d", len(args)+1))
		args = append(args, *params.WarehouseID)
	}
	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *params.Status)
	}
	if params.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", len(args)+1))
		args = append(args, *params.Priority)
	}
	if params.AssignedTo != nil {
		conditions = append(conditions, fmt.Sprintf("assigned_to = $%d", len(args)+1))
		args = append(args, *params.AssignedTo)
	}

	query += " WHERE " + strings.Join(conditions, " AND ")

	var count int
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count picking tasks")
	}
	return count, nil
}

func (r *pickingRepository) GetTaskItems(ctx context.Context, taskID int) ([]model.PickingTaskItemDTO, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, picking_task_id, sales_order_item_id, product_id, location_id, batch_id, quantity_required, quantity_picked, picked_at, status, created_at, updated_at
		FROM picking_task_items
		WHERE picking_task_id = $1
		ORDER BY id
	`, taskID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query picking task items")
	}
	defer rows.Close()

	var items []model.PickingTaskItemDTO
	for rows.Next() {
		var item model.PickingTaskItemDTO
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.PickingTaskID, &item.SalesOrderItemID, &item.ProductID, &item.LocationID, &item.BatchID, &item.QuantityRequired, &item.QuantityPicked, &item.PickedAt, &item.Status, &createdAt, &updatedAt); err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan picking task item")
		}
		item.ApplyNullScalars(createdAt, updatedAt)
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate picking task items")
	}

	return items, nil
}

func (r *pickingRepository) AssignTask(ctx context.Context, taskID int, assignedTo int, notes string) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start assign transaction")
	}
	defer tx.Rollback(ctx)

	task, err := scanPickingTaskWithTx(ctx, tx, `
		SELECT id, ref_code, sales_order_id, warehouse_id, status, priority, assigned_to, started_at, completed_at, notes, created_by, created_at, updated_at
		FROM picking_tasks
		WHERE id = $1
		FOR UPDATE
	`, taskID)
	if err != nil {
		return err
	}

	if task.Status != string(model.PickingStatusPending) && task.Status != string(model.PickingStatusInProgress) {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot assign this picking task")
	}

	var notesValue interface{}
	if notes != "" {
		notesValue = notes
	}

	if _, err := tx.Exec(ctx, `
		UPDATE picking_tasks
		SET assigned_to = $2, notes = COALESCE($3, notes), updated_at = NOW()
		WHERE id = $1
	`, taskID, assignedTo, notesValue); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to assign picking task")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit assign transaction")
	}

	return nil
}

func (r *pickingRepository) StartTask(ctx context.Context, taskID int) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start task transaction")
	}
	defer tx.Rollback(ctx)

	task, err := scanPickingTaskWithTx(ctx, tx, `
		SELECT id, ref_code, sales_order_id, warehouse_id, status, priority, assigned_to, started_at, completed_at, notes, created_by, created_at, updated_at
		FROM picking_tasks
		WHERE id = $1
		FOR UPDATE
	`, taskID)
	if err != nil {
		return err
	}

	if task.Status != string(model.PickingStatusPending) {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot start this picking task")
	}

	if _, err := tx.Exec(ctx, `
		UPDATE picking_tasks
		SET status = $2, started_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, taskID, model.PickingStatusInProgress); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start picking task")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit start transaction")
	}

	return nil
}

func (r *pickingRepository) PickItem(ctx context.Context, taskID int, itemID int, quantity int, locationID int, batchID *int, performedBy int) error {
	if quantity <= 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "quantity must be greater than 0")
	}

	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start pick transaction")
	}
	defer tx.Rollback(ctx)

	task, err := scanPickingTaskWithTx(ctx, tx, `
		SELECT id, ref_code, sales_order_id, warehouse_id, status, priority, assigned_to, started_at, completed_at, notes, created_by, created_at, updated_at
		FROM picking_tasks
		WHERE id = $1
		FOR UPDATE
	`, taskID)
	if err != nil {
		return err
	}

	if task.Status != string(model.PickingStatusInProgress) {
		return apperrors.New(apperrors.CodeInvalidOperation, "picking task is not in progress")
	}

	item, err := scanPickingTaskItemWithTx(ctx, tx, `
		SELECT id, picking_task_id, sales_order_item_id, product_id, location_id, batch_id, quantity_required, quantity_picked, picked_at, status, created_at, updated_at
		FROM picking_task_items
		WHERE id = $1 AND picking_task_id = $2
		FOR UPDATE
	`, itemID, taskID)
	if err != nil {
		return err
	}

	if item.Status == string(model.PickingStatusCompleted) {
		return apperrors.New(apperrors.CodeInvalidOperation, "item is already completed")
	}

	newPicked := item.QuantityPicked + quantity
	if newPicked > item.QuantityRequired {
		return apperrors.New(apperrors.CodeInvalidInput, "picked quantity exceeds required quantity")
	}

	// Deduct from inventory
	var sourceInventoryID, currentQuantity int
	var batchQuery interface{}
	if batchID != nil {
		batchQuery = *batchID
	} else {
		batchQuery = nil
	}

	err = tx.QueryRow(ctx, `
		SELECT id, quantity
		FROM inventories
		WHERE product_id = $1 AND warehouse_id = $2 AND location_id = $3 AND ((batch_id IS NULL AND $4 IS NULL) OR batch_id = $4)
		FOR UPDATE
	`, item.ProductID, task.WarehouseID, locationID, batchQuery).Scan(&sourceInventoryID, &currentQuantity)
	if err != nil {
		return apperrors.New(apperrors.CodeNotFound, "inventory not found at specified location")
	}

	if currentQuantity < quantity {
		return apperrors.ErrInsufficientStock
	}

	if _, err := tx.Exec(ctx, `
		UPDATE inventories
		SET quantity = quantity - $2, reserved_qty = reserved_qty - $2, version = version + 1, updated_at = NOW()
		WHERE id = $1
	`, sourceInventoryID, quantity); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to deduct inventory")
	}

	// Update picking task item
	itemStatus := string(model.PickingStatusInProgress)
	if newPicked == item.QuantityRequired {
		itemStatus = string(model.PickingStatusCompleted)
	}

	var pickedAt interface{}
	if itemStatus == string(model.PickingStatusCompleted) {
		pickedAt = time.Now().UTC()
	}

	if _, err := tx.Exec(ctx, `
		UPDATE picking_task_items
		SET quantity_picked = $2, status = $3, picked_at = COALESCE($4, picked_at), updated_at = NOW()
		WHERE id = $1
	`, itemID, newPicked, itemStatus, pickedAt); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update picking task item")
	}

	// Record stock movement
	var batchIDValue interface{}
	if batchID != nil {
		batchIDValue = *batchID
	}

	movementNotes := fmt.Sprintf("Picked %d of %d for SO#%d", newPicked, item.QuantityRequired, task.SalesOrderID)
	if _, err := tx.Exec(ctx, `
		INSERT INTO stock_movements (
			movement_type, product_id, warehouse_id, location_id_from, location_id_to, batch_id,
			quantity, reference_type, reference_id, performed_by, notes, created_at
		) VALUES ($1, $2, $3, $4, NULL, $5, $6, $7, $8, $9, $10, NOW())
	`, constants.MovementPick, item.ProductID, task.WarehouseID, locationID, batchIDValue, quantity, constants.ReferenceSalesOrder, task.SalesOrderID, performedBy, movementNotes); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to record stock movement")
	}

	// Update sales order item allocated quantity
	if _, err := tx.Exec(ctx, `
		UPDATE sales_order_items
		SET quantity_reserved = GREATEST(0, quantity_reserved - $2), updated_at = NOW()
		WHERE id = $1
	`, item.SalesOrderItemID, quantity); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update sales order item reservation")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit pick transaction")
	}

	return nil
}

func (r *pickingRepository) CompleteTask(ctx context.Context, taskID int, notes string) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start complete transaction")
	}
	defer tx.Rollback(ctx)

	task, err := scanPickingTaskWithTx(ctx, tx, `
		SELECT id, ref_code, sales_order_id, warehouse_id, status, priority, assigned_to, started_at, completed_at, notes, created_by, created_at, updated_at
		FROM picking_tasks
		WHERE id = $1
		FOR UPDATE
	`, taskID)
	if err != nil {
		return err
	}

	if task.Status != string(model.PickingStatusInProgress) {
		return apperrors.New(apperrors.CodeInvalidOperation, "picking task is not in progress")
	}

	// Check all items are picked
	var incompleteCount int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM picking_task_items
		WHERE picking_task_id = $1 AND status != $2
	`, taskID, model.PickingStatusCompleted).Scan(&incompleteCount); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check task completion")
	}

	if incompleteCount > 0 {
		return apperrors.New(apperrors.CodeInvalidOperation, "not all items are picked")
	}

	var notesValue interface{}
	if notes != "" {
		notesValue = notes
	}

	if _, err := tx.Exec(ctx, `
		UPDATE picking_tasks
		SET status = $2, completed_at = NOW(), notes = COALESCE($3, notes), updated_at = NOW()
		WHERE id = $1
	`, taskID, model.PickingStatusCompleted, notesValue); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to complete picking task")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit complete transaction")
	}

	return nil
}

func (r *pickingRepository) CancelTask(ctx context.Context, taskID int, notes string) error {
	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to start cancel transaction")
	}
	defer tx.Rollback(ctx)

	task, err := scanPickingTaskWithTx(ctx, tx, `
		SELECT id, ref_code, sales_order_id, warehouse_id, status, priority, assigned_to, started_at, completed_at, notes, created_by, created_at, updated_at
		FROM picking_tasks
		WHERE id = $1
		FOR UPDATE
	`, taskID)
	if err != nil {
		return err
	}

	if task.Status == string(model.PickingStatusCompleted) || task.Status == string(model.PickingStatusCancelled) {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot cancel this picking task")
	}

	var notesValue interface{}
	if notes != "" {
		notesValue = fmt.Sprintf("%s\nCancelled: %s", task.Notes.String, notes)
	} else {
		notesValue = fmt.Sprintf("Cancelled: %s", time.Now().Format(time.RFC3339))
	}

	if _, err := tx.Exec(ctx, `
		UPDATE picking_tasks
		SET status = $2, notes = $3, updated_at = NOW()
		WHERE id = $1
	`, taskID, model.PickingStatusCancelled, notesValue); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to cancel picking task")
	}

	if err := tx.Commit(ctx); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to commit cancel transaction")
	}

	return nil
}

func scanPickingTask(ctx context.Context, database *db.DB, query string, args ...interface{}) (*model.PickingTaskDTO, error) {
	var row model.PickingTaskDTO
	var startedAt, completedAt, createdAt, updatedAt sql.NullTime

	if err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID, &row.RefCode, &row.SalesOrderID, &row.WarehouseID, &row.Status, &row.Priority,
		&row.AssignedTo, &startedAt, &completedAt, &row.Notes, &row.CreatedBy, &createdAt, &updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPickingTaskNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load picking task")
	}

	row.ApplyNullScalars(startedAt, completedAt, createdAt, updatedAt)
	return &row, nil
}

func scanPickingTaskWithTx(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) (*model.PickingTaskDTO, error) {
	var row model.PickingTaskDTO
	var startedAt, completedAt, createdAt, updatedAt sql.NullTime

	if err := tx.QueryRow(ctx, query, args...).Scan(
		&row.ID, &row.RefCode, &row.SalesOrderID, &row.WarehouseID, &row.Status, &row.Priority,
		&row.AssignedTo, &startedAt, &completedAt, &row.Notes, &row.CreatedBy, &createdAt, &updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPickingTaskNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load picking task")
	}

	row.ApplyNullScalars(startedAt, completedAt, createdAt, updatedAt)
	return &row, nil
}

func scanPickingTasks(ctx context.Context, database *db.DB, query string, args ...interface{}) (model.PickingTaskDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to query picking tasks")
	}
	defer rows.Close()

	var tasks model.PickingTaskDTOs
	for rows.Next() {
		var row model.PickingTaskDTO
		var startedAt, completedAt, createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&row.ID, &row.RefCode, &row.SalesOrderID, &row.WarehouseID, &row.Status, &row.Priority, &row.AssignedTo, &startedAt, &completedAt, &row.Notes, &row.CreatedBy, &createdAt, &updatedAt); err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan picking task")
		}
		row.ApplyNullScalars(startedAt, completedAt, createdAt, updatedAt)
		tasks = append(tasks, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate picking tasks")
	}

	return tasks, nil
}

func scanPickingTaskItemWithTx(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) (*model.PickingTaskItemDTO, error) {
	var row model.PickingTaskItemDTO
	var createdAt, updatedAt sql.NullTime

	if err := tx.QueryRow(ctx, query, args...).Scan(
		&row.ID, &row.PickingTaskID, &row.SalesOrderItemID, &row.ProductID, &row.LocationID,
		&row.BatchID, &row.QuantityRequired, &row.QuantityPicked, &row.PickedAt, &row.Status, &createdAt, &updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("picking task item not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load picking task item")
	}

	row.ApplyNullScalars(createdAt, updatedAt)
	return &row, nil
}
