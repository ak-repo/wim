package model

import (
	"time"
)

type ProductDTO struct {
	ID            int       `db:"id"`
	SKU           string    `db:"sku"`
	Name          string    `db:"name"`
	Description   *string   `db:"description"`
	Category      *string   `db:"category"`
	UnitOfMeasure string    `db:"unit_of_measure"`
	Weight        *float64  `db:"weight"`
	Length        *float64  `db:"length"`
	Width         *float64  `db:"width"`
	Height        *float64  `db:"height"`
	Barcode       *string   `db:"barcode"`
	IsActive      bool      `db:"is_active"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type ProductResponse struct {
	ID            int       `json:"id"`
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

type ProductParams struct {
	Active   *bool  `json:"active"`
	Category string `json:"category,omitempty"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type ProductDTOs []*ProductDTO

// API conversion
func (m *ProductDTOs) ToAPIResponse() []*ProductResponse {
	var responses []*ProductResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *ProductDTO) ToAPIResponse() *ProductResponse {
	return &ProductResponse{
		ID:            m.ID,
		SKU:           m.SKU,
		Name:          m.Name,
		Description:   *m.Description,
		Category:      *m.Category,
		UnitOfMeasure: m.UnitOfMeasure,
		Weight:        m.Weight,
		Length:        m.Length,
		Width:         m.Width,
		Height:        m.Height,
		Barcode:       *m.Barcode,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
