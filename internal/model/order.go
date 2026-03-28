package model

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrderItemRequest struct {
	ProductID       uuid.UUID `json:"productId"`
	BatchNumber     string    `json:"batchNumber,omitempty"`
	QuantityOrdered int       `json:"quantityOrdered"`
	UnitPrice       *float64  `json:"unitPrice,omitempty"`
}

type CreatePurchaseOrderRequest struct {
	SupplierID   uuid.UUID                  `json:"supplierId"`
	WarehouseID  uuid.UUID                  `json:"warehouseId"`
	ExpectedDate *time.Time                 `json:"expectedDate,omitempty"`
	Notes        string                     `json:"notes,omitempty"`
	Items        []PurchaseOrderItemRequest `json:"items"`
}

type ReceivePurchaseOrderItemRequest struct {
	PurchaseOrderItemID uuid.UUID  `json:"purchaseOrderItemId"`
	QuantityReceived    int        `json:"quantityReceived"`
	LocationID          *uuid.UUID `json:"locationId,omitempty"`
	BatchID             *uuid.UUID `json:"batchId,omitempty"`
}

type ReceivePurchaseOrderRequest struct {
	ReceivedDate *time.Time                        `json:"receivedDate,omitempty"`
	Notes        string                            `json:"notes,omitempty"`
	Items        []ReceivePurchaseOrderItemRequest `json:"items"`
}

type SalesOrderItemRequest struct {
	ProductID       uuid.UUID `json:"productId"`
	QuantityOrdered int       `json:"quantityOrdered"`
	UnitPrice       *float64  `json:"unitPrice,omitempty"`
}

type CreateSalesOrderRequest struct {
	CustomerID      uuid.UUID               `json:"customerId"`
	WarehouseID     uuid.UUID               `json:"warehouseId"`
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
	SalesOrderItemID uuid.UUID  `json:"salesOrderItemId"`
	QuantityShipped  int        `json:"quantityShipped"`
	LocationID       *uuid.UUID `json:"locationId,omitempty"`
	BatchID          *uuid.UUID `json:"batchId,omitempty"`
}

type ShipSalesOrderRequest struct {
	ShippedDate *time.Time                  `json:"shippedDate,omitempty"`
	Notes       string                      `json:"notes,omitempty"`
	Items       []ShipSalesOrderItemRequest `json:"items"`
}
