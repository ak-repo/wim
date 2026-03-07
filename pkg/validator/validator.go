package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	v := validator.New()
	return &Validator{validate: v}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

type ProductCreateRequest struct {
	SKU           string   `json:"sku" validate:"required,max=50"`
	Name          string   `json:"name" validate:"required,max=255"`
	Description   string   `json:"description"`
	Category      string   `json:"category" validate:"max=100"`
	UnitOfMeasure string   `json:"unit_of_measure" validate:"required,max=20"`
	Weight        *float64 `json:"weight"`
	Barcode       string   `json:"barcode" validate:"max=50"`
	IsActive      *bool    `json:"is_active"`
}

type WarehouseCreateRequest struct {
	Code         string `json:"code" validate:"required,max=20"`
	Name         string `json:"name" validate:"required,max=255"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city" validate:"max=100"`
	State        string `json:"state" validate:"max=100"`
	PostalCode   string `json:"postal_code" validate:"max=20"`
	Country      string `json:"country" validate:"required,max=2"`
	IsActive     *bool  `json:"is_active"`
}

type LocationCreateRequest struct {
	WarehouseID  string   `json:"warehouse_id" validate:"required,uuid"`
	Zone         string   `json:"zone" validate:"required,max=20"`
	Aisle        string   `json:"aisle" validate:"required,max=20"`
	Rack         string   `json:"rack" validate:"required,max=20"`
	Bin          string   `json:"bin" validate:"required,max=20"`
	LocationCode string   `json:"location_code" validate:"required,max=50"`
	LocationType string   `json:"location_type" validate:"required,max=20"`
	IsPickFace   *bool    `json:"is_pick_face"`
	MaxWeight    *float64 `json:"max_weight"`
	IsActive     *bool    `json:"is_active"`
}

type InventoryAdjustRequest struct {
	ProductID    string  `json:"product_id" validate:"required,uuid"`
	WarehouseID  string  `json:"warehouse_id" validate:"required,uuid"`
	LocationID   string  `json:"location_id" validate:"required,uuid"`
	BatchID      *string `json:"batch_id" validate:"omitempty,uuid"`
	Quantity     int     `json:"quantity" validate:"required"`
	MovementType string  `json:"movement_type" validate:"required"`
	Notes        string  `json:"notes"`
}

type PurchaseOrderCreateRequest struct {
	SupplierID   string                     `json:"supplier_id" validate:"required,uuid"`
	WarehouseID  string                     `json:"warehouse_id" validate:"required,uuid"`
	ExpectedDate *string                    `json:"expected_date"`
	Notes        string                     `json:"notes"`
	Items        []PurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type PurchaseOrderItemRequest struct {
	ProductID       string   `json:"product_id" validate:"required,uuid"`
	QuantityOrdered int      `json:"quantity_ordered" validate:"required,min=1"`
	UnitPrice       *float64 `json:"unit_price"`
}

type SalesOrderCreateRequest struct {
	CustomerID      string                  `json:"customer_id" validate:"required,uuid"`
	WarehouseID     string                  `json:"warehouse_id" validate:"required,uuid"`
	RequiredDate    *string                 `json:"required_date"`
	ShippingMethod  string                  `json:"shipping_method"`
	ShippingAddress string                  `json:"shipping_address"`
	BillingAddress  string                  `json:"billing_address"`
	Notes           string                  `json:"notes"`
	Items           []SalesOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type SalesOrderItemRequest struct {
	ProductID string   `json:"product_id" validate:"required,uuid"`
	Quantity  int      `json:"quantity" validate:"required,min=1"`
	UnitPrice *float64 `json:"unit_price"`
}
