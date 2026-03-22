package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID
	SKU           string
	Name          string
	Description   string
	Category      string
	UnitOfMeasure string
	Weight        *float64
	Length        *float64
	Width         *float64
	Height        *float64
	Barcode       string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewProduct(sku, name, description, category, unitOfMeasure, barcode string) *Product {
	return &Product{
		ID:            uuid.New(),
		SKU:           sku,
		Name:          name,
		Description:   description,
		Category:      category,
		UnitOfMeasure: unitOfMeasure,
		Barcode:       barcode,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
