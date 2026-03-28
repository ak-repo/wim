package model

import (
	"time"

	"github.com/google/uuid"
)

type AdjustInventoryRequest struct {
	ProductID   uuid.UUID  `json:"productId"`
	WarehouseID uuid.UUID  `json:"warehouseId"`
	LocationID  uuid.UUID  `json:"locationId"`
	BatchID     *uuid.UUID `json:"batchId,omitempty"`
	Quantity    int        `json:"quantity"`
	Reason      string     `json:"reason"`
	Notes       string     `json:"notes,omitempty"`
}

type InventoryResponse struct {
	ID           uuid.UUID  `json:"id"`
	ProductID    uuid.UUID  `json:"productId"`
	WarehouseID  uuid.UUID  `json:"warehouseId"`
	LocationID   uuid.UUID  `json:"locationId"`
	BatchID      *uuid.UUID `json:"batchId,omitempty"`
	Quantity     int        `json:"quantity"`
	ReservedQty  int        `json:"reservedQty"`
	AvailableQty int        `json:"availableQty"`
	Version      int        `json:"version"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type StockMovementResponse struct {
	ID             uuid.UUID  `json:"id"`
	MovementType   string     `json:"movementType"`
	ProductID      uuid.UUID  `json:"productId"`
	WarehouseID    uuid.UUID  `json:"warehouseId"`
	LocationIDFrom *uuid.UUID `json:"locationIdFrom,omitempty"`
	LocationIDTo   *uuid.UUID `json:"locationIdTo,omitempty"`
	BatchID        *uuid.UUID `json:"batchId,omitempty"`
	Quantity       int        `json:"quantity"`
	ReferenceType  string     `json:"referenceType,omitempty"`
	ReferenceID    *uuid.UUID `json:"referenceId,omitempty"`
	PerformedBy    *uuid.UUID `json:"performedBy,omitempty"`
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
}
