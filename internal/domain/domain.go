package domain

import (
	"time"

	"github.com/google/uuid"
)

type Batch struct {
	ID                uuid.UUID
	BatchNumber       string
	ProductID         uuid.UUID
	SupplierID        *uuid.UUID
	ManufacturingDate *time.Time
	ExpiryDate        *time.Time
	OriginCountry     *string
	QuantityInitial   int
	QuantityRemaining int
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Inventory struct {
	ID             uuid.UUID
	ProductID      uuid.UUID
	WarehouseID    uuid.UUID
	LocationID     uuid.UUID
	BatchID        *uuid.UUID
	Quantity       int
	ReservedQty    int
	Version        int
	LastMovementID *uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (i *Inventory) AvailableQty() int {
	return i.Quantity - i.ReservedQty
}

type StockMovement struct {
	ID             uuid.UUID
	MovementType   string
	ProductID      uuid.UUID
	WarehouseID    uuid.UUID
	LocationIDFrom *uuid.UUID
	LocationIDTo   *uuid.UUID
	BatchID        *uuid.UUID
	Quantity       int
	ReferenceType  string
	ReferenceID    *uuid.UUID
	PerformedBy    *uuid.UUID
	Notes          string
	CreatedAt      time.Time
}

type PurchaseOrder struct {
	ID           uuid.UUID
	PONumber     string
	SupplierID   uuid.UUID
	WarehouseID  uuid.UUID
	OrderDate    time.Time
	ExpectedDate *time.Time
	ReceivedDate *time.Time
	Status       string
	TotalAmount  *float64
	Notes        string
	CreatedBy    *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PurchaseOrderItem struct {
	ID               uuid.UUID
	PurchaseOrderID  uuid.UUID
	ProductID        uuid.UUID
	BatchNumber      *string
	QuantityOrdered  int
	QuantityReceived int
	UnitPrice        *float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SalesOrder struct {
	ID               uuid.UUID
	OrderNumber      string
	CustomerID       uuid.UUID
	WarehouseID      uuid.UUID
	OrderDate        time.Time
	RequiredDate     *time.Time
	ShippedDate      *time.Time
	Status           string
	AllocationStatus string
	ShippingMethod   *string
	ShippingAddress  *string
	BillingAddress   *string
	Subtotal         *float64
	TaxAmount        *float64
	ShippingAmount   *float64
	TotalAmount      *float64
	Notes            *string
	CreatedBy        *uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SalesOrderItem struct {
	ID                uuid.UUID
	SalesOrderID      uuid.UUID
	ProductID         uuid.UUID
	BatchID           *uuid.UUID
	LocationID        *uuid.UUID
	QuantityOrdered   int
	QuantityAllocated int
	QuantityPicked    int
	QuantityShipped   int
	UnitPrice         *float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Transfer struct {
	ID                uuid.UUID
	TransferNumber    string
	SourceWarehouseID uuid.UUID
	DestWarehouseID   uuid.UUID
	Status            string
	RequestedBy       *uuid.UUID
	ApprovedBy        *uuid.UUID
	ShippedDate       *time.Time
	ReceivedDate      *time.Time
	Notes             *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type TransferItem struct {
	ID                uuid.UUID
	TransferID        uuid.UUID
	ProductID         uuid.UUID
	BatchID           *uuid.UUID
	QuantityRequested int
	QuantityShipped   int
	QuantityReceived  int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Barcode struct {
	ID           uuid.UUID
	ProductID    uuid.UUID
	BarcodeValue string
	BarcodeType  string
	IsPrimary    bool
	CreatedAt    time.Time
}

type AuditLog struct {
	ID         uuid.UUID
	EntityType string
	EntityID   uuid.UUID
	Action     string
	UserID     *uuid.UUID
	OldValues  *string
	NewValues  *string
	IPAddress  *string
	UserAgent  *string
	CreatedAt  time.Time
}
