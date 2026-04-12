package model

import (
	"database/sql"
	"time"
)

type InventoryDTO struct {
	ID          int       `db:"id"`
	ProductID   int       `db:"product_id"`
	WarehouseID int       `db:"warehouse_id"`
	LocationID  int       `db:"location_id"`
	BatchID     *int      `db:"batch_id"`
	Quantity    int       `db:"quantity"`
	ReservedQty int       `db:"reserved_qty"`
	Version     int       `db:"version"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type InventoryParams struct {
	ProductID   *int `json:"productId"`
	WarehouseID *int `json:"warehouseId"`
	LocationID  *int `json:"locationId"`
	BatchID     *int `json:"batchId"`
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
}

type InventoryDTOs []*InventoryDTO

func (m *InventoryDTOs) ToAPIResponse() []*InventoryResponse {
	var responses []*InventoryResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *InventoryDTO) ToAPIResponse() *InventoryResponse {
	available := m.Quantity - m.ReservedQty
	if available < 0 {
		available = 0
	}
	return &InventoryResponse{
		ID:           m.ID,
		ProductID:    m.ProductID,
		WarehouseID:  m.WarehouseID,
		LocationID:   m.LocationID,
		BatchID:      m.BatchID,
		Quantity:     m.Quantity,
		ReservedQty:  m.ReservedQty,
		AvailableQty: available,
		Version:      m.Version,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (m *InventoryDTO) ApplyNullScalars(createdAt, updatedAt sql.NullTime) {
	if createdAt.Valid {
		m.CreatedAt = createdAt.Time
	} else {
		m.CreatedAt = time.Time{}
	}
	if updatedAt.Valid {
		m.UpdatedAt = updatedAt.Time
	} else {
		m.UpdatedAt = time.Time{}
	}
}

type StockMovementDTO struct {
	ID             int       `db:"id"`
	MovementType   string    `db:"movement_type"`
	ProductID      int       `db:"product_id"`
	WarehouseID    int       `db:"warehouse_id"`
	LocationIDFrom *int      `db:"location_id_from"`
	LocationIDTo   *int      `db:"location_id_to"`
	BatchID        *int      `db:"batch_id"`
	Quantity       int       `db:"quantity"`
	ReferenceType  *string   `db:"reference_type"`
	ReferenceID    *int      `db:"reference_id"`
	PerformedBy    *int      `db:"performed_by"`
	Notes          *string   `db:"notes"`
	CreatedAt      time.Time `db:"created_at"`
}

type StockMovementParams struct {
	MovementType  *string `json:"movementType"`
	ProductID     *int    `json:"productId"`
	WarehouseID   *int    `json:"warehouseId"`
	LocationID    *int    `json:"locationId"`
	BatchID       *int    `json:"batchId"`
	ReferenceType *string `json:"referenceType"`
	ReferenceID   *int    `json:"referenceId"`
	Page          int     `json:"page"`
	Limit         int     `json:"limit"`
}

type StockMovementDTOs []*StockMovementDTO

func (m *StockMovementDTOs) ToAPIResponse() []*StockMovementResponse {
	var responses []*StockMovementResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *StockMovementDTO) ToAPIResponse() *StockMovementResponse {
	refType := ""
	if m.ReferenceType != nil {
		refType = *m.ReferenceType
	}
	notes := ""
	if m.Notes != nil {
		notes = *m.Notes
	}
	return &StockMovementResponse{
		ID:             m.ID,
		MovementType:   m.MovementType,
		ProductID:      m.ProductID,
		WarehouseID:    m.WarehouseID,
		LocationIDFrom: m.LocationIDFrom,
		LocationIDTo:   m.LocationIDTo,
		BatchID:        m.BatchID,
		Quantity:       m.Quantity,
		ReferenceType:  refType,
		ReferenceID:    m.ReferenceID,
		PerformedBy:    m.PerformedBy,
		Notes:          notes,
		CreatedAt:      m.CreatedAt,
	}
}

func (m *StockMovementDTO) ApplyNullScalars(createdAt sql.NullTime) {
	if createdAt.Valid {
		m.CreatedAt = createdAt.Time
	} else {
		m.CreatedAt = time.Time{}
	}
}

// inventory
type AdjustInventoryRequest struct {
	ProductID   int    `json:"productId"`
	WarehouseID int    `json:"warehouseId"`
	LocationID  int    `json:"locationId"`
	BatchID     *int   `json:"batchId,omitempty"`
	Quantity    int    `json:"quantity"`
	Reason      string `json:"reason"`
	Notes       string `json:"notes,omitempty"`
}

type InventoryResponse struct {
	ID           int       `json:"id"`
	ProductID    int       `json:"productId"`
	WarehouseID  int       `json:"warehouseId"`
	LocationID   int       `json:"locationId"`
	BatchID      *int      `json:"batchId,omitempty"`
	Quantity     int       `json:"quantity"`
	ReservedQty  int       `json:"reservedQty"`
	AvailableQty int       `json:"availableQty"`
	Version      int       `json:"version"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type StockMovementResponse struct {
	ID             int       `json:"id"`
	MovementType   string    `json:"movementType"`
	ProductID      int       `json:"productId"`
	WarehouseID    int       `json:"warehouseId"`
	LocationIDFrom *int      `json:"locationIdFrom,omitempty"`
	LocationIDTo   *int      `json:"locationIdTo,omitempty"`
	BatchID        *int      `json:"batchId,omitempty"`
	Quantity       int       `json:"quantity"`
	ReferenceType  string    `json:"referenceType,omitempty"`
	ReferenceID    *int      `json:"referenceId,omitempty"`
	PerformedBy    *int      `json:"performedBy,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}
