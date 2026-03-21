package transfer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/event"
	"github.com/ak-repo/wim/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrTransferNotFound  = errors.New("transfer not found")
	ErrInvalidState      = errors.New("invalid transfer state")
	ErrInsufficientStock = errors.New("insufficient stock for transfer")
)

type Service struct {
	transferRepo      postgres.TransferRepository
	inventoryRepo     postgres.InventoryRepository
	stockMovementRepo postgres.StockMovementRepository
	auditRepo         postgres.AuditLogRepository
	publisher         event.EventPublisher
	db                *postgres.DB
}

func NewService(
	transferRepo postgres.TransferRepository,
	inventoryRepo postgres.InventoryRepository,
	stockMovementRepo postgres.StockMovementRepository,
	auditRepo postgres.AuditLogRepository,
	publisher event.EventPublisher,
	db *postgres.DB,
) *Service {
	return &Service{
		transferRepo:      transferRepo,
		inventoryRepo:     inventoryRepo,
		stockMovementRepo: stockMovementRepo,
		auditRepo:         auditRepo,
		publisher:         publisher,
		db:                db,
	}
}

type CreateTransferItemInput struct {
	ProductID uuid.UUID
	BatchID   *uuid.UUID
	Quantity  int
}

type CreateTransferInput struct {
	SourceWarehouseID uuid.UUID
	DestWarehouseID   uuid.UUID
	RequestedBy       *uuid.UUID
	Notes             *string
	Items             []CreateTransferItemInput
}

func (s *Service) CreateTransfer(ctx context.Context, input CreateTransferInput) (*domain.Transfer, error) {
	transfer := &domain.Transfer{
		ID:                uuid.New(),
		TransferNumber:    generateTransferNumber(),
		SourceWarehouseID: input.SourceWarehouseID,
		DestWarehouseID:   input.DestWarehouseID,
		Status:            "PENDING",
		RequestedBy:       input.RequestedBy,
		Notes:             input.Notes,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.transferRepo.Create(ctx, transfer); err != nil {
		return nil, err
	}

	for _, in := range input.Items {
		item := &domain.TransferItem{
			ID:                uuid.New(),
			TransferID:        transfer.ID,
			ProductID:         in.ProductID,
			BatchID:           in.BatchID,
			QuantityRequested: in.Quantity,
			QuantityShipped:   0,
			QuantityReceived:  0,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		if err := s.transferRepo.CreateItem(ctx, item); err != nil {
			return nil, err
		}
	}

	s.writeAudit(ctx, "TRANSFER", transfer.ID, "CREATE", nil, transfer)
	s.publishEvent(ctx, event.EventTransferCreated, transfer)

	return transfer, nil
}

func (s *Service) GetTransfer(ctx context.Context, id uuid.UUID) (*domain.Transfer, []*domain.TransferItem, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, ErrTransferNotFound
	}

	items, err := s.transferRepo.GetItems(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return transfer, items, nil
}

func (s *Service) ListTransfers(ctx context.Context, filter postgres.TransferFilter) ([]*domain.Transfer, error) {
	return s.transferRepo.List(ctx, filter)
}

func (s *Service) ApproveTransfer(ctx context.Context, id uuid.UUID, approvedBy *uuid.UUID) (*domain.Transfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTransferNotFound
	}

	if transfer.Status != "PENDING" {
		return nil, ErrInvalidState
	}

	transfer.Status = "APPROVED"
	transfer.ApprovedBy = approvedBy
	transfer.UpdatedAt = time.Now()

	if err := s.transferRepo.Update(ctx, transfer); err != nil {
		return nil, err
	}

	s.writeAudit(ctx, "TRANSFER", transfer.ID, "UPDATE", nil, transfer)

	return transfer, nil
}

func (s *Service) ShipTransfer(ctx context.Context, id uuid.UUID) (*domain.Transfer, error) {
	tx, err := s.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var transfer domain.Transfer
	err = tx.QueryRow(ctx, `
		SELECT id, transfer_number, source_warehouse_id, dest_warehouse_id, status, requested_by, approved_by, shipped_date, received_date, notes, created_at, updated_at
		FROM transfers WHERE id = $1
		FOR UPDATE`, id).Scan(
		&transfer.ID, &transfer.TransferNumber, &transfer.SourceWarehouseID, &transfer.DestWarehouseID,
		&transfer.Status, &transfer.RequestedBy, &transfer.ApprovedBy, &transfer.ShippedDate,
		&transfer.ReceivedDate, &transfer.Notes, &transfer.CreatedAt, &transfer.UpdatedAt,
	)
	if err != nil {
		return nil, ErrTransferNotFound
	}

	if transfer.Status != "APPROVED" {
		return nil, ErrInvalidState
	}

	rows, err := tx.Query(ctx, `
		SELECT id, product_id, batch_id, quantity_requested
		FROM transfer_items
		WHERE transfer_id = $1
		FOR UPDATE`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type itemRow struct {
		ID       uuid.UUID
		Product  uuid.UUID
		BatchID  *uuid.UUID
		Quantity int
	}
	var items []itemRow
	for rows.Next() {
		var it itemRow
		if err := rows.Scan(&it.ID, &it.Product, &it.BatchID, &it.Quantity); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	for _, item := range items {
		invRows, err := tx.Query(ctx, `
			SELECT id, location_id, batch_id, quantity, reserved_quantity
			FROM inventory
			WHERE product_id = $1 AND warehouse_id = $2
			ORDER BY created_at ASC
			FOR UPDATE`, item.Product, transfer.SourceWarehouseID)
		if err != nil {
			return nil, err
		}

		type invRow struct {
			ID       uuid.UUID
			Location uuid.UUID
			BatchID  *uuid.UUID
			Qty      int
			Res      int
		}
		var invs []invRow
		for invRows.Next() {
			var inv invRow
			if err := invRows.Scan(&inv.ID, &inv.Location, &inv.BatchID, &inv.Qty, &inv.Res); err != nil {
				invRows.Close()
				return nil, err
			}
			invs = append(invs, inv)
		}
		invRows.Close()

		remaining := item.Quantity
		for _, inv := range invs {
			if remaining <= 0 {
				break
			}
			if item.BatchID != nil {
				if inv.BatchID == nil || *inv.BatchID != *item.BatchID {
					continue
				}
			}
			available := inv.Qty - inv.Res
			if available <= 0 {
				continue
			}
			moveQty := available
			if moveQty > remaining {
				moveQty = remaining
			}

			if _, err := tx.Exec(ctx, `
				UPDATE inventory
				SET quantity = quantity - $1, version = version + 1, updated_at = $2
				WHERE id = $3`, moveQty, time.Now(), inv.ID); err != nil {
				return nil, err
			}

			if _, err := tx.Exec(ctx, `
				INSERT INTO stock_movements (id, movement_type, product_id, warehouse_id, location_id_from, batch_id, quantity, reference_type, reference_id, notes, created_at)
				VALUES ($1, 'TRANSFER_OUT', $2, $3, $4, $5, $6, 'transfer', $7, 'transfer shipped from source warehouse', $8)`,
				uuid.New(), item.Product, transfer.SourceWarehouseID, inv.Location, item.BatchID, moveQty, transfer.ID, time.Now()); err != nil {
				return nil, err
			}

			remaining -= moveQty
		}

		if remaining > 0 {
			return nil, fmt.Errorf("%w: product %s", ErrInsufficientStock, item.Product)
		}

		if _, err := tx.Exec(ctx, `
			UPDATE transfer_items
			SET quantity_shipped = quantity_requested, updated_at = $2
			WHERE id = $1`, item.ID, time.Now()); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	if _, err := tx.Exec(ctx, `
		UPDATE transfers
		SET status = 'IN_TRANSIT', shipped_date = $2, updated_at = $2
		WHERE id = $1`, id, now); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	updated, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.writeAudit(ctx, "TRANSFER", updated.ID, "UPDATE", nil, updated)

	return updated, nil
}

func (s *Service) ReceiveTransfer(ctx context.Context, id uuid.UUID, destinationLocationID uuid.UUID) (*domain.Transfer, error) {
	tx, err := s.db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var transfer domain.Transfer
	err = tx.QueryRow(ctx, `
		SELECT id, transfer_number, source_warehouse_id, dest_warehouse_id, status, requested_by, approved_by, shipped_date, received_date, notes, created_at, updated_at
		FROM transfers WHERE id = $1
		FOR UPDATE`, id).Scan(
		&transfer.ID, &transfer.TransferNumber, &transfer.SourceWarehouseID, &transfer.DestWarehouseID,
		&transfer.Status, &transfer.RequestedBy, &transfer.ApprovedBy, &transfer.ShippedDate,
		&transfer.ReceivedDate, &transfer.Notes, &transfer.CreatedAt, &transfer.UpdatedAt,
	)
	if err != nil {
		return nil, ErrTransferNotFound
	}

	if transfer.Status != "IN_TRANSIT" {
		return nil, ErrInvalidState
	}

	rows, err := tx.Query(ctx, `
		SELECT id, product_id, batch_id, quantity_shipped
		FROM transfer_items
		WHERE transfer_id = $1
		FOR UPDATE`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type itemRow struct {
		ID      uuid.UUID
		Product uuid.UUID
		BatchID *uuid.UUID
		Qty     int
	}
	var items []itemRow
	for rows.Next() {
		var it itemRow
		if err := rows.Scan(&it.ID, &it.Product, &it.BatchID, &it.Qty); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	for _, item := range items {
		if item.Qty <= 0 {
			continue
		}

		var invID uuid.UUID
		err := tx.QueryRow(ctx, `
			SELECT id
			FROM inventory
			WHERE product_id = $1 AND warehouse_id = $2 AND location_id = $3
			AND (($4::uuid IS NULL AND batch_id IS NULL) OR batch_id = $4)
			FOR UPDATE`, item.Product, transfer.DestWarehouseID, destinationLocationID, item.BatchID).Scan(&invID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				if _, err := tx.Exec(ctx, `
					INSERT INTO inventory (id, product_id, warehouse_id, location_id, batch_id, quantity, reserved_quantity, version, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, 0, 1, $7, $7)`,
					uuid.New(), item.Product, transfer.DestWarehouseID, destinationLocationID, item.BatchID, item.Qty, time.Now()); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			if _, err := tx.Exec(ctx, `
				UPDATE inventory
				SET quantity = quantity + $1, version = version + 1, updated_at = $2
				WHERE id = $3`, item.Qty, time.Now(), invID); err != nil {
				return nil, err
			}
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO stock_movements (id, movement_type, product_id, warehouse_id, location_id_to, batch_id, quantity, reference_type, reference_id, notes, created_at)
			VALUES ($1, 'TRANSFER_IN', $2, $3, $4, $5, $6, 'transfer', $7, 'transfer received in destination warehouse', $8)`,
			uuid.New(), item.Product, transfer.DestWarehouseID, destinationLocationID, item.BatchID, item.Qty, transfer.ID, time.Now()); err != nil {
			return nil, err
		}

		if _, err := tx.Exec(ctx, `
			UPDATE transfer_items
			SET quantity_received = quantity_shipped, updated_at = $2
			WHERE id = $1`, item.ID, time.Now()); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	if _, err := tx.Exec(ctx, `
		UPDATE transfers
		SET status = 'COMPLETED', received_date = $2, updated_at = $2
		WHERE id = $1`, id, now); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	updated, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.writeAudit(ctx, "TRANSFER", updated.ID, "UPDATE", nil, updated)
	s.publishEvent(ctx, event.EventTransferCompleted, updated)

	return updated, nil
}

func generateTransferNumber() string {
	return "TR-" + time.Now().Format("20060102150405")
}

func (s *Service) writeAudit(ctx context.Context, entityType string, entityID uuid.UUID, action string, oldValue any, newValue any) {
	if s.auditRepo == nil {
		return
	}

	var oldJSON *string
	if oldValue != nil {
		if b, err := json.Marshal(oldValue); err == nil {
			v := string(b)
			oldJSON = &v
		}
	}

	var newJSON *string
	if newValue != nil {
		if b, err := json.Marshal(newValue); err == nil {
			v := string(b)
			newJSON = &v
		}
	}

	_ = s.auditRepo.Create(ctx, &domain.AuditLog{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		OldValues:  oldJSON,
		NewValues:  newJSON,
		CreatedAt:  time.Now(),
	})
}

func (s *Service) publishEvent(ctx context.Context, eventType event.EventType, payload any) {
	if s.publisher == nil {
		return
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return
	}

	_ = s.publisher.Publish(ctx, event.Event{
		ID:        uuid.NewString(),
		Type:      eventType,
		Payload:   b,
		Timestamp: time.Now(),
	})
}
