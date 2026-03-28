package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	SKU           string   `json:"sku"`
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Category      string   `json:"category,omitempty"`
	UnitOfMeasure string   `json:"unitOfMeasure"`
	Weight        *float64 `json:"weight,omitempty"`
	Length        *float64 `json:"length,omitempty"`
	Width         *float64 `json:"width,omitempty"`
	Height        *float64 `json:"height,omitempty"`
	Barcode       string   `json:"barcode,omitempty"`
}

type UpdateProductRequest struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Category      string   `json:"category,omitempty"`
	UnitOfMeasure string   `json:"unitOfMeasure,omitempty"`
	Weight        *float64 `json:"weight,omitempty"`
	Length        *float64 `json:"length,omitempty"`
	Width         *float64 `json:"width,omitempty"`
	Height        *float64 `json:"height,omitempty"`
	Barcode       string   `json:"barcode,omitempty"`
	IsActive      *bool    `json:"isActive,omitempty"`
}

type ProductResponse struct {
	ID            uuid.UUID `json:"id"`
	SKU           string    `json:"sku"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Category      string    `json:"category,omitempty"`
	UnitOfMeasure string    `json:"unitOfMeasure"`
	Weight        *float64  `json:"weight,omitempty"`
	Length        *float64  `json:"length,omitempty"`
	Width         *float64  `json:"width,omitempty"`
	Height        *float64  `json:"height,omitempty"`
	Barcode       string    `json:"barcode,omitempty"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
