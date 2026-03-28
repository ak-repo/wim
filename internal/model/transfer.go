package model

import (
	"time"

	"github.com/google/uuid"
)

type TransferItemRequest struct {
	ProductID         uuid.UUID  `json:"productId"`
	BatchID           *uuid.UUID `json:"batchId,omitempty"`
	QuantityRequested int        `json:"quantityRequested"`
}

type CreateTransferRequest struct {
	SourceWarehouseID uuid.UUID             `json:"sourceWarehouseId"`
	DestWarehouseID   uuid.UUID             `json:"destWarehouseId"`
	Notes             string                `json:"notes,omitempty"`
	Items             []TransferItemRequest `json:"items"`
}

type ShipTransferItemRequest struct {
	TransferItemID   uuid.UUID  `json:"transferItemId"`
	QuantityShipped  int        `json:"quantityShipped"`
	SourceLocationID *uuid.UUID `json:"sourceLocationId,omitempty"`
}

type ShipTransferRequest struct {
	ShippedDate *time.Time                `json:"shippedDate,omitempty"`
	Notes       string                    `json:"notes,omitempty"`
	Items       []ShipTransferItemRequest `json:"items"`
}

type ReceiveTransferItemRequest struct {
	TransferItemID        uuid.UUID  `json:"transferItemId"`
	QuantityReceived      int        `json:"quantityReceived"`
	DestinationLocationID *uuid.UUID `json:"destinationLocationId,omitempty"`
}

type ReceiveTransferRequest struct {
	ReceivedDate *time.Time                   `json:"receivedDate,omitempty"`
	Notes        string                       `json:"notes,omitempty"`
	Items        []ReceiveTransferItemRequest `json:"items"`
}
