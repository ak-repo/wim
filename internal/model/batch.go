package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateBatchRequest struct {
	BatchNumber       string     `json:"batchNumber"`
	ProductID         uuid.UUID  `json:"productId"`
	SupplierID        *uuid.UUID `json:"supplierId,omitempty"`
	ManufacturingDate *time.Time `json:"manufacturingDate,omitempty"`
	ExpiryDate        *time.Time `json:"expiryDate,omitempty"`
	OriginCountry     *string    `json:"originCountry,omitempty"`
	QuantityInitial   int        `json:"quantityInitial"`
}

type BatchResponse struct {
	ID                uuid.UUID  `json:"id"`
	BatchNumber       string     `json:"batchNumber"`
	ProductID         uuid.UUID  `json:"productId"`
	SupplierID        *uuid.UUID `json:"supplierId,omitempty"`
	ManufacturingDate *time.Time `json:"manufacturingDate,omitempty"`
	ExpiryDate        *time.Time `json:"expiryDate,omitempty"`
	OriginCountry     *string    `json:"originCountry,omitempty"`
	QuantityInitial   int        `json:"quantityInitial"`
	QuantityRemaining int        `json:"quantityRemaining"`
	IsActive          bool       `json:"isActive"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
