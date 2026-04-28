package model

import (
	"database/sql"
	"time"
)

type PurchaseOrderDTO struct {
	ID           int        `db:"id"`
	RefCode      string     `db:"ref_code"`
	SupplierID   int        `db:"supplier_id"`
	WarehouseID  int        `db:"warehouse_id"`
	Status       string     `db:"status"`
	ExpectedDate *time.Time `db:"expected_date"`
	ReceivedDate *time.Time `db:"received_date"`
	Notes        *string    `db:"notes"`
	CreatedBy    *int       `db:"created_by"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

type PurchaseOrderItemDTO struct {
	ID               int       `db:"id"`
	PurchaseOrderID  int       `db:"purchase_order_id"`
	ProductID        int       `db:"product_id"`
	QuantityOrdered  int       `db:"quantity_ordered"`
	QuantityReceived int       `db:"quantity_received"`
	BatchNumber      *string   `db:"batch_number"`
	UnitPrice        *float64  `db:"unit_price"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type PurchaseOrderParams struct {
	SupplierID  *int    `json:"supplierId"`
	WarehouseID *int    `json:"warehouseId"`
	Status      *string `json:"status"`
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
}

type PurchaseOrderResponse struct {
	ID           int                          `json:"id"`
	RefCode      string                       `json:"refCode"`
	SupplierID   int                          `json:"supplierId"`
	WarehouseID  int                          `json:"warehouseId"`
	Status       string                       `json:"status"`
	ExpectedDate *time.Time                   `json:"expectedDate,omitempty"`
	ReceivedDate *time.Time                   `json:"receivedDate,omitempty"`
	Notes        *string                      `json:"notes,omitempty"`
	CreatedBy    *int                         `json:"createdBy,omitempty"`
	CreatedAt    time.Time                    `json:"createdAt"`
	UpdatedAt    time.Time                    `json:"updatedAt"`
	Items        []*PurchaseOrderItemResponse `json:"items,omitempty"`
}

type PurchaseOrderItemResponse struct {
	ID               int       `json:"id"`
	PurchaseOrderID  int       `json:"purchaseOrderId"`
	ProductID        int       `json:"productId"`
	QuantityOrdered  int       `json:"quantityOrdered"`
	QuantityReceived int       `json:"quantityReceived"`
	BatchNumber      *string   `json:"batchNumber,omitempty"`
	UnitPrice        *float64  `json:"unitPrice,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type PurchaseOrderDTOs []*PurchaseOrderDTO

func (m *PurchaseOrderDTOs) ToAPIResponse() []*PurchaseOrderResponse {
	var responses []*PurchaseOrderResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *PurchaseOrderDTO) ToAPIResponse() *PurchaseOrderResponse {
	return &PurchaseOrderResponse{
		ID:           m.ID,
		RefCode:      m.RefCode,
		SupplierID:   m.SupplierID,
		WarehouseID:  m.WarehouseID,
		Status:       m.Status,
		ExpectedDate: m.ExpectedDate,
		ReceivedDate: m.ReceivedDate,
		Notes:        m.Notes,
		CreatedBy:    m.CreatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (m *PurchaseOrderItemDTOs) ToAPIResponse() []*PurchaseOrderItemResponse {
	var responses []*PurchaseOrderItemResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

type PurchaseOrderItemDTOs []*PurchaseOrderItemDTO

func (m *PurchaseOrderItemDTO) ToAPIResponse() *PurchaseOrderItemResponse {
	return &PurchaseOrderItemResponse{
		ID:               m.ID,
		PurchaseOrderID:  m.PurchaseOrderID,
		ProductID:        m.ProductID,
		QuantityOrdered:  m.QuantityOrdered,
		QuantityReceived: m.QuantityReceived,
		BatchNumber:      m.BatchNumber,
		UnitPrice:        m.UnitPrice,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func (m *PurchaseOrderDTO) ApplyNullScalars(expectedDate, receivedDate, createdAt, updatedAt sql.NullTime) {
	if expectedDate.Valid {
		m.ExpectedDate = &expectedDate.Time
	}
	if receivedDate.Valid {
		m.ReceivedDate = &receivedDate.Time
	}
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

func (m *PurchaseOrderItemDTO) ApplyNullScalars(createdAt, updatedAt sql.NullTime) {
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

type PutAwayPurchaseOrderItemRequest struct {
	PurchaseOrderItemID int  `json:"purchaseOrderItemId"`
	Quantity            int  `json:"quantity"`
	FromLocationID      int  `json:"fromLocationId"`
	ToLocationID        int  `json:"toLocationId"`
	BatchID             *int `json:"batchId,omitempty"`
}

type PutAwayPurchaseOrderRequest struct {
	Notes string                            `json:"notes,omitempty"`
	Items []PutAwayPurchaseOrderItemRequest `json:"items"`
}
