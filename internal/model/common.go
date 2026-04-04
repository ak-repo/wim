package model

import (
	"time"
)

type AuditLogResponse struct {
	ID         int       `json:"id"`
	EntityType string    `json:"entityType"`
	EntityID   int       `json:"entityId"`
	Action     string    `json:"action"`
	UserID     *int      `json:"userId,omitempty"`
	OldValues  *string   `json:"oldValues,omitempty"`
	NewValues  *string   `json:"newValues,omitempty"`
	IPAddress  *string   `json:"ipAddress,omitempty"`
	UserAgent  *string   `json:"userAgent,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type BatchResponse struct {
	ID                int        `json:"id"`
	BatchNumber       string     `json:"batchNumber"`
	ProductID         int        `json:"productId"`
	SupplierID        *int       `json:"supplierId,omitempty"`
	ManufacturingDate *time.Time `json:"manufacturingDate,omitempty"`
	ExpiryDate        *time.Time `json:"expiryDate,omitempty"`
	OriginCountry     *string    `json:"originCountry,omitempty"`
	QuantityInitial   int        `json:"quantityInitial"`
	QuantityRemaining int        `json:"quantityRemaining"`
	IsActive          bool       `json:"isActive"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
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

// order

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
	RequiredDate    *time.Time              `json:"requiredDate,omitempty"`
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

// report
type InventoryReportRequest struct {
	WarehouseID string     `json:"warehouseId,omitempty"`
	ProductID   string     `json:"productId,omitempty"`
	FromDate    *time.Time `json:"fromDate,omitempty"`
	ToDate      *time.Time `json:"toDate,omitempty"`
}

type MovementReportRequest struct {
	MovementType string     `json:"movementType,omitempty"`
	WarehouseID  string     `json:"warehouseId,omitempty"`
	FromDate     *time.Time `json:"fromDate,omitempty"`
	ToDate       *time.Time `json:"toDate,omitempty"`
}

type ExpiryReportRequest struct {
	WarehouseID string `json:"warehouseId,omitempty"`
	DaysAhead   int    `json:"daysAhead,omitempty"`
}

// transfer

type TransferItemRequest struct {
	ProductID         int  `json:"productId"`
	BatchID           *int `json:"batchId,omitempty"`
	QuantityRequested int  `json:"quantityRequested"`
}

type CreateTransferRequest struct {
	SourceWarehouseID int                   `json:"sourceWarehouseId"`
	DestWarehouseID   int                   `json:"destWarehouseId"`
	Notes             string                `json:"notes,omitempty"`
	Items             []TransferItemRequest `json:"items"`
}

type ShipTransferItemRequest struct {
	TransferItemID   int  `json:"transferItemId"`
	QuantityShipped  int  `json:"quantityShipped"`
	SourceLocationID *int `json:"sourceLocationId,omitempty"`
}

type ShipTransferRequest struct {
	ShippedDate *time.Time                `json:"shippedDate,omitempty"`
	Notes       string                    `json:"notes,omitempty"`
	Items       []ShipTransferItemRequest `json:"items"`
}

type ReceiveTransferItemRequest struct {
	TransferItemID        int  `json:"transferItemId"`
	QuantityReceived      int  `json:"quantityReceived"`
	DestinationLocationID *int `json:"destinationLocationId,omitempty"`
}

type ReceiveTransferRequest struct {
	ReceivedDate *time.Time                   `json:"receivedDate,omitempty"`
	Notes        string                       `json:"notes,omitempty"`
	Items        []ReceiveTransferItemRequest `json:"items"`
}
