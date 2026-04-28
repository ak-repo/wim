package model

import (
	"database/sql"
	"strings"
	"time"
)

type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

// SalesOrderDTO represents a sales order in the database
type SalesOrderDTO struct {
	ID               int        `db:"id"`
	RefCode          string     `db:"ref_code"`
	CustomerID       int        `db:"customer_id"`
	WarehouseID      int        `db:"warehouse_id"`
	Status           string     `db:"status"`
	AllocationStatus string     `db:"allocation_status"`
	OrderDate        time.Time  `db:"order_date"`
	RequiredDate     *time.Time `db:"required_date"`
	ShippedDate      *time.Time `db:"shipped_date"`
	ShippingMethod   *string    `db:"shipping_method"`
	ShippingAddress  *string    `db:"shipping_address"`
	BillingAddress   *string    `db:"billing_address"`
	Notes            *string    `db:"notes"`
	CreatedBy        *int       `db:"created_by"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

// SalesOrderItemDTO represents a sales order item in the database
type SalesOrderItemDTO struct {
	ID                  int       `db:"id"`
	SalesOrderID        int       `db:"sales_order_id"`
	ProductID           int       `db:"product_id"`
	QuantityOrdered     int       `db:"quantity_ordered"`
	QuantityShipped     int       `db:"quantity_shipped"`
	QuantityReserved    int       `db:"quantity_reserved"`
	UnitPrice           *float64  `db:"unit_price"`
	AllocationStatus    string    `db:"allocation_status"`
	BatchID             *int      `db:"batch_id"`
	AllocatedLocationID *int      `db:"allocated_location_id"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}

// SalesOrderAllocationDTO represents inventory allocation for a sales order
type SalesOrderAllocationDTO struct {
	ID                int       `db:"id"`
	SalesOrderID      int       `db:"sales_order_id"`
	SalesOrderItemID  int       `db:"sales_order_item_id"`
	InventoryID       int       `db:"inventory_id"`
	ProductID         int       `db:"product_id"`
	WarehouseID       int       `db:"warehouse_id"`
	LocationID        int       `db:"location_id"`
	BatchID           *int      `db:"batch_id"`
	QuantityAllocated int       `db:"quantity_allocated"`
	IsActive          bool      `db:"is_active"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// SalesOrderParams for listing sales orders
type SalesOrderParams struct {
	CustomerID       *int    `json:"customerId"`
	WarehouseID      *int    `json:"warehouseId"`
	Status           *string `json:"status"`
	AllocationStatus *string `json:"allocationStatus"`
	Page             int     `json:"page"`
	Limit            int     `json:"limit"`
}

// SalesOrderResponse represents the API response
type SalesOrderResponse struct {
	ID               int                       `json:"id"`
	RefCode          string                    `json:"refCode"`
	CustomerID       int                       `json:"customerId"`
	WarehouseID      int                       `json:"warehouseId"`
	Status           string                    `json:"status"`
	AllocationStatus string                    `json:"allocationStatus"`
	OrderDate        time.Time                 `json:"orderDate"`
	RequiredDate     *time.Time                `json:"requiredDate,omitempty"`
	ShippedDate      *time.Time                `json:"shippedDate,omitempty"`
	ShippingMethod   *string                   `json:"shippingMethod,omitempty"`
	ShippingAddress  *string                   `json:"shippingAddress,omitempty"`
	BillingAddress   *string                   `json:"billingAddress,omitempty"`
	Notes            *string                   `json:"notes,omitempty"`
	CreatedBy        *int                      `json:"createdBy,omitempty"`
	CreatedAt        time.Time                 `json:"createdAt"`
	UpdatedAt        time.Time                 `json:"updatedAt"`
	Items            []*SalesOrderItemResponse `json:"items,omitempty"`
}

// SalesOrderItemResponse represents the API response for order items
type SalesOrderItemResponse struct {
	ID                  int       `json:"id"`
	SalesOrderID        int       `json:"salesOrderId"`
	ProductID           int       `json:"productId"`
	QuantityOrdered     int       `json:"quantityOrdered"`
	QuantityShipped     int       `json:"quantityShipped"`
	QuantityReserved    int       `json:"quantityReserved"`
	UnitPrice           *float64  `json:"unitPrice,omitempty"`
	AllocationStatus    string    `json:"allocationStatus"`
	BatchID             *int      `json:"batchId,omitempty"`
	AllocatedLocationID *int      `json:"allocatedLocationId,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// SalesOrderAllocationResponse represents allocation API response
type SalesOrderAllocationResponse struct {
	ID                int       `json:"id"`
	SalesOrderID      int       `json:"salesOrderId"`
	SalesOrderItemID  int       `json:"salesOrderItemId"`
	InventoryID       int       `json:"inventoryId"`
	ProductID         int       `json:"productId"`
	WarehouseID       int       `json:"warehouseId"`
	LocationID        int       `json:"locationId"`
	BatchID           *int      `json:"batchId,omitempty"`
	QuantityAllocated int       `json:"quantityAllocated"`
	IsActive          bool      `json:"isActive"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// SalesOrderDTOs slice with API conversion
type SalesOrderDTOs []*SalesOrderDTO

func (m *SalesOrderDTOs) ToAPIResponse() []*SalesOrderResponse {
	var responses []*SalesOrderResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *SalesOrderDTO) ToAPIResponse() *SalesOrderResponse {
	return &SalesOrderResponse{
		ID:               m.ID,
		RefCode:          m.RefCode,
		CustomerID:       m.CustomerID,
		WarehouseID:      m.WarehouseID,
		Status:           m.Status,
		AllocationStatus: m.AllocationStatus,
		OrderDate:        m.OrderDate,
		RequiredDate:     m.RequiredDate,
		ShippedDate:      m.ShippedDate,
		ShippingMethod:   m.ShippingMethod,
		ShippingAddress:  m.ShippingAddress,
		BillingAddress:   m.BillingAddress,
		Notes:            m.Notes,
		CreatedBy:        m.CreatedBy,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

// SalesOrderItemDTOs slice with API conversion
type SalesOrderItemDTOs []*SalesOrderItemDTO

func (m *SalesOrderItemDTOs) ToAPIResponse() []*SalesOrderItemResponse {
	var responses []*SalesOrderItemResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *SalesOrderItemDTO) ToAPIResponse() *SalesOrderItemResponse {
	return &SalesOrderItemResponse{
		ID:                  m.ID,
		SalesOrderID:        m.SalesOrderID,
		ProductID:           m.ProductID,
		QuantityOrdered:     m.QuantityOrdered,
		QuantityShipped:     m.QuantityShipped,
		QuantityReserved:    m.QuantityReserved,
		UnitPrice:           m.UnitPrice,
		AllocationStatus:    m.AllocationStatus,
		BatchID:             m.BatchID,
		AllocatedLocationID: m.AllocatedLocationID,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func (m *SalesOrderDTO) ApplyNullScalars(createdAt, updatedAt sql.NullTime) {
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

func (m *SalesOrderItemDTO) ApplyNullScalars(createdAt, updatedAt sql.NullTime) {
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

type PurchaseOrderItemRequest struct {
	ProductID       int      `json:"productId"`
	BatchNumber     string   `json:"batchNumber,omitempty"`
	QuantityOrdered int      `json:"quantityOrdered"`
	UnitPrice       *float64 `json:"unitPrice,omitempty"`
}

type CreatePurchaseOrderRequest struct {
	SupplierID   int                        `json:"supplierId"`
	WarehouseID  int                        `json:"warehouseId"`
	ExpectedDate *time.Time                 `json:"expectedDate,omitempty"`
	Notes        string                     `json:"notes,omitempty"`
	Items        []PurchaseOrderItemRequest `json:"items"`
}

type ReceivePurchaseOrderItemRequest struct {
	PurchaseOrderItemID int  `json:"purchaseOrderItemId"`
	QuantityReceived    int  `json:"quantityReceived"`
	LocationID          *int `json:"locationId,omitempty"`
	BatchID             *int `json:"batchId,omitempty"`
}

type ReceivePurchaseOrderRequest struct {
	ReceivedDate *time.Time                        `json:"receivedDate,omitempty"`
	Notes        string                            `json:"notes,omitempty"`
	Items        []ReceivePurchaseOrderItemRequest `json:"items"`
}

type SalesOrderItemRequest struct {
	ProductID       int      `json:"productId"`
	QuantityOrdered int      `json:"quantityOrdered"`
	UnitPrice       *float64 `json:"unitPrice,omitempty"`
}

type CreateSalesOrderRequest struct {
	CustomerID      int                     `json:"customerId"`
	WarehouseID     int                     `json:"warehouseId"`
	RequiredDate    *DateOnly               `json:"requiredDate,omitempty"`
	ShippingMethod  string                  `json:"shippingMethod,omitempty"`
	ShippingAddress string                  `json:"shippingAddress,omitempty"`
	BillingAddress  string                  `json:"billingAddress,omitempty"`
	Notes           string                  `json:"notes,omitempty"`
	Items           []SalesOrderItemRequest `json:"items"`
}

type AllocateSalesOrderRequest struct {
	Strategy string `json:"strategy,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

type ShipSalesOrderItemRequest struct {
	SalesOrderItemID int  `json:"salesOrderItemId"`
	QuantityShipped  int  `json:"quantityShipped"`
	LocationID       *int `json:"locationId,omitempty"`
	BatchID          *int `json:"batchId,omitempty"`
}

type ShipSalesOrderRequest struct {
	ShippedDate *time.Time                  `json:"shippedDate,omitempty"`
	Notes       string                      `json:"notes,omitempty"`
	Items       []ShipSalesOrderItemRequest `json:"items"`
}
