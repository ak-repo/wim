package postgres

import (
	"context"
	"fmt"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

type TransferRepository interface {
	Create(ctx context.Context, transfer *domain.Transfer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transfer, error)
	List(ctx context.Context, filter TransferFilter) ([]*domain.Transfer, error)
	Update(ctx context.Context, transfer *domain.Transfer) error

	CreateItem(ctx context.Context, item *domain.TransferItem) error
	GetItems(ctx context.Context, transferID uuid.UUID) ([]*domain.TransferItem, error)
	UpdateItem(ctx context.Context, item *domain.TransferItem) error
}

type TransferFilter struct {
	SourceWarehouseID *uuid.UUID
	DestWarehouseID   *uuid.UUID
	Status            string
	Limit             int
	Offset            int
}

type transferRepo struct {
	db *DB
}

func NewTransferRepository(db *DB) TransferRepository {
	return &transferRepo{db: db}
}

func (r *transferRepo) Create(ctx context.Context, transfer *domain.Transfer) error {
	query := `
		INSERT INTO transfers (id, transfer_number, source_warehouse_id, dest_warehouse_id, status, requested_by, approved_by, shipped_date, received_date, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Pool.Exec(ctx, query,
		transfer.ID, transfer.TransferNumber, transfer.SourceWarehouseID, transfer.DestWarehouseID,
		transfer.Status, transfer.RequestedBy, transfer.ApprovedBy, transfer.ShippedDate,
		transfer.ReceivedDate, transfer.Notes, transfer.CreatedAt, transfer.UpdatedAt,
	)
	return err
}

func (r *transferRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transfer, error) {
	query := `
		SELECT id, transfer_number, source_warehouse_id, dest_warehouse_id, status, requested_by, approved_by, shipped_date, received_date, notes, created_at, updated_at
		FROM transfers WHERE id = $1`

	var t domain.Transfer
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.TransferNumber, &t.SourceWarehouseID, &t.DestWarehouseID, &t.Status,
		&t.RequestedBy, &t.ApprovedBy, &t.ShippedDate, &t.ReceivedDate, &t.Notes,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *transferRepo) List(ctx context.Context, filter TransferFilter) ([]*domain.Transfer, error) {
	query := `
		SELECT id, transfer_number, source_warehouse_id, dest_warehouse_id, status, requested_by, approved_by, shipped_date, received_date, notes, created_at, updated_at
		FROM transfers WHERE 1=1`

	args := []interface{}{}
	argNum := 1

	if filter.SourceWarehouseID != nil {
		query += fmt.Sprintf(" AND source_warehouse_id = $%d", argNum)
		args = append(args, *filter.SourceWarehouseID)
		argNum++
	}

	if filter.DestWarehouseID != nil {
		query += fmt.Sprintf(" AND dest_warehouse_id = $%d", argNum)
		args = append(args, *filter.DestWarehouseID)
		argNum++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, filter.Status)
		argNum++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argNum)
		args = append(args, filter.Limit)
		argNum++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argNum)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []*domain.Transfer
	for rows.Next() {
		var t domain.Transfer
		err := rows.Scan(
			&t.ID, &t.TransferNumber, &t.SourceWarehouseID, &t.DestWarehouseID, &t.Status,
			&t.RequestedBy, &t.ApprovedBy, &t.ShippedDate, &t.ReceivedDate, &t.Notes,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, &t)
	}

	return transfers, nil
}

func (r *transferRepo) Update(ctx context.Context, transfer *domain.Transfer) error {
	query := `
		UPDATE transfers SET
			transfer_number = $2, source_warehouse_id = $3, dest_warehouse_id = $4,
			status = $5, requested_by = $6, approved_by = $7, shipped_date = $8,
			received_date = $9, notes = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		transfer.ID, transfer.TransferNumber, transfer.SourceWarehouseID, transfer.DestWarehouseID,
		transfer.Status, transfer.RequestedBy, transfer.ApprovedBy, transfer.ShippedDate,
		transfer.ReceivedDate, transfer.Notes, transfer.UpdatedAt,
	)
	return err
}

func (r *transferRepo) CreateItem(ctx context.Context, item *domain.TransferItem) error {
	query := `
		INSERT INTO transfer_items (id, transfer_id, product_id, batch_id, quantity_requested, quantity_shipped, quantity_received, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.TransferID, item.ProductID, item.BatchID, item.QuantityRequested,
		item.QuantityShipped, item.QuantityReceived, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *transferRepo) GetItems(ctx context.Context, transferID uuid.UUID) ([]*domain.TransferItem, error) {
	query := `
		SELECT id, transfer_id, product_id, batch_id, quantity_requested, quantity_shipped, quantity_received, created_at, updated_at
		FROM transfer_items WHERE transfer_id = $1`

	rows, err := r.db.Pool.Query(ctx, query, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.TransferItem
	for rows.Next() {
		var item domain.TransferItem
		err := rows.Scan(
			&item.ID, &item.TransferID, &item.ProductID, &item.BatchID,
			&item.QuantityRequested, &item.QuantityShipped, &item.QuantityReceived,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

func (r *transferRepo) UpdateItem(ctx context.Context, item *domain.TransferItem) error {
	query := `
		UPDATE transfer_items SET
			product_id = $2, batch_id = $3, quantity_requested = $4,
			quantity_shipped = $5, quantity_received = $6, updated_at = $7
		WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.ProductID, item.BatchID, item.QuantityRequested,
		item.QuantityShipped, item.QuantityReceived, item.UpdatedAt,
	)
	return err
}
